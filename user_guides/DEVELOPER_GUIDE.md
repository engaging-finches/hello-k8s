### clone repo

### create a cluster for managing your API and resources
kind create cluster -n ghrunner

### Install cert-manager in cluster
kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.15.0/cert-manager.yaml  

### generate manifests for CRD and sample resource
make manifests

### install CRD in cluster
make install (verify that crd has been installed with k get crds -- should see ghrunners... in output)

### activate controller in cluster
make deploy IMG=imagename

### create a resource from CRD
kubectl apply -f config/samples/somecrmanifestname.yaml

---

## Applying changes to the controller
- Make any desired changes to ghrunner_controller.go. 
- make 
- kubectl delete deployment ghrunner-controller-manager -n ghrunner-system
- make docker-build docker-push IMG=imagename
- make deploy IMG=imagename