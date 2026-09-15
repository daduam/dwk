# the-project

The project app: a todo frontend, a todo backend, and the PostgreSQL database
they share, running in the `project` namespace of the local k3d cluster. The
ingress serves `/` from the frontend and `/api` from the backend, stripping the
`/api` prefix.

```
the-project/
├── manifests/            namespace, PV/PVC, ingress, strip-api middleware
│   ├── postgres/         secret.enc.yaml, service.yaml, statefulset.yaml
│   ├── todos/            frontend deployment and service
│   └── todos-backend/    backend deployment and service
├── todos/                frontend source and Dockerfile
└── todos-backend/        backend source and Dockerfile
```

- **`todos/`** fetches the todo list from the backend pod-to-pod and renders it
  server-side on the `shared-data` volume's cached image.
- **`todos-backend/`** serves `GET /todos` and `POST /todos`.

## Prerequisites

```bash
k3d cluster create k3s-default --agents 2 --port 8082:30080@agent:0 --port 8081:80@loadbalancer
```

Also needs `docker`, `kubectl`, `sops`, and `age`.

## Secrets

The credentials are committed encrypted; the plaintext and the private key are not.

```bash
age-keygen -o ~/.config/sops/age/keys.txt                        # once per machine
sops manifests/postgres/secret.enc.yaml                          # edit, re-encrypted on save
sops -d manifests/postgres/secret.enc.yaml | kubectl apply -f -   # apply
```

Encryption protects the git copy only; anyone with read access to Secrets in the
cluster can still print the values.

## Deploy

Images are built locally and imported; the manifests use fixed tags with
`imagePullPolicy: IfNotPresent`, so a rebuilt image needs a restart.

```bash
cd the-project

kubectl apply -f manifests/                       # namespace, PV/PVC, ingress, middleware

sops -d manifests/postgres/secret.enc.yaml | kubectl apply -f -
kubectl apply -f manifests/postgres/service.yaml
kubectl apply -f manifests/postgres/statefulset.yaml
kubectl -n project rollout status statefulset/postgres-stset --timeout=180s

docker build -t daduam/dwk-the-project-todos-backend todos-backend
k3d image import daduam/dwk-the-project-todos-backend -c k3s-default
kubectl apply -f manifests/todos-backend/
kubectl -n project rollout restart deployment/the-project-todos-backend
kubectl -n project rollout status deployment/the-project-todos-backend --timeout=120s

docker build -t daduam/dwk-the-project-todos todos
k3d image import daduam/dwk-the-project-todos -c k3s-default
kubectl apply -f manifests/todos/
kubectl -n project rollout status deployment/the-project-todos --timeout=120s
```

The StatefulSet `rollout status` waits for the `pg_isready` readiness probe, so
the backend does not start against a database that is still coming up.

## Access

```bash
open http://localhost:8081
curl http://localhost:8081/api/todos
curl -X POST http://localhost:8081/api/todos \
  -H 'Content-Type: application/json' -d '{"content":"hello"}'
```

## Notes

- `kubectl apply` is not recursive: `kubectl apply -f manifests/` applies the flat
  files only, so the service directories go in individually. `manifests/postgres/`
  cannot be applied as a directory at all, since it holds ciphertext; apply the
  Secret from the decrypted stream and the other two files by name.
- Todos live in the `data-postgres-stset-0` volume and the image cache in
  `shared-data-pvc`, so deleting either pod keeps the data.
- The StatefulSet sets only `POSTGRES_PASSWORD`, so the image defaults to user
  `postgres` and database `postgres`. `DATABASE_URL` must agree, and addresses a
  single pod rather than the service:
  `postgres://postgres:<password>@postgres-stset-0.postgres-svc.project:5432/postgres`.
