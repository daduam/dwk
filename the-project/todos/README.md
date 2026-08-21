# the-project/todos

The todos app frontend: serves server-rendered HTML and the daily image.
It depends on [the-project/todos-backend](../todos-backend) for the todo list.

- The page fetches the todo list from the backend pod-to-pod via `TODO_API_URL` and renders it server-side.
- The form submits to `/api/todos`, which the ingress routes to the backend's `POST /todos` (the `strip-api` middleware removes the `/api` prefix).

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

Build both images:

```bash
docker build -t daduam/dwk-the-project-todos .
cd ../todos-backend && docker build -t daduam/dwk-the-project-todos-backend .
```

Upload them to the k3d cluster:

```bash
k3d image import daduam/dwk-the-project-todos
k3d image import daduam/dwk-the-project-todos-backend
```

Deploy the shared volume, the backend, and the frontend:

```bash
kubectl apply -f ../../manifests
kubectl apply -f ../todos-backend/manifests
kubectl apply -f manifests
```

The frontend requires `TODO_API_URL` to be set; the deployment manifest points it at `http://the-project-todos-backend-svc:1234`.

## Access the app

Once deployed with the ingress in place, open `http://localhost:8081` in your browser or run:

```bash
curl http://localhost:8081
```

The todo list is rendered server-side from the backend API. Adding a todo posts to `/api/todos`, which the ingress routes to the backend.

## Access the API directly

The backend API is exposed under `/api` through the same ingress:

```bash
curl http://localhost:8081/api/todos
```
