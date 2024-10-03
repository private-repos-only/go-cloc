package bitbucket

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_clone_CreateDiscoverURLBitbucket(t *testing.T) {

	organization := "organization"
	apiURL := CreateDiscoverURLBitbucket(organization)
	expected := "https://api.bitbucket.org/2.0/repositories/organization"
	// Assert
	assert.Equal(t, expected, apiURL)
}
func Test_clone_CreateCloneURLBitbucket(t *testing.T) {

	organization := "organization"
	repository := "repository"
	accessToken := "accessToken"
	cloneUrl := CreateCloneURLBitbucket(accessToken, organization, repository)
	// Assert
	assert.Equal(t, "https://x-token-auth:accessToken@bitbucket.org/organization/repository.git", cloneUrl)
}
