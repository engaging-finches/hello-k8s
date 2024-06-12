# hello-k8s

## Description
**MVP: CustomResourceDefinition (CRD), WebHooks, and associated controller to create a Deployment for Pods to run GitHub Actions self-hosted runners.**

The goal of this project is to allow users to easily manage self-hosted GitHub Actions runners in Kubernetes.

## Running Locally

### Configure commit hooks

**Generate PAT for Husky commit-msg**
PAT requires read access to issues in the repo that you're using.  Save PAT to a variable called `GITHUB_TOKEN` and store this in a `.env` file in the root directory of your repo.

**Install packages from package.json**
`npm install`

