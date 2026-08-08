# ping-pong

A simple HTTP server that returns `pong ${count}` on each request to `/pingpong`, where the count starts at 0 and increments by 1 after every response.

## Build, import, and deploy

Build docker image with `docker build -t daduam/dwk-ping-pong .`

Upload to k3d cluster with `k3d image import daduam/dwk-ping-pong`

Deploy with `kubectl apply -f manifests`

## Access the app

Once deployed with the ingress (this application shares the [log-output-ingress](../log-output/manifests/ingress.yaml)) in place, open `http://localhost:8081/pingpong` in your browser or run:

```bash
curl http://localhost:8081/pingpong
```
