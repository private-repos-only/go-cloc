package bitbucket

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_bitbucket_CreateDiscoverURLBitbucket(t *testing.T) {

	organization := "organization"
	pageNum := 1
	pageSize := 100
	apiURL := CreateDiscoverURLBitbucket(organization, pageNum, pageSize)
	expected := "https://api.bitbucket.org/2.0/repositories/organization?pagelen=100&page=1"
	// Assert
	assert.Equal(t, expected, apiURL)
}
func Test_bitbucket_CreateCloneURLBitbucket(t *testing.T) {

	organization := "organization"
	repository := "repository"
	accessToken := "accessToken"
	cloneUrl := CreateCloneURLBitbucket(accessToken, organization, repository)
	// Assert
	assert.Equal(t, "https://x-token-auth:accessToken@bitbucket.org/organization/repository.git", cloneUrl)
}
