### clone repo

### create a cluster
kind create cluster -n ghrunner

### install CRD and controller in cluster
kubectl apply -f dist/install.yaml    

### you should now be able to apply manifests for GhRunner resources.  