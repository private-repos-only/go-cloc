package main

import (
	"go-cloc/azuredevops"
	"go-cloc/bitbucket"
	"go-cloc/clone"
	"go-cloc/devops"
	"go-cloc/github"
	"go-cloc/gitlab"
	"go-cloc/logger"
	"go-cloc/report"
	"go-cloc/scanner"
	"go-cloc/utilities"
	"os"
	"path/filepath"
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
	initialNumReposFound := len(repositoryInfoArr)
	logger.Info("Discovered ", initialNumReposFound, " repositories in ", args.Organization)

	// Filter repositories
	logger.Info("Including / Excluding repositories...")
	fitleredRepoInfoArr := []devops.RepoInfo{}
	for _, repoInfo := range repositoryInfoArr {
		logger.Debug("Checking repo ", repoInfo.RepositoryName, " for exclusion")
		// check if we should include or exclude this repo
		if contains(args.ExcludeRepositories, repoInfo.RepositoryName) {
			logger.Debug("Excluding ", repoInfo.RepositoryName, " as it is in the exclude list")
			continue
		}

		// check if we should include this repo
		if len(args.IncludeRepositories) > 0 && !contains(args.IncludeRepositories, repoInfo.RepositoryName) {
			logger.Debug("Excluding ", repoInfo.RepositoryName, " as it is NOT in the include list")
			continue
		}

		logger.Debug("Including ", repoInfo.RepositoryName)
		fitleredRepoInfoArr = append(fitleredRepoInfoArr, repoInfo)
	}
	numRepos := len(fitleredRepoInfoArr)

	// create output folder
	if args.DumpCSVs {
		logger.Debug("Creating folder ", args.ResultsDirectoryPath, " to store results")
		err := os.Mkdir(args.ResultsDirectoryPath, 0777)
		if err != nil {
			logger.LogStackTraceAndExit(err)
		}
	}

	failedRepos := []devops.RepoInfo{}
	allRepoResults := []report.RepoTotal{}
	// for each repo, clone and scan
	for index, repoInfo := range fitleredRepoInfoArr {
		// set directory
		clonedRepoDir := ""
		logger.Debug("Setting directory for ", repoInfo.RepositoryName, " to begin scanning")
		if args.Mode == utilities.LOCAL {
			// set directory or file to local file
			clonedRepoDir = args.LocalScanFilePath
			logger.Debug("Local file scan path is ", args.LocalScanFilePath)
		} else {

			// print status
			logger.Info((index + 1), "/", len(fitleredRepoInfoArr), " cloning respository ", repoInfo.RepositoryName, "...")

			// TODO: add support for cloning using zip for more platforms
			if args.CloneRepoUsingZip && args.Mode == utilities.GITHUB {
				logger.Debug("Cloning using zip")

				// clone repo
				defaultBranch := github.DiscoverDefaultBranchForRepoGithub(repoInfo.OrganizationName, repoInfo.RepositoryName, args.AccessToken)
				zipUrl := github.CreateZipURLGithub(repoInfo.OrganizationName, repoInfo.RepositoryName, defaultBranch)
				clonedRepoDir = clone.DonwloadAndUnzip(zipUrl, repoInfo.RepositoryName, args.AccessToken)
				if clonedRepoDir == "" {
					// Failed to clone repo, save metadata for later reporting
					failedRepos = append(failedRepos, repoInfo)
					logger.Error("Failed to clone repo ", repoInfo.RepositoryName)
					// skip to the next repo
					continue
				}
			} else {
				logger.Debug("Cloning using git clone")
				// clone repo
				clonedRepoDir = CloneRepo(args.Mode, args.AccessToken, args.Organization, repoInfo)

				// Handle failed repos
				if clonedRepoDir == "" {
					// Failed to clone repo, save metadata for later reporting
					failedRepos = append(failedRepos, repoInfo)
					logger.Error("Failed to clone repo ", repoInfo.RepositoryName)
					// skip to the next repo
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
		if args.DumpCSVs {
			outputCsvFilePath := filepath.Join(args.ResultsDirectoryPath, repoInfo.Id+".csv")
			logger.Debug("Dumping results to ", outputCsvFilePath)
			codeLineCount := report.OutputCSV(resultsArr, outputCsvFilePath)
			allRepoResults = append(allRepoResults, report.RepoTotal{RepositoryName: repoInfo.RepositoryName, Total: codeLineCount})

			// TODO error checking
			logger.Info("Done! Results for ", repoInfo.RepositoryName, " can be found ", outputCsvFilePath)
		}

		// clean up cloned repo after scan completes
		if args.Mode == utilities.LOCAL {
			// do not delete the directory if we are scanning a local file or directory
		} else {
			// delete the cloned repo directory after scanning
			logger.Debug("Deleting directory ", clonedRepoDir)
			err := os.RemoveAll(clonedRepoDir)
			if err != nil {
				logger.Error("Failed to remove directory: ", clonedRepoDir)
			}
		}
	}

	// print failed repos
	numFailedRepos := len(failedRepos)
	if numFailedRepos > 0 {
		logger.Info(numFailedRepos, "/", numRepos, " failed to process. See below for a list")
		for _, failedRepo := range failedRepos {
			logger.Info(failedRepo.RepositoryName, " - ", failedRepo.Id)
		}
	} else {
		logger.Info("0 repos failed to scan.")
	}

	// sum total LOC for all repos
	totalLoc := 0
	for _, repoResult := range allRepoResults {
		totalLoc += repoResult.Total
	}

	// dump combined csv reports
	if args.DumpCSVs {
		logger.Debug("Combining results...")
		combinedReportsCSVFilePath := filepath.Join(args.ResultsDirectoryPath, "AAA-combined-total-lines.csv")
		totalLoc = report.OutputCombinedCSV(allRepoResults, combinedReportsCSVFilePath)
		logger.Info("Total LOC results can be found ", combinedReportsCSVFilePath)
	}

	logger.Info("Total LOC for ", args.Organization, " is ", totalLoc)

	// return total LOC as exit code for external use
	os.Exit(totalLoc)
}

func CloneRepo(mode string, accessToken string, organization string, repoInfo devops.RepoInfo) string {
	cloneRepoUrl := ""
	clonedRepoDir := ""
	if mode == utilities.GITHUB {
		cloneRepoUrl = github.CreateCloneURLGithub(accessToken, organization, repoInfo.RepositoryName)
		clonedRepoDir = clone.CloneRepo(cloneRepoUrl, accessToken, repoInfo.RepositoryName)

	} else if mode == utilities.AZUREDEVOPS {
		cloneRepoUrl = azuredevops.CreateCloneURLAzureDevOps(accessToken, organization, repoInfo.ProjectName, repoInfo.RepositoryName)
		clonedRepoDir = clone.CloneRepoAzureDevOps(cloneRepoUrl, accessToken, repoInfo.RepositoryName)
	} else if mode == utilities.GITLAB {
		cloneRepoUrl = gitlab.CreateCloneURLGitLab(accessToken, organization, repoInfo.RepositoryName)
		clonedRepoDir = clone.CloneRepo(cloneRepoUrl, accessToken, repoInfo.RepositoryName)
	} else if mode == utilities.BITBUCKET {
		cloneRepoUrl = bitbucket.CreateCloneURLBitbucket(accessToken, organization, repoInfo.RepositoryName)
		clonedRepoDir = clone.CloneRepo(cloneRepoUrl, accessToken, repoInfo.RepositoryName)
	} else {
		logger.Error("Mode ", mode, " is not supported")
	}
	return clonedRepoDir
}

func DiscoverRepositories(mode string, accessToken string, organization string) []devops.RepoInfo {
	repositoryInfoArr := []devops.RepoInfo{}
	if mode == utilities.LOCAL {
		repositoryInfo := devops.NewRepoInfo("local-org", "", "local", "")
		repositoryInfoArr = append(repositoryInfoArr, repositoryInfo)
	} else if mode == utilities.GITHUB {
		repositoryInfoArr = github.DiscoverReposGithub(organization, accessToken)
	} else if mode == utilities.AZUREDEVOPS {
		repositoryInfoArr = azuredevops.DiscoverReposAzureDevOps(organization, accessToken)
	} else if mode == utilities.GITLAB {
		repositoryInfoArr = gitlab.DiscoverReposGitlab(organization, accessToken)
	} else if mode == utilities.BITBUCKET {
		repositoryInfoArr = bitbucket.DiscoverReposBitbucket(organization, accessToken)
	} else {
		logger.LogStackTraceAndExit("Mode " + mode + " is not supported")
	}
	return repositoryInfoArr
}
