# log-output

A simple HTTP server that returns a timestamp and a random UUID (generated on startup) on each request.

## Build, import, and deploy

Build docker image with `docker build -t daduam/dwk-log-output .`

Upload to k3d cluster with `k3d image import daduam/dwk-log-output`

Deploy with `kubectl apply -f manifests`

## Access the app

Once deployed with the ingress in place, open `http://localhost:8081` in your browser or run:

```bash
curl http://localhost:8081
```
