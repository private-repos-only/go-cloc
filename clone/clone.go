package clone

import (
	"archive/zip"
	"go-cloc/logger"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5"
)

// UnzipAndRename extracts and renames the top-level directory in the zip file.
func UnzipAndRename(src string, dest string, newFolderName string) error {
	// Open the zip file
	r, err := zip.OpenReader(src)
	if err != nil {
		logger.Error("Error opening zip file: ", err)
		return err
	}
	defer r.Close()

	var baseDir string

	// Iterate through the files in the archive
	for _, f := range r.File {
		// Get the file's path within the archive
		fpath := f.Name

		// Normalize the path to use forward slashes
		fpath = strings.ReplaceAll(fpath, "\\", "/")

		// Identify the base directory (i.e., the top-level folder)
		if baseDir == "" {
			parts := strings.Split(fpath, "/")
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
			logger.Error("Error creating directory: ", err)
			return err
		}

		// Open the file for writing
		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			logger.Error("Error opening file: ", err)
			return err
		}
		defer outFile.Close()

		// Open the zip file entry for reading
		rc, err := f.Open()
		if err != nil {
			logger.Error("Error opening zip file entry: ", err)
			return err
		}
		defer rc.Close()

		// Copy the file's contents to the output file
		_, err = io.Copy(outFile, rc)
		if err != nil {
			logger.Error("Error copying file contents: ", err)
			return err
		}
	}
	return nil
}

/*
@return The resulting directory name of the repository
*/
func DonwloadAndUnzip(getUrl string, repoName string, accessToken string) string {

	logger.Debug("Cloning using url: ", getUrl)

	// Make API call
	req, _ := http.NewRequest("GET", getUrl, nil)
	req.Header.Set("Authorization", "Bearer "+accessToken)

	// Perform the request using the default HTTP client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logger.LogStackTraceAndExit(err)
	}
	defer resp.Body.Close()

	// Check if the status code is 200
	if resp.StatusCode != http.StatusOK {
		logger.Error(resp.Status, " ", resp.StatusCode, " ", resp.Body)
		return ""
	}

	// Check if the Content-Type is application/zip, if not, return
	if resp.Header.Get("Content-Type") != "application/zip" {
		logger.Error("Unexpected Content-Type: ", resp.Header.Get("Content-Type"))
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
	err = UnzipAndRename(zipFilePath, "", repoName)
	if err != nil {
		logger.Error("Error unzipping file: ", err)
		logger.LogStackTraceAndExit(err)
	}
	logger.Debug("File unzipped successfully!")

	// remove the zip file
	os.Remove(zipFilePath)

	// return the resulting directory
	return repoName
}

func CloneRepo(url string, accessToken string, repoName string) string {
	// Clone the given repository to the given directory
	dir := "./" + repoName // Directory where repo will be cloned

	logger.Debug("Cloning url: ", url, " into directory: ", dir)
	// .Printf("Cloning %s into %s...\n", url, dir)

	// Clone repository to specified directory with authentication
	_, err := git.PlainClone(dir, false, &git.CloneOptions{
		URL:          url,
		SingleBranch: true,
		Depth:        1,
	})

	// Check to see if there was an error cloning the repo
	if err != nil {
		logger.Error("Error cloning repository: ", repoName, " : ", err)
		return ""
	}

	logger.Debug("Repository successfully cloned!")
	return repoName
}
