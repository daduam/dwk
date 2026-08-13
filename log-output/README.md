# log-output

Two-app sidecar setup:

- **writer**: generates a random UUID on startup and appends `timestamp: uuid` lines to a file every 5 seconds.
- **reader**: an HTTP server that serves the contents of that file on `GET /`, with `Ping / Pongs: ${count}` appended, where the count is read from the ping-pong app's count file.

Both containers share an `emptyDir` volume mounted at `/app/logs`, so the reader always sees the latest log written by the writer. The count file lives on the PVC `shared-data-pvc` (shared with the [ping-pong](../ping-pong) app) at `/app/shared/requests_count.txt`.

## Configuration

### writer

| Env var       | Default               | Description               |
| ------------- | --------------------- | ------------------------- |
| `OUTPUT_FILE` | `/app/logs/output.log` | File the lines are appended to |

### reader

| Env var               | Default                              | Description                                    |
| --------------------- | ------------------------------------ | ---------------------------------------------- |
| `PORT`                | `8080`                               | Port the server listens on                     |
| `OUTPUT_FILE`         | `/app/logs/output.log`               | Log file to serve                              |
| `REQUESTS_COUNT_FILE` | `requests_count.txt`                 | File the ping-pong request count is stored in  |

If the count file is missing or unreadable, the reader logs a warning and prints `Ping / Pongs: 0`.

## Run locally

```bash
cd writer && OUTPUT_FILE=./output.log go run .   # appends to ./output.log
cd reader && OUTPUT_FILE=../writer/output.log go run .   # serves the writer's log on localhost:8080
curl localhost:8080                # log contents + "Ping / Pongs: 0"
```

Point the reader at the writer's file and the ping-pong count file with `OUTPUT_FILE` and `REQUESTS_COUNT_FILE` if they aren't in the default locations.

## Build and import into k3d

```bash
# writer
cd writer
docker build -t daduam/dwk-log-output-writer .

# reader
cd reader
docker build -t daduam/dwk-log-output-reader .

k3d image import daduam/dwk-log-output-writer
k3d image import daduam/dwk-log-output-reader
```

## Deploy

The deployment mounts the PVC `shared-data-pvc` into the reader at `/app/shared` (where it reads the ping-pong count) and shares an `emptyDir` with the writer at `/app/logs`. Create the PV and PVC first:

```bash
kubectl apply -f ../manifests/persistentvolume.yaml
kubectl apply -f ../manifests/persistentvolumeclaim.yaml
kubectl apply -f manifests
```

## Access

Once the ingress is in place, open `http://localhost:8081` or run:

```bash
curl http://localhost:8081   # log lines followed by "Ping / Pongs: ${count}"
```
