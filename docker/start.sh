#!/bin/bash
echo "Current User; $(whoami)"
echo "Environment Variables:"
echo 'PAT: ' $PAT
echo 'REPO_OWNER: ' $REPO_OWNER
echo 'REPO_NAME: ' $REPO_NAME
# Use the GitHub API to create a self-hosted runner token
RUNNER_TOKEN=$(curl -X POST -H "Authorization: token ${PAT}" \
  -H "Accept: application/vnd.github.v3+json" \
  "https://api.github.com/repos/${REPO_OWNER}/${REPO_NAME}/actions/runners/registration-token" | jq .token --raw-output)

# If there is no runner, then configure one
if [ ! -f ".runner" ]; then
  ./config.sh --url https://github.com/${REPO_OWNER}/${REPO_NAME} --token ${RUNNER_TOKEN}
  touch .runner # Save a marker file indicating that the runner is configured
  ./run.sh
else # else if there is already a runner, then just start the runner
  ./run.sh
fi