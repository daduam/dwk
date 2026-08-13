# ping-pong

A simple HTTP server that returns `pong ${count}` on each request, where `count` is the total number of requests the app has served.

The count is persisted to a file on every request and restored from it at startup, so it survives pod restarts. In the cluster the file lives on a PersistentVolume so the count is kept even if the pod is rescheduled.

## Configuration

| Env var               | Default                 | Description                                 |
| --------------------- | ----------------------- | ------------------------------------------- |
| `PORT`                | `8080`                  | Port the server listens on                  |
| `REQUESTS_COUNT_FILE` | `requests_count.txt`    | File the request count is stored in         |

## Run locally

```bash
go run .
curl localhost:8081   # pong 1
curl localhost:8081   # pong 2
```

## Build and import into k3d

```bash
docker build -t daduam/dwk-ping-pong .
k3d image import daduam/dwk-ping-pong
```

## Deploy

The deployment mounts the PVC `shared-data-pvc` at `/app/shared` and sets
`REQUESTS_COUNT_FILE=/app/shared/requests_count.txt`. Create the PV and PVC first:

```bash
kubectl apply -f ../manifests/persistentvolume.yaml
kubectl apply -f ../manifests/persistentvolumeclaim.yaml
kubectl apply -f manifests
```

## Access the app

The deployment shares the [log-output-ingress](../log-output/manifests/ingress.yaml). Once it's in place, open `http://localhost:8081/pingpong` or run:

```bash
curl http://localhost:8081/pingpong   # pong ${count}
```
