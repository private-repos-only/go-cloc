package utilities

import (
	"flag"
	"go-cloc/logger"
	"go-cloc/scanner"
	"os"
)

// Modes
const (
	LOCAL       string = "Local"
	GITHUB      string = "GitHub"
	AZUREDEVOPS string = "AzureDevOps"
	GITLAB      string = "GitLab"
	BITBUCKET   string = "Bitbucket"
)

type CLIArgs struct {
	LogLevel            string
	Mode                string
	LocalScanFilePath   string
	AccessToken         string
	Organization        string
	IgnorePatterns      []string
	ExcludeRepositories []string
	IncludeRepositories []string
	CloneRepoUsingZip   bool
}

func ParseArgsFromCLI() CLIArgs {

	// mandatory arguments
	modeArg := flag.String("devops", LOCAL, "flag : <GitHub>||<AzureDevOps>||<Bitbucket>||<GitLab>||<File>")
	accessTokenArg := flag.String("accessToken", "", "Your DevOps personal access token used for discovering and downloading repositories in your organization")
	organizationArg := flag.String("organization", "", "Your DevOps organization name")
	// optional arguments
	logLevelArg := flag.String("log-level", "INFO", "Log level (DEBUG, INFO, WARN, ERROR)")
	localScanFilePathArg := flag.String("local-file-path", "", "Path to your local file or directory that you want to scan")
	ignoreFilePathArg := flag.String("ignore-file", "", "(Optional) Path to your ignore file to exclude directories and files. Please see the README.md for how to format your ignore configuration")
	excludeRepositoriesFilePathArg := flag.String("exclude-repositories-file", "", "(Optional) Path to your exclude repositories file to exclude repositories. Please see the README.md for how to format your exclude repositories configuration")
	includeRepositoriesFilePathArg := flag.String("include-repositories-file", "", "(Optional) Path to your include repositories file to include repositories. Please see the README.md for how to format your include repositories configuration")
	cloneRepoUsingZipArg := flag.Bool("clone-repo-using-zip", false, "Flag to clone repositories using zip files instead of git clone for performance improvements")

	// parse the CLI arguments
	flag.Parse()

	// dereference all CLI args to make it easier to use
	logLevel := *logLevelArg
	mode := *modeArg
	localScanFilePath := *localScanFilePathArg
	accessToken := *accessTokenArg
	organization := *organizationArg
	ignoreFilePath := *ignoreFilePathArg
	excludeRepositoriesFilePath := *excludeRepositoriesFilePathArg
	includeRepositoriesFilePath := *includeRepositoriesFilePathArg
	cloneRepoUsingZip := *cloneRepoUsingZipArg

	// set log level
	logger.SetLogLevel(logger.ConvertStringToLogLevel(logLevel))
	logger.SetOutput(os.Stdout)

	logger.Info("Setting Log Level to " + logLevel)
	logger.Info("Parsing CLI arguments")

	// print out arguments
	logger.Debug("Mode: ", mode)
	logger.Debug("clone-repo-using-zip: ", cloneRepoUsingZip)

	// validate mandatory arguments
	logger.Debug("Validating mandatory arguments")
	if mode == LOCAL {
		if localScanFilePath == "" {
			logger.Error("Mode ", mode, " requires : --local-file-path")
			os.Exit(-1)
		}
	} else {
		if organization == "" || accessToken == "" {
			logger.Error("Mode ", mode, " requires : --organization & --accessToken")
			os.Exit(-1)
		}
	}

	// validate optional arguments

	// parse ignore patterns
	ignorePatterns := []string{}
	if ignoreFilePath != "" {
		logger.Debug("Parsing ignore-file ", ignoreFilePath)
		ignorePatterns = scanner.ReadIgnoreFile(ignoreFilePath)
		logger.Debug("Successfully read in the ignore-file ", ignoreFilePath)
		logger.Debug("Ignore Patterns: ", ignorePatterns)
	}

	// parse exclude repositories
	excludeRepositories := []string{}
	if excludeRepositoriesFilePath != "" {
		logger.Debug("Parsing exclude-repositories-file ", excludeRepositoriesFilePath)
		excludeRepositories = scanner.ReadIgnoreFile(excludeRepositoriesFilePath)
		logger.Debug("Successfully read in the exclude-repositories-file ", excludeRepositoriesFilePath)
		logger.Debug("Exclude Repositories: ", excludeRepositories)
	}

	// parse include repositories
	includeRepositories := []string{}
	if includeRepositoriesFilePath != "" {
		logger.Debug("Parsing include-repositories-file ", includeRepositoriesFilePath)
		includeRepositories = scanner.ReadIgnoreFile(includeRepositoriesFilePath)
		logger.Debug("Successfully read in the include-repositories-file ", includeRepositoriesFilePath)
		logger.Debug("Include Repositories: ", includeRepositories)
	}

	args := CLIArgs{
		LogLevel:            logLevel,
		Mode:                mode,
		LocalScanFilePath:   localScanFilePath,
		AccessToken:         accessToken,
		Organization:        organization,
		IgnorePatterns:      ignorePatterns,
		ExcludeRepositories: excludeRepositories,
		IncludeRepositories: includeRepositories,
		CloneRepoUsingZip:   cloneRepoUsingZip,
	}

	return args
}
