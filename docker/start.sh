#!/bin/bash

RUNNER_TOKEN=$(curl -X POST -H "Authorization: token ${PAT}" \
  -H "Accept: application/vnd.github.v3+json" \
  "https://api.github.com/repos/${REPO_OWNER}/${REPO_NAME}/actions/runners/registration-token" | jq .token --raw-output)

./config.sh --url https://github.com/${REPO_OWNER}/${REPO_NAME} --token ${RUNNER_TOKEN}

cleanup() {
    echo "Removing runner..."
    ./config.sh remove --unattended --token ${RUNNER_TOKEN}
}

trap 'cleanup; exit 130' INT
trap 'cleanup; exit 143' TERM

./run.sh & wait $!