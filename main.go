package main

import (
	"go-cloc/clone"
	"go-cloc/logger"
	"go-cloc/report"
	"go-cloc/scanner"
	"go-cloc/utilities"
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
	// parse CLI arguments and store them in a struct
	args := utilities.ParseArgsFromCLI()

	// Discover repositories
	logger.Info("Discovering repositories...")
	repositoryInfoArr := DiscoverRepositories(args.Mode, args.AccessToken, args.Organization)
	num_repos_found := len(repositoryInfoArr)
	logger.Info("Discovered ", num_repos_found, " repositories in ", args.Organization)

	// create output folder with time stamp
	timeDir := time.Now().Format("20060102_150405") // Format: YYYYMMDD_HHMMSS
	logger.Debug("Creating folder ", timeDir, " to store results")
	err := os.Mkdir(timeDir, 0777)
	if err != nil {
		logger.LogStackTraceAndExit(err)
	}

	failedRepos := []clone.RepoInfo{}
	csvFilePaths := []string{}

	// for each repo, clone and scan
	for index, repoInfo := range repositoryInfoArr {
		// set directory
		clonedRepoDir := ""
		logger.Debug("Setting directory for ", repoInfo.RepositoryName, " to begin scanning")
		if args.Mode == utilities.LOCAL {
			// set directory or file to local file
			clonedRepoDir = args.LocalScanFilePath
			logger.Debug("Local file scan path is ", args.LocalScanFilePath)
		} else {
			logger.Debug("Checking repo ", repoInfo.RepositoryName, " for exclusion")
			// check if we should include or exclude this repo
			if contains(args.ExcludeRepositories, repoInfo.RepositoryName) {
				logger.Info((index + 1), "/", len(repositoryInfoArr), " skipping ", repoInfo.RepositoryName, " as it is in the exclude list")
				continue
			}

			// check if we should include this repo
			if len(args.IncludeRepositories) > 0 && !contains(args.IncludeRepositories, repoInfo.RepositoryName) {
				logger.Info((index + 1), "/", len(repositoryInfoArr), " skipping ", repoInfo.RepositoryName, " as it is not in the include list")
				continue
			}

			// print status
			logger.Info((index + 1), "/", len(repositoryInfoArr), " cloning respository ", repoInfo.RepositoryName, "...")

			if args.CloneRepoUsingZip {
				logger.Debug("Cloning using zip")

				// clone repo
				zipUrl := clone.CreateZipURLGithub(repoInfo.OrganizationName, repoInfo.RepositoryName)
				clonedRepoDir = clone.DonwloadAndUnzip(zipUrl, repoInfo.RepositoryName, args.AccessToken)
				if clonedRepoDir == "" {
					// Failed to clone repo, save metadata for later reporting
					// TODO dump this to CSV
					failedRepos = append(failedRepos, repoInfo)
					continue
				}
			} else {
				logger.Debug("Cloning using git clone")
				// clone repo
				clonedRepoDir = CloneRepoMain(args.Mode, args.AccessToken, args.Organization, repoInfo)

				// Handle failed repos
				if clonedRepoDir == "" {
					// Failed to clone repo, save metadata for later reporting
					// TODO dump this to CSV
					failedRepos = append(failedRepos, repoInfo)
					continue
				}
			}
		}

		// scan LOC for the directory
		logger.Info("Scanning ", clonedRepoDir, "...")
		filePaths := scanner.WalkDirectory(clonedRepoDir, args.IgnorePatterns)
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

		if args.Mode == utilities.LOCAL {
			// do not delete the directory
		} else {
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
	if args.Mode == utilities.LOCAL {
		// nothing should have failed
	} else {
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
	if args.Mode == utilities.LOCAL {
		reportTotals := report.ParseTotalsFromCSVs(csvFilePaths)
		totalLoc = reportTotals[0].Total
	} else {
		// combine csv reports
		logger.Debug("Combining results...")
		repoResults := report.ParseTotalsFromCSVs(csvFilePaths)
		combinedReportsCSVFilePath := filepath.Join(timeDir, "AAA-combined-total-lines.csv")
		totalLoc = report.OutputCombinedCSV(repoResults, combinedReportsCSVFilePath)
		logger.Info("Total LOC results can be found ", combinedReportsCSVFilePath)
	}

	logger.Info("Total LOC for ", args.Organization, " is ", totalLoc)
}

func CloneRepoMain(mode string, accessToken string, organization string, repoInfo clone.RepoInfo) string {
	cloneRepoUrl := ""
	clonedRepoDir := ""
	if mode == utilities.GITHUB {
		cloneRepoUrl = clone.CreateCloneURLGithub(accessToken, organization, repoInfo.RepositoryName)
		clonedRepoDir = clone.CloneRepo(cloneRepoUrl, accessToken, repoInfo.RepositoryName)

	} else if mode == utilities.AZUREDEVOPS {
		cloneRepoUrl = clone.CreateCloneURLAzureDevOps(accessToken, organization, repoInfo.ProjectName, repoInfo.RepositoryName)
		// cloneRepoUrl = clone.CloneRepoAzureDevOps(accessToken, repoInfo.OrganizationName, repoInfo.ProjectName, repoInfo.RepositoryName)
		clonedRepoDir = clone.CloneRepoAzureDevOps(cloneRepoUrl, accessToken, repoInfo.RepositoryName)
	} else if mode == utilities.GITLAB {
		cloneRepoUrl = clone.CreateCloneURLGitLab(accessToken, organization, repoInfo.RepositoryName)
		clonedRepoDir = clone.CloneRepo(cloneRepoUrl, accessToken, repoInfo.RepositoryName)
	} else {
		logger.Error("Mode ", mode, " is not supported")
	}
	return clonedRepoDir
}

func DiscoverRepositories(mode string, accessToken string, organization string) []clone.RepoInfo {
	repositoryInfoArr := []clone.RepoInfo{}
	if mode == utilities.LOCAL {
		repositoryInfo := clone.NewRepoInfo("local-org", "", "local")
		repositoryInfoArr = append(repositoryInfoArr, repositoryInfo)
	} else if mode == utilities.GITHUB {
		repositoryInfoArr = clone.DiscoverReposGithub(organization, accessToken)
	} else if mode == utilities.AZUREDEVOPS {
		repositoryInfoArr = clone.DiscoverReposAzureDevOps(organization, accessToken)
	} else if mode == utilities.GITLAB {
		repositoryInfoArr = clone.DiscoverReposGitlab(organization, accessToken)
	} else if mode == utilities.BITBUCKET {
		repositoryInfoArr = clone.DiscoverReposBitbucket(organization, accessToken)
	} else {
		logger.LogStackTraceAndExit("Mode " + mode + " is not supported")
	}
	return repositoryInfoArr
}
