# the-project/todos-backend

A simple Go HTTP backend for managing todos. Todos are stored in memory.
The server starts with three example todos.

## Run locally

```bash
go build -o todos-backend .
./todos-backend
```

The server listens on port **8080** by default. Set the `PORT` environment variable to change it:

```bash
PORT=3000 ./todos-backend
```

## Run with Docker

```bash
docker build -t daduam/dwk-the-project-todos-backend .
docker run -p 8080:8080 daduam/dwk-the-project-todos-backend
```

## Deploy to Kubernetes

Build the docker image with `docker build -t daduam/dwk-the-project-todos-backend .`

Upload it to the k3d cluster with `k3d image import daduam/dwk-the-project-todos-backend`

Deploy with `kubectl apply -f manifests`

Once the the-project ingress is deployed, it routes `/api` to this backend (stripping the `/api` prefix). From the host:

```bash
curl http://localhost:8081/api/todos
```

## API

The backend serves the API at its root. Through the cluster ingress, the same endpoints are available under `/api`.

### List todos

```bash
curl http://localhost:8080/todos
```

Returns a JSON array of todos:

```json
[
  {
    "id": 1,
    "content": "Write a README",
    "done": false,
    "createdAt": "2025-08-21T09:55:00Z"
  }
]
```

### Create a todo

```bash
curl -X POST http://localhost:8080/todos \
  -H "Content-Type: application/json" \
  -d '{"content": "Write a README"}'
```

Returns the created todo with status `201`:

```json
{
  "id": 1,
  "content": "Write a README",
  "done": false,
  "createdAt": "2025-08-21T09:55:00Z"
}
```

- `content` is required; requests without it get `400 Bad Request`.
- Malformed request bodies get `400 Bad Request`.

## Data model

| Field       | Type     | Description                            |
| ----------- | -------- | -------------------------------------- |
| `id`        | `int`    | Auto-incremented todo ID               |
| `content`   | `string` | Todo text                              |
| `done`      | `bool`   | Completion status, defaults to `false` |
| `createdAt` | `time`   | Creation timestamp                     |

Todos live in an in-memory store, so data is lost when the server restarts.
