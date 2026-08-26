# log-output

Two-app sidecar setup:

- **writer**: generates a random UUID on startup and appends `timestamp: uuid` lines to a file every 5 seconds.
- **reader**: an HTTP server that on `GET /` serves the contents of the information file and the `MESSAGE` env variable, followed by the writer's log file with `Ping / Pongs: ${count}` appended, where the count is fetched over HTTP from the ping-pong app's `GET /pings` endpoint.

Both containers share an `emptyDir` volume mounted at `/app/logs`, so the reader always sees the latest log written by the writer. The count comes straight from the ping-pong service; the reader no longer mounts any shared volume. The information file is provided to the reader through a ConfigMap mounted read-only at `/config`.

## Configuration

### writer

| Env var       | Default               | Description               |
| ------------- | --------------------- | ------------------------- |
| `OUTPUT_FILE` | `/app/logs/output.log` | File the lines are appended to |

### reader

| Env var            | Default                     | Description                                     |
| ------------------ | --------------------------- | ----------------------------------------------- |
| `PORT`             | `8080`                      | Port the server listens on                      |
| `OUTPUT_FILE`      | `/app/logs/output.log`      | Log file to serve                               |
| `PING_PONG_URL`    | `http://ping-pong-svc:3456` | Base URL of the ping-pong app (fetches `/pings`) |
| `INFORMATION_FILE` | `/config/information.txt`   | File whose content is printed on `GET /`        |
| `MESSAGE`          | _(empty)_                   | Value printed as the `MESSAGE` env variable     |

If the ping-pong app is unreachable or returns an error, the reader logs a warning and prints `Ping / Pongs: 0`.

## Run locally

```bash
cd writer && OUTPUT_FILE=./output.log go run .   # appends to ./output.log
echo "this text is from file" > information.txt
cd reader && OUTPUT_FILE=../writer/output.log PING_PONG_URL=http://localhost:8081 INFORMATION_FILE=../information.txt MESSAGE=hello go run .   # serves the writer's log on localhost:8080
curl localhost:8080                # info file content + "env variable: MESSAGE=hello" + log contents + "Ping / Pongs: 0"
```

Point the reader at the writer's file with `OUTPUT_FILE` and at the ping-pong app with `PING_PONG_URL` (the reader requests `${PING_PONG_URL}/pings`). The information file must exist (its default `/config/information.txt` is provided by the ConfigMap in the cluster, not on the local machine); `MESSAGE` is optional and prints empty when unset.

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

The deployment shares an `emptyDir` volume between the writer and the reader at `/app/logs`, sets `PING_PONG_URL` so the reader fetches the count from the ping-pong service, and uses the `log-output-cm` ConfigMap to provide the reader with `MESSAGE` and the information file mounted at `/config`:

```bash
kubectl apply -f manifests
```

## Access

Once the ingress is in place, open `http://localhost:8081` or run:

```bash
curl http://localhost:8081   # info file content, MESSAGE, log lines, then "Ping / Pongs: ${count}"
```
