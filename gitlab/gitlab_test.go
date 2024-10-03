package gitlab

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_clone_CreateCloneURLGitLab(t *testing.T) {

	organization := "organization"
	repository := "respository"
	accessToken := "accesstoken"

	cloneUrl := CreateCloneURLGitLab(accessToken, organization, repository)
	// Assert
	assert.Equal(t, "https://oauth2:accesstoken@gitlab.com/organization/respository.git", cloneUrl)

}
