package github

import (
	"encoding/json"
	"go-cloc/devops"
	"go-cloc/logger"
	"io"
	"net/http"
)

func CreateCloneURLGithub(accessToken string, organization string, repoName string) string {
	return "https://oauth2:" + accessToken + "@github.com/" + organization + "/" + repoName + ".git"
}

func CreateZipURLGithub(organization string, repoName string) string {
	// getUrl := "https://oauth2:" + accessToken + "@api.github.com/repos/" + organization + "/" + repoName + "/zipball"

	return "https://api.github.com/repos/" + organization + "/" + repoName + "/zipball"
}

// Define a struct with only the fields you care about
type item struct {
	Name string `json:"name"`
}

func CreateDiscoverURLGitHub(organization string) string {
	return "https://api.github.com/orgs/" + organization + "/repos?per_page=100&page=1"
}

func DiscoverReposGithub(organization string, accessToken string) []devops.RepoInfo {
	apiURL := CreateDiscoverURLGitHub(organization)
	logger.Debug("GET: " + apiURL)
	logger.Debug("Using access token: " + accessToken)

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
	logger.Debug(result)

	repoNames := []devops.RepoInfo{}
	for _, item := range result {
		repoInfo := devops.NewRepoInfo(organization, "", item.Name)
		repoNames = append(repoNames, repoInfo)
	}
	return repoNames
}
