FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum* ./
RUN go mod download && go mod verify

COPY *.go ./
RUN go build -ldflags="-s -w" -o todos-backend .

FROM scratch

WORKDIR /app

COPY --from=builder /app/todos-backend .

USER 65534

ENTRYPOINT ["./todos-backend"]
