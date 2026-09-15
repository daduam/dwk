# ping-pong

A simple HTTP server that returns `pong ${count}` on each request, where `count` is the total number of requests the app has served.

The count lives in a PostgreSQL table and is updated with a single atomic upsert, so it survives pod restarts and stays correct under concurrent requests without any application-side locking.

## Endpoints

| Route        | Response        |
| ------------ | --------------- |
| `GET /`      | `pong ${count}` |
| `GET /pings` | `${count}`      |

## Configuration

| Env var        | Default | Description                |
| -------------- | ------- | -------------------------- |
| `PORT`         | `8080`  | Port the server listens on |
| `DATABASE_URL` | —       | PostgreSQL DSN, required   |

The app exits non-zero if `DATABASE_URL` is missing or the database is unreachable, so Kubernetes restarts it until the database is up.

## Schema

The app creates its own table on startup, idempotently, so an existing count is never reset: a `counters` table with a database-assigned `id`, a unique `title` slug, and a `value`. The single counter is addressed by the slug `requests-count`, and the unique index on `title` is what makes the `ON CONFLICT` upsert work.

Incrementing and reading happen in one atomic statement, which is what removes the need for an application mutex. The SQL lives in [`store.go`](store.go).

## Run locally

```bash
docker run -d --name ping-pong-postgres \
  -e POSTGRES_PASSWORD=postgres \
  -p 5433:5432 postgres:17.6-alpine

export DATABASE_URL='postgres://postgres:postgres@localhost:5433/postgres'
export PORT=8090
go run .

curl localhost:8090         # pong 1
curl localhost:8090         # pong 2
curl localhost:8090/pings   # 2
```

## Build and import into k3d

```bash
docker build -t daduam/dwk-ping-pong .
k3d image import daduam/dwk-ping-pong -c k3s-default
```

## Deploy

Needs the `exercises` namespace and the PostgreSQL StatefulSet in `manifests/postgres/`:

```bash
kubectl apply -f ../manifests/namespaces.yaml

kubectl apply -f manifests/postgres/
kubectl -n exercises rollout status statefulset/postgres-stset --timeout=180s

kubectl apply -f manifests/deployment.yaml
kubectl -n exercises rollout status deploy/ping-pong --timeout=120s
```

`kubectl apply -f manifests` on its own skips `postgres/`, because `apply` is not recursive.

`DATABASE_URL` addresses a single pod:

```
postgres://postgres:postgres@postgres-stset-0.postgres-svc.exercises:5432/postgres
```

A StatefulSet names its pods `<statefulset>-<ordinal>.<service>.<namespace>`, so renaming the StatefulSet changes this host.

## Access the app

The route lives in [log-output's ingress](../log-output/manifests/ingress.yaml), which serves `/` from `log-output-svc` and `/pingpong` from `ping-pong-svc`.

```bash
curl http://localhost:8081/pingpong   # pong ${count}
```

`GET /pings` is not reachable through the ingress. Use `kubectl -n exercises port-forward deploy/ping-pong 9090:8081` and `curl localhost:9090/pings`.
