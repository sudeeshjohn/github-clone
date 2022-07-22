# github-clone
To clone GitHub issues with labels , milestones etc from a source repository to empty target .
This script except the target repository with out any milestones, labels and issues

Run:

1. Update information in  .env file
   ```
   export GITHUB_OAUTH_TOKEN=<github api token>
   export GITHUB_SOURCE_ORG=<source org/owner>
   export GITHUB_SOURCE_REPO=<source repository>
   export GITHUB_TARGET_ORG=<target org/owner>
   export GITHUB_TARGET_REPO=<target repository>
   export GITHUB_ENTERPRISE_URL=<github enterprise url>
   ```


set GITHUB_ENTERPRISE_URL only if you are working with an enterprise github

2. Run  source .env
3. make
4. ./github-clone

TODO:
Update assignee in the target issues