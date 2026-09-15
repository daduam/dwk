# the-project/todos-backend

A Go HTTP backend for managing todos, backed by PostgreSQL, so the todos survive
pod restarts. The app creates its table on startup and seeds three example todos
the first time it finds the table empty.

## Configuration

| Env var        | Default | Description                |
| -------------- | ------- | -------------------------- |
| `PORT`         | `8080`  | Port the server listens on |
| `DATABASE_URL` | —       | PostgreSQL DSN, required   |

The app exits non-zero if `DATABASE_URL` is missing or the database is
unreachable, so Kubernetes restarts it until the database is up.

## Run locally

```bash
docker run -d --name todos-postgres \
  -e POSTGRES_PASSWORD=postgres -e POSTGRES_DB=todos \
  -p 5434:5432 postgres:17.6-alpine

export DATABASE_URL='postgres://postgres:postgres@localhost:5434/todos'
export PORT=8090
go run .
```

## API

The backend serves the API at its root. Through the cluster ingress the same
endpoints are available under `/api`.

### List todos

```bash
curl localhost:8090/todos
```

Returns a JSON array of todos, newest first:

```json
[{ "id": 1, "content": "Write a README", "done": false, "createdAt": "2025-08-21T09:55:00Z" }]
```

### Create a todo

```bash
curl -X POST localhost:8090/todos \
  -H "Content-Type: application/json" \
  -d '{"content": "Write a README"}'
```

Returns the created todo with status `201`. `content` is required, and malformed
request bodies get `400 Bad Request`.

## Schema

`todos` has a database-assigned `bigint id`, the `text content`, a `boolean done`
that defaults to `false`, and a `timestamptz created_at` that defaults to
`now()`. The app creates the table on startup, idempotently, so existing data is
never reset, and seeds only when the table is empty, so restarts neither
duplicate the examples nor wipe real rows. The SQL lives in [`store.go`](store.go).

## Deploy to Kubernetes

The backend needs the PostgreSQL StatefulSet and the `postgres-credentials`
Secret from [../manifests/postgres](../manifests/postgres), which is where the
deployment reads `DATABASE_URL` from. See [../README.md](../README.md) for the
full sequence; from this directory:

```bash
docker build -t daduam/dwk-the-project-todos-backend .
k3d image import daduam/dwk-the-project-todos-backend -c k3s-default
kubectl apply -f ../manifests/todos-backend/
kubectl -n project rollout restart deployment/the-project-todos-backend
```

`DATABASE_URL` addresses a single pod rather than the service, and its database
name must match `POSTGRES_DB`, which is `postgres` by default because the
StatefulSet sets only `POSTGRES_PASSWORD`:
`postgres://postgres:<password>@postgres-stset-0.postgres-svc.project:5432/postgres`.
A StatefulSet names its pods `<statefulset>-<ordinal>.<service>.<namespace>`, so
renaming the StatefulSet or its headless service changes this host.

Through the ingress, from the host:

```bash
curl http://localhost:8081/api/todos
```
