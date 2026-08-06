# the-project/todos

## Prerequisites: create a k3d cluster with port mapping

Create a new k3d cluster with an agent node that maps host port `8080` to the
NodePort defined in `manifests/service.yml` (port `30080`):

```bash
k3d cluster create --agents 1 -p "8080:30080@agent:0"
```

This maps host port `8080` to agent node port `30080`, so once the service is
deployed the app will be reachable at `http://localhost:8080`.

## Build, import, and deploy

Build docker image with `docker build -t daduam/dwk-the-project-todos .`

Upload to k3d cluster with `k3d image import daduam/dwk-the-project-todos`

Deploy with `kubectl apply -f manifests`

## Access the app

Once deployed, open `http://localhost:8080` in your browser.

Alternatively, run `kubectl port-forward` to forward port for local testing. 

