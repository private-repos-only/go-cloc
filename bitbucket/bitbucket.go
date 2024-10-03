package bitbucket

import (
	"encoding/json"
	"go-cloc/devops"
	"go-cloc/logger"
	"io"
	"log"
	"net/http"
)

// Define the nested struct types
type item struct {
	Name string `json:"name"`
}

type response struct {
	Value []item `json:"values"`
}

func CreateCloneURLBitbucket(accessToken string, organization string, respository string) string {
	// Create the URL
	return "https://x-token-auth:" + accessToken + "@bitbucket.org/" + organization + "/" + respository + ".git"
}

func CreateDiscoverURLBitbucket(organization string) string {
	return "https://api.bitbucket.org/2.0/repositories/" + organization
}
func DiscoverReposBitbucket(organization string, accessToken string) []devops.RepoInfo {
	apiURL := CreateDiscoverURLBitbucket(organization)
	logger.Debug("Discovering repos using url: ", apiURL)

	// Create a new HTTP request
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		log.Fatalf("Failed to create HTTP request: %v", err)
	}

	// Set basic auth
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
		logger.Error("Unexpected status code: ", resp.StatusCode, ", expected 200")
		logger.Error("Response: ", resp.Status)
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to read response body: %v", err)
	}

	// Unmarshal the JSON data into the Response struct
	var r response
	err = json.Unmarshal([]byte(body), &r)
	if err != nil {
		log.Fatalf("Error unmarshalling JSON: %v", err)
	}

	logger.Debug("Response: ", r)

	repoNames := []devops.RepoInfo{}
	// Access the nested Name field
	for _, item := range r.Value {
		projectName := item.Name
		logger.Debug("Project Name:", projectName)
		// 'https://dev.azure.com/{organization}/{project_name}/_apis/git/repositories?api-version=7.0'
		apiURL := "https://dev.azure.com/" + organization + "/" + projectName + "/_apis/git/repositories?api-version=7.0"
		// Create a new HTTP request
		req, err := http.NewRequest("GET", apiURL, nil)
		if err != nil {
			log.Fatalf("Failed to create HTTP request: %v", err)
		}

		// Set basic auth
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
			logger.Error("Unexpected status code: ", resp.StatusCode, ", expected 200")
			logger.Error("Response: ", resp.Status)
		}

		// Read the response body
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatalf("Failed to read response body: %v", err)
		}

		// Unmarshal the JSON data into the Response struct
		var r response
		err = json.Unmarshal([]byte(body), &r)
		if err != nil {
			log.Fatalf("Error unmarshalling JSON: %v", err)
		}
		for _, item := range r.Value {
			repoName := item.Name
			repoInfo := devops.NewRepoInfo(organization, projectName, repoName)
			repoNames = append(repoNames, repoInfo)
		}
	}

	return repoNames
}
