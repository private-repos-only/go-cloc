# Go-Cloc

## Overview

This tool simplifies the process of obtaining an accurate Lines of Code (LOC) count for an organization's DevOps platform. It can automatically discover repositories and calculate the total LOC with a single executable.

Please download the appropriate [artifact]() for your platform.

Simply run the below command to discover all repositories in your **DevOps Organization**.
```sh
./go-cloc --devops <DevOpsPlatform>  --organization <YourOrganizationName>  --accessToken <YourPersonalAccessToken>
```
This will output the total Lines of Code (LOC) count for the entire organization. See example below.
```
2024/09/29 17:37:04 [INFO] Discovering repositories in  MyExampleOrganization
2024/09/29 17:37:04 [INFO] Discovered  50  repositories in  MyExampleOrganization
2024/09/29 17:37:04 [INFO] 1 / 50  cloning respository  example-repo ...
2024/09/29 17:37:05 [INFO] Scanning  example-repo ...
2024/09/29 17:37:05 [INFO] Done! Results for  example-repo  can be found  MyExampleOrganization-example-repo.csv
...
...
...
2024/09/29 17:37:05 [INFO] 0 repos failed to scan.
2024/09/29 17:37:05 [INFO] Total LOC results can be found  AAA-combined-total-lines.csv
2024/09/29 17:37:05 [INFO] Total LOC for  MyExampleOrganization  is  23005
```

## Requirements
1. An **Access Token** for your appropriate DevOps platform (GitHub, Azure DevOps, GitLab, or Bitbucket) with **read** access for each of the repositories within the organization.

## Options
```sh
prompt> ./go-cloc --help
```
-  `-accessToken`
       Your DevOps personal access token used for discovering and downloading repositories in your organization
-  `-clone-repo-using-zip`
       (Optional) Flag to clone repositories using zip files instead of git clone for faster downloads. Default is false.
-  `-devops`
       flag : <GitHub>||<AzureDevOps>||<Bitbucket>||<GitLab>||<File> (default "Local")
-  `-dump-csvs`
       (Optional) Flag to output CSV files. Default is true, but can be set to false to disable file dumps (default true)
-  `-exclude-repositories-file`
       (Optional) Path to your exclude repositories file to exclude repositories. Please see the README.md for how to format your exclude repositories configuration
-  `-ignore-file`
       (Optional) Path to your ignore file to exclude directories and files. Please see the README.md for how to format your ignore configuration
-  `-include-repositories-file`
       (Optional) Path to your include repositories file to include repositories. Please see the README.md for how to format your include repositories configuration
-  `-local-file-path`
       Path to your local file or directory that you want to scan
-  `-log-level`
       Log level (DEBUG, INFO, WARN, ERROR) (default "INFO")
-  `-organization`
       Your DevOps organization name
-  `-results-directory-path`
       (Optional) Path to a new directory for storing the results. By default the tool will create one

## Examples
Github
```sh
prompt> ./go-cloc --devops GitHub --organization MyExampleOrganization --accessToken abcdefg1234 
```
Local
```sh
prompt> ./go-cloc main.js 
```
## Extensibility
The tool will return an exit code of the total lines of code (LOC) count if successful, for example `103230`. If it fails, it will return an exit code of `-1`.This allows for easy integration with scripts or other 3rd party tools.

## Ignore Files

The ignore file is a simple text file used to exclude certain directories and files from processing. You can use a wildcard (`*`) to match patterns, similar to regular expressions. However, you can only use one `*` wildcard at a time. Make sure to place your ignore patterns in the ignore file, one per line, to apply them effectively.

This same configuration format applies to **exclude** or **include** repositories using the **devops** flag. Note: if using the `--devops` flag, these patterns will apply to all repositories.

- To ignore all files in a specific directory:

```sh
/path/to/directory/*
```

- To ignore all `.log` and `.js` files:
```sh
*.log
*.js
```

* Combined examples
```sh
# Local scan with ignoring certain files or directoreis
$ ./go-cloc src/ --ignore-file ignore.txt

# DevOps scan ignoring certain repositores 
$ ./go-cloc --devops GitHub \
      --organization MyExampleOrganization \
      --accessToken abcdefg1234 \
      --exclude-repositories-file github_repos_to_ignore.txt

# DevOps scan only including certain repositories
$ ./go-cloc --devops GitHub \
      --organization MyExampleOrganization \
      --accessToken abcdefg1234 \
      --include-repositories-file github_repos_to_include.txt
```

## Personal Access Tokens

Personal Access Tokens (PATs) are used to authenticate and authorize access to your DevOps platform. They are necessary for the tool to discover and clone repositories within your organization. Below are the steps to generate a PAT for different DevOps platforms:

### GitHub
1. Navigate to [GitHub Settings](https://github.com/settings/tokens).
2. Click on **Generate new token**.
3. Under **Select Scopes**, select **repo**.
5. Click **Generate token** and copy the token for use.

### Azure DevOps
1. Navigate to [Azure DevOps](https://dev.azure.com).
2. Click on **User Settings** and select **Personal Access Token**.
3. Click on **New Token**.
4. Set the name, organization, and scopes for the token.
5. Ensure that **Code -> Read** is selected as a scope.
6. Click **Create** and copy the token for use.

### GitLab
1. Navigate to [GitLab](https://gitlab.com).
2. Select your **Organization**.
3. Click on **Settings** in the left sidebar and select **Access Tokens**.
4. Provide a name and expiration date for the token.
5. Select the scopes `read_api` and `read_repository` to grant the necessary permissions.
6. Click **Create personal access token** and copy the token for use.

### Bitbucket
1. Navigate to [Bitbucket](https://bitbucket.org).
2. Select your **Organization**.
3. In the left sidebar and select **Access Tokens**.
4. Provide a name and expiration date for the token.
5. Select the scopes **Repository** to **Read**.to grant the necessary permissions.
6. Click **Create** and copy the token for use.
