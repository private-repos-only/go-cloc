package clone

import (
	"go-cloc/logger"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/protocol/packp/capability"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
)

func CloneRepoAzureDevOps(url string, accessToken string, repoName string) string {
	// Clone the given repository to the given directory
	dir := "./" + repoName // Directory where repo will be cloned

	logger.Debug("Cloning url: ", url, " into directory: ", dir)

	// New commits and pushes against a remote worked without any issues.
	transport.UnsupportedCapabilities = []capability.Capability{
		capability.ThinPack,
	}

	_, err := git.PlainClone(dir, false, &git.CloneOptions{
		Auth: &http.BasicAuth{
			Username: "",
			Password: accessToken,
		},
		URL: url,
	})
	if err != nil {
		logger.Error("Failed to clone repo: ", repoName)
		logger.LogStackTraceAndExit(err)
	}

	return repoName
}
