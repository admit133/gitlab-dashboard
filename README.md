# Gitlab Dashboard

Gitlab Dashboard allow to manage GitLab environments. 

# Features

* Easy deploy any branch on any projects in an environment.
* Deploy history
* Deploy a specific branch on all project in an environment (i.e. deploy master on all projects)
* Redeploy current branch
* OAuth with Gitlab Server
* Quick link on logs/jobs/pipelines/projects etc.

List of environments:
![Screenshot 2020-08-21 at 10 18 29](https://user-images.githubusercontent.com/2131624/90863533-df0d9e80-e397-11ea-909e-7206f20f7fa0.png)

List of projects in an environment and last deployed branches: 
![Screenshot 2020-08-21 at 10 18 52](https://user-images.githubusercontent.com/2131624/90863527-dcab4480-e397-11ea-97e0-25f4a4fabbd7.png)

# Settings

* `GITLAB_BASE_URL` (default: `https://gitlab.com`) - URL on GitLab server
* `GITLAB_TOKEN` - GitLab token for API request (should have access to deploy projects)
* `GITLAB_PROJECT_IDS` - IDs of Gitlab Project (only these projects you will see on Dashboard)
* `ENVIRONMENT_UPDATE_DURATION` (default: `1m`) - Frequency of updating environment list
* `OAUTH_ENABLED` (default: `0`) - Enable Gitlab OAuth application (you should create an application in GitLab and specify `GITLAB_APP_ID` and `GITLAB_APP_SECRET`)
* `GITLAB_APP_ID` - App ID for OAuth
* `GITLAB_APP_SECRET` - App Secret for OAuth
* `PROTECTED_ENVIRONMENTS` - List of protected environments (these environments will be hidden on Dashboard)