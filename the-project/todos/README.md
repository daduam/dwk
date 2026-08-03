# the-project/todos

Build docker image with `docker build -t daduam/dwk-the-project-todos .`

Upload to k3d cluster with `k3d image import daduam/dwk-the-project-todos`

Deploy with `kubectl apply -f manifests`

Run `kubectl port-forward` to foward port for local testing. 

