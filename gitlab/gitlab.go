package gitlab

import (
	"encoding/json"
	"go-cloc/devops"
	"go-cloc/logger"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func CreateCloneURLGitLab(accessToken string, organization string, respository string) string {
	// Create the URL
	return "https://oauth2:" + accessToken + "@gitlab.com/" + organization + "/" + respository + ".git"
}

// Discovers projects in a GitLab organization
func CreateDiscoverURLGitLab(accessToken string, organization string, pageNum int, pageSize int) string {
	return "https://" + accessToken + "@gitlab.com/api/v4/groups/" + organization + "/projects?per_page=" + strconv.Itoa(pageSize) + "&page=" + strconv.Itoa(pageNum)
}

// Define the nested struct types
type item struct {
	Name string `json:"name"`
}

func DiscoverReposGitlab(organization string, accessToken string) []devops.RepoInfo {
	pageSize := 100
	pageNum := 1
	repoNames := []devops.RepoInfo{}
	// pageNum -1 means there are no more pages to discover
	for pageNum != -1 {
		apiURL := CreateDiscoverURLGitLab(accessToken, organization, pageNum, pageSize)
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
		for _, item := range result {
			repoInfo := devops.NewRepoInfo(organization, "", item.Name)
			repoNames = append(repoNames, repoInfo)
		}

		// Get the next page URL
		link := resp.Header.Get("Link")
		logger.Debug("Link header: ", link)
		logger.Debug(strings.Contains(link, `rel="last"`))
		// If there is no next page, stop the loop
		if link == "" || !strings.Contains(link, `rel="next"`) {
			pageNum = -1
		} else {
			pageNum = pageNum + 1
		}
	}

	return repoNames
}
