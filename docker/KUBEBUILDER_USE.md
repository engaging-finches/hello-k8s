<!-- create a cluster for managing your API and resources -->
kind create cluster -n ghrunner-test 

<!-- Install cert-manager in cluster -->
kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.15.0/cert-manager.yaml  

<!-- create a Docker image for controller and use this image to install the controller in your cluster -->
make docker-build docker-push IMG=repo/controllername:tag

<!-- generate manifests for CRD and sample resource -->
make manifests

<!-- install CRD in cluster -->
make install

<!-- activate controller in cluster -->
make deploy IMG=repo/controllername:tag

<!-- create a resource from CRD -->
kubectl apply -f config/samples/somecrmanifestname.yaml

