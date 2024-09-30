package clone

import (
	"archive/zip"
	"encoding/json"
	"go-cloc/logger"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
)

// UnzipAndRename extracts and renames the top-level directory in the zip file.
func UnzipAndRename(src string, dest string, newFolderName string) error {
	// Open the zip file
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	var baseDir string

	// Iterate through the files in the archive
	for _, f := range r.File {
		// Get the file's path within the archive
		fpath := f.Name

		// Identify the base directory (i.e., the top-level folder)
		if baseDir == "" {
			parts := strings.Split(fpath, string(os.PathSeparator))
			baseDir = parts[0] // the first part is the top-level folder
		}

		// Replace the base directory name with the user-provided new name
		fpath = strings.Replace(fpath, baseDir, newFolderName, 1)

		// Create the full destination path for the file
		fpath = filepath.Join(dest, fpath)

		// Check if the file is a directory, create it if necessary
		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}

		// Create the file's directory if it doesn't exist
		if err := os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return err
		}

		// Open the file for writing
		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}
		defer outFile.Close()

		// Open the zip file entry for reading
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer rc.Close()

		// Copy the file's contents to the output file
		_, err = io.Copy(outFile, rc)
		if err != nil {
			return err
		}
	}
	return nil
}

/*
@return The resulting directory name of the repository
*/
func CloneGithubRepoViaZip(repoInfo RepoInfo, accessToken string) string {
	// extract data
	organization := repoInfo.OrganizationName
	repoName := repoInfo.RepositoryName

	// make getURL
	// getUrl := "https://api.github.com/repos/cole-gannaway-sonarsource/cloc-wrapper/zipball"
	getUrl := "https://oauth2:" + accessToken + "@api.github.com/repos/" + organization + "/" + repoName + "/zipball"
	logger.Debug("Cloning using url: ", getUrl)

	// Make API call
	resp, err := http.Get(getUrl)
	if err != nil {
		log.Fatalln(err)
	}

	// Check if the status code is 200
	if resp.StatusCode != http.StatusOK {
		logger.Error(resp.Status, " ", resp.StatusCode, " ", resp.Body)
		return ""
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	// Save the file to a local zip file
	zipFilePath := repoName + ".zip"
	err = os.WriteFile(zipFilePath, body, 0644)
	if err != nil {
		log.Fatalln(err)
	}

	// unzip the file
	// TODO check for errors here
	UnzipAndRename(zipFilePath, "", repoName)
	logger.Debug("File unzipped successfully!")

	// remove the zip file
	os.Remove(zipFilePath)

	// return the resulting directory
	return repoName
}
func CloneRepo(organization string, repoName string, accessToken string) {
	// Clone the given repository to the given directory
	url := "https://oauth2:" + accessToken + "@github.com/" + organization + "/" + repoName + ".git"
	dir := "./" + repoName // Directory where repo will be cloned

	log.Printf("Cloning %s into %s...\n", url, dir)

	// Clone repository to specified directory
	_, err := git.PlainClone(dir, false, &git.CloneOptions{
		URL:      url,
		Progress: log.Writer(), // Display progress in the log
	})

	if err != nil {
		log.Fatalf("Error cloning repository: %s", err)
	}

	logger.Debug("Repository successfully cloned!")
}

// Define a struct with only the fields you care about
type GithubAPIItem struct {
	Name string `json:"name"`
}

func DiscoverReposGithub(organization string, authToken string) []RepoInfo {
	apiURL := "https://api.github.com/orgs/" + organization + "/repos?per_page=100&page=1"

	// Create a new HTTP request
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		log.Fatalf("Failed to create HTTP request: %v", err)
	}

	// Add the Authorization header
	req.Header.Set("Authorization", "Bearer "+authToken)

	// Perform the request using the default HTTP client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Failed to fetch data from API: %v", err)
	}
	defer resp.Body.Close()

	// Check if the status code is 200
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Unexpected status code: %d, expected 200", resp.StatusCode)
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

func NowTime() string {
	now := time.Now()
	timestamp := now.Format(time.RFC3339)
	return timestamp
}
