package devops

type RepoInfo struct {
	Id               string
	RepositoryName   string
	OrganizationName string
	ProjectName      string
	DefaultBranch    string
}

/*
* Constructor for RepoInfo
 */
func NewRepoInfo(organization string, project string, repositoryName string, defaultBranch string) RepoInfo {
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
	repoInfo.DefaultBranch = defaultBranch
	return repoInfo
}
