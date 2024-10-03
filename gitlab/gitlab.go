package gitlab

import (
	"encoding/json"
	"go-cloc/devops"
	"go-cloc/logger"
	"io"
	"log"
	"net/http"
)

func CreateCloneURLGitLab(accessToken string, organization string, respository string) string {
	// Create the URL
	return "https://oauth2:" + accessToken + "@gitlab.com/" + organization + "/" + respository + ".git"
}

// Discovers projects in a GitLab organization
func CreateDiscoverURLGitLab(accessToken string, organization string) string {
	return "https://" + accessToken + "@gitlab.com/api/v4/groups/" + organization + "/projects?include_subgroups=true"
}

// Define the nested struct types
type item struct {
	Name string `json:"name"`
}

func DiscoverReposGitlab(organization string, accessToken string) []devops.RepoInfo {
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
	var result []item
	if err := json.Unmarshal(body, &result); err != nil {
		log.Fatalf("Failed to parse JSON: %v", err)
	}

	// Print the parsed data
	repoNames := []devops.RepoInfo{}
	for _, item := range result {
		repoInfo := devops.NewRepoInfo(organization, "", item.Name)
		repoNames = append(repoNames, repoInfo)
	}
	return repoNames
}
