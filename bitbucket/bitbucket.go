package bitbucket

import (
	"encoding/json"
	"go-cloc/devops"
	"go-cloc/logger"
	"io"
	"log"
	"net/http"
	"strconv"
)

// Define the nested struct types
type item struct {
	Name    string  `json:"name"`
	Project project `json:"project"`
}
type project struct {
	Name string `json:"name"`
}

type response struct {
	Value []item `json:"values"`
	// Next is string or null
	Next *string `json:"next"`
}

func CreateCloneURLBitbucket(accessToken string, organization string, respository string) string {
	// Create the URL
	return "https://x-token-auth:" + accessToken + "@bitbucket.org/" + organization + "/" + respository + ".git"
}

func CreateDiscoverURLBitbucket(organization string, pageNum int, pageSize int) string {
	return "https://api.bitbucket.org/2.0/repositories/" + organization + "?pagelen=" + strconv.Itoa(pageSize) + "&page=" + strconv.Itoa(pageNum)
}
func DiscoverReposBitbucket(organization string, accessToken string) []devops.RepoInfo {
	pageSize := 100
	pageNum := 1
	repoNames := []devops.RepoInfo{}

	for pageNum != -1 {
		apiURL := CreateDiscoverURLBitbucket(organization, pageNum, pageSize)
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

		for _, item := range r.Value {
			repoInfo := devops.NewRepoInfo(organization, item.Project.Name, item.Name, "")
			repoNames = append(repoNames, repoInfo)
		}
		// If there is no next page, stop the loop
		if r.Next == nil {
			pageNum = -1
		} else {
			pageNum = pageNum + 1
		}
	}

	return repoNames
}
