package clone

import (
	"archive/zip"
	"go-cloc/logger"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5"
)

// Unzip extracts the contents of the zip file to a folder with the same name as the zip file.
func Unzip(zipFilePath string) error {
	// Open the zip file
	r, err := zip.OpenReader(zipFilePath)
	if err != nil {
		logger.Error("Error opening zip file: ", err)
		return err
	}
	defer r.Close()

	// Create the destination directory based on the zip file name
	dest := strings.TrimSuffix(zipFilePath, filepath.Ext(zipFilePath))
	if err := os.MkdirAll(dest, os.ModePerm); err != nil {
		logger.Error("Error creating destination directory: ", err)
		return err
	}

	// Iterate through the files in the archive
	for _, f := range r.File {
		// Get the file's path within the archive
		fpath := f.Name

		// Normalize the path to use forward slashes
		fpath = strings.ReplaceAll(fpath, "\\", "/")

		// Strip the top-level directory
		parts := strings.SplitN(fpath, "/", 2)
		if len(parts) > 1 {
			fpath = parts[1]
		} else {
			fpath = parts[0]
		}

		// Create the full path for the destination file
		fpath = filepath.Join(dest, fpath)

		if f.FileInfo().IsDir() {
			// Create directory
			if err := os.MkdirAll(fpath, os.ModePerm); err != nil {
				logger.Error("Error creating directory: ", err)
				return err
			}
		} else {
			// Create file
			if err := os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
				logger.Error("Error creating directory for file: ", err)
				return err
			}

			outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				logger.Error("Error opening file for writing: ", err)
				return err
			}

			rc, err := f.Open()
			if err != nil {
				logger.Error("Error opening file in zip: ", err)
				return err
			}

			_, err = io.Copy(outFile, rc)
			if err != nil {
				logger.Error("Error copying file content: ", err)
				return err
			}

			outFile.Close()
			rc.Close()
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
	contentType := resp.Header.Get("Content-Type")
	if !strings.ContainsAny(contentType, "application/zip") || !strings.ContainsAny(contentType, "application/octet-stream") {
		logger.Error("Unexpected Content-Type: ", resp.Header.Get("Content-Type"))
		return ""
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.LogStackTraceAndExit(err)
	}

	// Save the file to a local zip file
	zipFilePath := repoName + ".zip"
	err = os.WriteFile(zipFilePath, body, 0644)
	if err != nil {
		logger.LogStackTraceAndExit(err)
	}

	// unzip the file
	err = Unzip(zipFilePath)

	// remove the zip file regardless of the result
	os.Remove(zipFilePath)

	// Check if there was an error unzipping the file
	if err != nil {
		logger.Error("Error unzipping file: ", err)
		return ""
	} else {
		logger.Debug("File unzipped successfully!")
	}

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
