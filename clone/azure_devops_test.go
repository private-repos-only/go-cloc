package clone

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_clone_CreateCloneURL(t *testing.T) {
	accessToken := "abcdefg"
	organization := "organization"
	projectName := "project"
	repoName := "repo"
	azdoCloneURL := CreateCloneURLAzureDevOps(accessToken, organization, projectName, repoName)

	// Assert
	assert.Equal(t, "https://abcdefg@dev.azure.com/organization/project/_git/repo", azdoCloneURL)
}
