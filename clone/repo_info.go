package clone

type RepoInfo struct {
	Id               string
	RepositoryName   string
	OrganizationName string
	ProjectName      string
}

/*
* Constructor for RepoInfo
 */
func NewRepoInfo(organization string, project string, repositoryName string) RepoInfo {
	repoInfo := RepoInfo{}
	// Github does not have a concept of projects
	if project == "" {
		repoInfo.Id = organization + "-" + repositoryName
	} else {
		repoInfo.Id = organization + "-" + project + "-" + repositoryName
	}
	repoInfo.ProjectName = project
	repoInfo.OrganizationName = organization
	repoInfo.RepositoryName = repositoryName
	return repoInfo
}
