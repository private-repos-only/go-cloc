package clone

import (
	"encoding/json"
	"go-cloc/logger"
	"io"
	"log"
	"net/http"
)

// https://gitlab.com/cole-gannaway-sonarsource/php-vulnerabilities.git
func CreateCloneURLGitLab(accessToken string, organization string, respository string) string {
	// Create the URL
	return "https://oauth2:" + accessToken + "@gitlab.com/" + organization + "/" + respository + ".git"
}

// Discovers projects in a GitLab organization
func CreateDiscoverURLGitLab(accessToken string, organization string) string {
	return "https://" + accessToken + "@gitlab.com/api/v4/groups/" + organization + "/projects?include_subgroups=true"
	// https://gitlab.com/api/v4/groups/cole-gannaway-sonarsource/projects?include_subgroups=true

}
func DiscoverReposGitlab(organization string, accessToken string) []RepoInfo {
	apiURL := CreateDiscoverURLGitLab(accessToken, organization)
	logger.Debug("Discovering repos using url: ", apiURL)

	// Create a new HTTP request
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		log.Fatalf("Failed to create HTTP request: %v", err)
	}

	// Add the Authorization header
	req.Header.Set("Authorization", "Bearer "+accessToken)

	// Perform the request using the default HTTP client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Failed to fetch data from API: %v", err)
	}
	defer resp.Body.Close()

	// Check if the status code is 200
	if resp.StatusCode != http.StatusOK {
		logger.Error("Failed to fetch data from API: ", resp.Status)
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to read response body: %v", err)
	}

	// Parse the JSON response into a slice of Item
	var result []GithubAPIItem
	if err := json.Unmarshal(body, &result); err != nil {
		log.Fatalf("Failed to parse JSON: %v", err)
	}

	// Print the parsed data
	repoNames := []RepoInfo{}
	for _, item := range result {
		repoInfo := NewRepoInfo(organization, "", item.Name)
		repoNames = append(repoNames, repoInfo)
	}
	return repoNames
}

// func DiscoverReposBitbucket(organization string, accessToken string) []RepoInfo {
// // repositories/%s/%s/?fields=size

// }
