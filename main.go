package main

import (
	"flag"
	"go-cloc/clone"
	"go-cloc/logger"
	"go-cloc/report"
	"go-cloc/scanner"
	"log"
	"os"
	"path/filepath"
	"time"
)

// pseduocode
// discover repos, should return a list of repos

// for each repo
// // clone repo
// // perform a scan
// // dump a csv report

// combine csv reports
// report on any failed repos

// Modes
const (
	LOCAL string = "Local"
	GH    string = "GitHub"
)

func DiscoverRepositories(mode string, accessToken string, organization string) []clone.RepoInfo {
	logger.Info("Discovering repositories in ", organization)
	repositoryInfoArr := clone.DiscoverReposGithub(organization, accessToken)

	return repositoryInfoArr
}

// contains checks if a slice contains a specific string
func contains(slice []string, item string) bool {
	for _, str := range slice {
		if str == item {
			return true
		}
	}
	return false
}

func main() {
	// Parse command-line flags
	logLevelArg := flag.String("log-level", "INFO", "Log level (DEBUG, INFO, WARN, ERROR)")
	modeArg := flag.String("devops", LOCAL, "flag : <GitHub>||<File>")
	localScanFilePathArg := flag.String("local-file-path", "", "Path to your local file or directory that you want to scan")
	accessTokenArg := flag.String("accessToken", "", "Your DevOps personal access token used for discovering and downloading repositories in your organization")
	organizationArg := flag.String("organization", "", "Your DevOps organization name")
	ignoreFilePathArg := flag.String("ignore-file", "", "(Optional) Path to your ignore file to exclude directories and files. Please see the README.md for how to format your ignore configuration")
	excludeRepositoriesFilePathArg := flag.String("exclude-repositories-file", "", "(Optional) Path to your exclude repositories file to exclude repositories. Please see the README.md for how to format your exclude repositories configuration")
	includeRepositoriesFilePathArg := flag.String("include-repositories-file", "", "(Optional) Path to your include repositories file to include repositories. Please see the README.md for how to format your include repositories configuration")

	flag.Parse()

	// dereference all CLI args to make it easier to use later
	logLevel := *logLevelArg
	mode := *modeArg
	localScanFilePath := *localScanFilePathArg
	accessToken := *accessTokenArg
	organization := *organizationArg
	ignoreFilePath := *ignoreFilePathArg
	excludeRepositoriesFilePath := *excludeRepositoriesFilePathArg
	includeRepositoriesFilePath := *includeRepositoriesFilePathArg

	// set log level
	logger.Info("Setting Log Level to " + logLevel)
	logger.SetLogLevel(logger.ConvertStringToLogLevel(logLevel))
	logger.SetOutput(os.Stdout)

	// validate mandatory arguments
	if mode == GH {
		if organization == "" || accessToken == "" {
			logger.Error("Mode ", mode, " requires : --organization & --accessToken")
			os.Exit(-1)
		}
	} else if mode == LOCAL {
		if localScanFilePath == "" {
			logger.Error("Mode ", mode, " requires : --local-file-path")
			os.Exit(-1)
		}
	}

	// validate optional arguments

	// parse ignore patterns
	ignorePatterns := []string{}
	if ignoreFilePath != "" {
		temp := scanner.ReadIgnoreFile(ignoreFilePath)
		if temp == nil {
			logger.Error("Error reading ignore-file ", ignoreFilePath)
			os.Exit(-1)
		}
		logger.Debug("Successfully read in the ignore-file ", ignoreFilePath)
		ignorePatterns = temp
	}

	// parse exclude repositories
	excludeRepositories := []string{}
	if excludeRepositoriesFilePath != "" {
		temp := scanner.ReadIgnoreFile(excludeRepositoriesFilePath)
		if temp == nil {
			logger.Error("Error reading exclude-repositories-file ", excludeRepositoriesFilePath)
			os.Exit(-1)
		}
		logger.Debug("Successfully read in the exclude-repositories-file ", excludeRepositoriesFilePath)
		excludeRepositories = temp
	}

	// parse include repositories
	includeRepositories := []string{}
	if includeRepositoriesFilePath != "" {
		temp := scanner.ReadIgnoreFile(includeRepositoriesFilePath)
		if temp == nil {
			logger.Error("Error reading include-repositories-file ", includeRepositoriesFilePath)
			os.Exit(-1)
		}
		logger.Debug("Successfully read in the include-repositories-file ", includeRepositoriesFilePath)
		includeRepositories = temp
	}

	// Discover repositories
	num_repos_found := 0
	repositoryInfoArr := []clone.RepoInfo{}
	if mode == LOCAL {
		num_repos_found = 1
		repositoryInfo := clone.NewRepoInfo("local-org", "", "local")
		repositoryInfoArr = append(repositoryInfoArr, repositoryInfo)

	} else if mode == GH {
		repositoryInfoArr = DiscoverRepositories(mode, accessToken, organization)
		num_repos_found = len(repositoryInfoArr)

		logger.Info("Discovered ", num_repos_found, " repositories in ", organization)
	}

	// create output folder with time stamp
	timeDir := time.Now().Format("20060102_150405") // Format: YYYYMMDD_HHMMSS
	logger.Debug("Creating output folder ", timeDir)
	err := os.Mkdir(timeDir, 0777)
	if err != nil {
		log.Fatalln(err)
	}

	failedRepos := []clone.RepoInfo{}
	csvFilePaths := []string{}

	// for each repo, clone and scan
	for index, repoInfo := range repositoryInfoArr {
		clonedRepoDir := ""
		if mode == LOCAL {
			// set directory or file to local file
			clonedRepoDir = localScanFilePath
			logger.Debug("Local file scan path is ", localScanFilePath)
		} else if mode == GH {
			// check if we should include or exclude this repo
			if contains(excludeRepositories, repoInfo.RepositoryName) {
				logger.Info((index + 1), "/", len(repositoryInfoArr), " skipping ", repoInfo.RepositoryName, " as it is in the exclude list")
				continue
			}

			// check if we should include this repo
			if len(includeRepositories) > 0 && !contains(includeRepositories, repoInfo.RepositoryName) {
				logger.Info((index + 1), "/", len(repositoryInfoArr), " skipping ", repoInfo.RepositoryName, " as it is not in the include list")
				continue
			}

			// print status
			logger.Info((index + 1), "/", len(repositoryInfoArr), " cloning respository ", repoInfo.RepositoryName, "...")

			// clone repo
			clonedRepoDir = clone.CloneGithubRepoViaZip(repoInfo, accessToken)
			if clonedRepoDir == "" {
				// Failed to clone repo, save metadata for later reporting
				// TODO dump this to CSV
				failedRepos = append(failedRepos, repoInfo)
				continue
			}
		}

		// scan LOC for the directory
		logger.Info("Scanning ", clonedRepoDir, "...")
		filePaths := scanner.WalkDirectory(clonedRepoDir, ignorePatterns)
		resultsArr := []scanner.FileScanResults{}
		for _, filePath := range filePaths {
			resultsArr = append(resultsArr, scanner.ScanFile(filePath))
		}

		// Dump results by file in a csv
		outputCSV_fullPath := filepath.Join(timeDir, repoInfo.Id+".csv")
		logger.Debug("Dumping results to ", outputCSV_fullPath)
		report.OutputCSV(resultsArr, outputCSV_fullPath)
		// TODO error checking
		csvFilePaths = append(csvFilePaths, outputCSV_fullPath)
		logger.Info("Done! Results for ", repoInfo.RepositoryName, " can be found ", outputCSV_fullPath)

		if mode == LOCAL {
			// do not delete the directory
		} else if mode == GH {
			// Attempt to remove the directory and its contents
			logger.Debug("Deleting directory ", clonedRepoDir)
			err := os.RemoveAll(clonedRepoDir)
			if err != nil {
				logger.Warn("Failed to remove directory: %v", err)
			}
		}
	}

	// print failed repos
	// TODO dump to csv
	if mode == LOCAL {
		// nothing should have failed
	} else if mode == GH {
		numFailedRepos := len(failedRepos)
		if numFailedRepos > 0 {
			logger.Info(numFailedRepos, "/", num_repos_found, " failed to process. See below for a list")
			for _, failedRepo := range failedRepos {
				logger.Info(failedRepo.RepositoryName, " - ", failedRepo.Id)
			}
		} else {
			logger.Info("0 repos failed to scan.")
		}
	}

	totalLoc := 0
	if mode == LOCAL {
		reportTotals := report.ParseTotalsFromCSVs(csvFilePaths)
		totalLoc = reportTotals[0].Total
	} else if mode == GH {
		// combine csv reports
		logger.Debug("Combining results...")
		repoResults := report.ParseTotalsFromCSVs(csvFilePaths)
		combinedReportsCSVFilePath := filepath.Join(timeDir, "AAA_combined_results.csv")
		totalLoc = report.OutputCombinedCSV(repoResults, combinedReportsCSVFilePath)
		logger.Info("Total LOC results can be found ", combinedReportsCSVFilePath)
	}

	logger.Info("Total LOC for ", organization, " is ", totalLoc)
}
