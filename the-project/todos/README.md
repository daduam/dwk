# the-project/todos

## Create the k3d cluster

```bash
k3d cluster create dwk-cluster \
  --agents 2 \
  --port 8082:30080@agent:0 \
  --port 8081:80@loadbalancer
```

This creates a cluster with two agent nodes:
- **agent-0**: NodePort 30080 mapped to host port **8082**
- **loadbalancer**: Port 80 mapped to host port **8081**

## Build, import, and deploy

Build docker image with `docker build -t daduam/dwk-the-project-todos .`

Upload to k3d cluster with `k3d image import daduam/dwk-the-project-todos`

Deploy with `kubectl apply -f manifests`

## Access the app

Once deployed with the ingress in place, open `http://localhost:8081` in your browser or run:

```bash
curl http://localhost:8081
```
