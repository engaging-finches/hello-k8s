[![Build and Test Go Binary](https://github.com/engaging-finches/hello-k8s/actions/workflows/build-test.yaml/badge.svg)](https://github.com/engaging-finches/hello-k8s/actions/workflows/build-test.yaml)

[![Verify Lint and Format](https://github.com/engaging-finches/hello-k8s/actions/workflows/verify-lint-fmt.yaml/badge.svg)](https://github.com/engaging-finches/hello-k8s/actions/workflows/verify-lint-fmt.yaml)

[![Verify Pre Commit Standards](https://github.com/engaging-finches/hello-k8s/actions/workflows/pre-commit-ci.yaml/badge.svg)](https://github.com/engaging-finches/hello-k8s/actions/workflows/pre-commit-ci.yaml)

[![CodeQL](https://github.com/engaging-finches/hello-k8s/actions/workflows/github-code-scanning/codeql/badge.svg)](https://github.com/engaging-finches/hello-k8s/actions/workflows/github-code-scanning/codeql)

[![Release](https://github.com/engaging-finches/hello-k8s/actions/workflows/release.yaml/badge.svg)](https://github.com/engaging-finches/hello-k8s/actions/workflows/release.yaml)


# hello-k8s

## Description
**MVP: CustomResourceDefinition (CRD), WebHooks, and associated controller to create a Deployment for Pods to run GitHub Actions self-hosted runners.**

The goal of this project is to allow users to easily manage self-hosted GitHub Actions runners in Kubernetes.
<br></br>
## Running Locally


(Clone repo)

#### Create a cluster
`kind create cluster -n ghrunner
`

#### Install cert-manager
`kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.15.0/cert-manager.yaml`

#### Install CRD and controller in cluster
`kubectl apply -f dist/install.yaml`

#### Verify controller is running
`k get pods -n ghrunner-system`

#### Apply sample CR manifest
`k apply -f config/samples/ghrunner_v1_ghrunner.yaml `

#### You should now be able to apply manifests for GhRunner resources.
<br> </br>
## Developer Guide
(Clone repo)

#### Create a cluster for managing your API and resources
`kind create cluster -n ghrunner
`
#### Install cert-manager in cluster
`kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.15.0/cert-manager.yaml`

#### Generate manifests for CRD and sample resource
`make manifests`

#### Install CRD in cluster
`make install` (verify that crd has been installed with `k get crds` -- should see ghrunners... in output)

#### Create a Docker image for controller and use this image to install the controller in your cluster
`make docker-build docker-push IMG=imagename
`
#### Activate controller in cluster
`make deploy IMG=imagename`

#### Create a resource from CRD
`kubectl apply -f config/samples/ghrunner_v1_ghrunner.yaml`


### How to apply changes to the controller
- Make any desired changes to ghrunner_controller.go.
- `make`
- `kubectl delete deployment ghrunner-controller-manager -n ghrunner-system`
- `make docker-build docker-push IMG=imagename`
- `make deploy IMG=imagename`
- `make install`

<br></br>
## Configure commit hooks

**Generate PAT for Husky commit-msg**
PAT requires read access to issues in the repo that you're using.  Save PAT to a variable called `GITHUB_TOKEN` and store this in a `.env` file in the root directory of your repo.

**Install packages from package.json**
`npm install`

