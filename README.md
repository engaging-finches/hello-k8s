[![Build and Test Go Binary](https://github.com/engaging-finches/hello-k8s/actions/workflows/build-test.yaml/badge.svg)](https://github.com/engaging-finches/hello-k8s/actions/workflows/build-test.yaml)

[![Verify Lint and Format](https://github.com/engaging-finches/hello-k8s/actions/workflows/verify-lint-fmt.yaml/badge.svg)](https://github.com/engaging-finches/hello-k8s/actions/workflows/verify-lint-fmt.yaml)

[![Verify Pre Commit Standards](https://github.com/engaging-finches/hello-k8s/actions/workflows/pre-commit-ci.yaml/badge.svg)](https://github.com/engaging-finches/hello-k8s/actions/workflows/pre-commit-ci.yaml)

[![CodeQL](https://github.com/engaging-finches/hello-k8s/actions/workflows/github-code-scanning/codeql/badge.svg)](https://github.com/engaging-finches/hello-k8s/actions/workflows/github-code-scanning/codeql)

[![Release](https://github.com/engaging-finches/hello-k8s/actions/workflows/release.yaml/badge.svg)](https://github.com/engaging-finches/hello-k8s/actions/workflows/release.yaml)


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

