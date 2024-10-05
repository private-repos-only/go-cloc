package github

import (
	"encoding/json"
	"go-cloc/devops"
	"go-cloc/logger"
	"io"
	"net/http"
	"strconv"
	"strings"
)

// Define a struct with only the fields you care about
type item struct {
	Name string `json:"name"`
}

type repo struct {
	DefaultBranch string `json:"default_branch"`
}

func CreateCloneURLGithub(accessToken string, organization string, repoName string) string {
	return "https://oauth2:" + accessToken + "@github.com/" + organization + "/" + repoName + ".git"
}

func CreateZipURLGithub(organization string, repoName string, defaultBranch string) string {
	return "https://api.github.com/repos/" + organization + "/" + repoName + "/zipball/" + defaultBranch
}

func CreateGetDefaultBranchURLGitHub(organization string, repoName string) string {
	return "https://api.github.com/repos/" + organization + "/" + repoName
}

func CreateDiscoverURLGitHub(organization string, pageNum int, pageSize int) string {
	return "https://api.github.com/orgs/" + organization + "/repos?per_page=" + strconv.Itoa(pageSize) + "&page=" + strconv.Itoa(pageNum)
}

func DiscoverReposGithub(organization string, accessToken string) []devops.RepoInfo {
	pageSize := 100
	pageNum := 1
	repoNames := []devops.RepoInfo{}

	// pageNum -1 means there are no more pages to discover
	for pageNum != -1 {
		apiURL := CreateDiscoverURLGitHub(organization, pageNum, pageSize)
		logger.Debug("GET: " + apiURL)

		// Create a new HTTP request
		req, _ := http.NewRequest("GET", apiURL, nil)
		req.Header.Set("Authorization", "Bearer "+accessToken)

		// Perform the request using the default HTTP client
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			logger.LogStackTraceAndExit(err)
		}
		defer resp.Body.Close()

		// Read the response body
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			logger.Error("Failed to read response body: ", body)
			logger.LogStackTraceAndExit(err)
		}

		// Check if the status code is 200
		if resp.StatusCode != http.StatusOK {
			logger.Error("Response status code: ", resp.StatusCode, " expected 200")
			logger.LogStackTraceAndExit(nil)
		}

		// Parse the JSON response into a slice of Item
		var result []item
		if err := json.Unmarshal(body, &result); err != nil {
			logger.Error("Failed to parse JSON: ", err)
			logger.LogStackTraceAndExit(err)
		}

		// Print the parsed data
		logger.Debug("Default branch is: ", result)

		for _, item := range result {
			repoInfo := devops.NewRepoInfo(organization, "", item.Name)
			repoNames = append(repoNames, repoInfo)
		}

		// Get the next page URL
		link := resp.Header.Get("Link")
		logger.Debug("Link header: ", link)
		if link == "" || !strings.Contains(link, `rel="last"`) {
			pageNum = -1
		} else {
			pageNum = pageNum + 1
		}
	}

	return repoNames
}

func DiscoverDefaultBranchForRepoGithub(organization string, repoName string, accessToken string) string {
	logger.Debug("Getting default branch for ", organization, "/", repoName)

	url := CreateGetDefaultBranchURLGitHub(organization, repoName)
	logger.Debug("GET: " + url)

	// Create a new HTTP request
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+accessToken)

	// Perform the request using the default HTTP client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logger.LogStackTraceAndExit(err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error("Failed to read response body: ", body)
		logger.LogStackTraceAndExit(err)
	}

	// Check if the status code is 200
	if resp.StatusCode != http.StatusOK {
		logger.Error("Response status code: ", resp.StatusCode, " expected 200")
		logger.LogStackTraceAndExit(nil)
	}

	// Parse the JSON response into a slice of Item
	var repoResult repo
	if err := json.Unmarshal(body, &repoResult); err != nil {
		logger.Error("Failed to parse JSON: ", err)
		logger.LogStackTraceAndExit(err)
	}

	// Print the parsed data
	logger.Debug("Default branch is: ", repoResult.DefaultBranch)

	return repoResult.DefaultBranch
}
