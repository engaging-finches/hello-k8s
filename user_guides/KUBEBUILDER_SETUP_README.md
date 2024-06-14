## This is a guide for starting a fresh project.  To use a clone of the hello-k8s repo please refer to DEVELOPER_GUIDE.md

### create a cluster for managing your API and resources
kind create cluster -n ghrunner

### Install cert-manager in cluster
kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.15.0/cert-manager.yaml  

### generate scaffolding for CRD and controller
mkdir ghrunner && cd ghrunner
kubebuilder init --domain ghrunner  --repo ghrunner/selfhosted 


### plugins are used to generate deployment config for container -- IMPORTANT
kubebuilder create api \
  --group ghrunner \
  --version v1 \
  --kind GhRunner \
  --image=meherliatrio/selfhosted:latest \
  --plugins="deploy-image/v1-alpha" \
  --run-as-user="1000" 
<!-- NOT CERTAIN HOW USEFUL THESE ARE
  --image-container-command="memcached,-m=64,modern,-v" \
  --image-container-port="11211" \ -->

### modify api/v1/ghrunner_types.go, internal/controller/ghrunner_controller.go -- replicate file contents from engaging-finches/hello-k8s main branch

### apply changes to ghrunner_types.go and/or controller
make

### create a Docker image for controller and use this image to install the controller in your cluster
make docker-build docker-push IMG=repo/controllername:tag

### generate manifests for CRD and sample resource
make manifests

### install CRD in cluster
make install

### activate controller in cluster
make deploy IMG=repo/controllername:tag

### create a resource from CRD
kubectl apply -f config/samples/somecrmanifestname.yaml
