FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum* ./
RUN go mod download && go mod verify

COPY main.go ./
RUN go build -ldflags="-s -w" -o todos main.go

FROM scratch

WORKDIR /app

COPY --from=builder /app/todos .
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

USER 65534

ENTRYPOINT ["./todos"]
