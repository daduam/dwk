# log-output

Two-app sidecar setup:

- **writer**: generates a random UUID on startup and appends `timestamp: uuid` lines to a file every 5 seconds.
- **reader**: an HTTP server that serves the contents of that file on `GET /`.

Both containers share an `emptyDir` volume mounted at `/app/logs`, so the reader always sees the latest log written by the writer.

## Build

```bash
# writer
cd writer
docker build -t daduam/dwk-log-output-writer .

# reader
cd reader
docker build -t daduam/dwk-log-output-reader .
```

## Import into k3d

```bash
k3d image import daduam/dwk-log-output-writer
k3d image import daduam/dwk-log-output-reader
```

## Deploy

```bash
kubectl apply -f manifests/
```

## Access

Once the ingress is in place, open `http://localhost:8081` or run:

```bash
curl http://localhost:8081
```
