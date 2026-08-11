FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum* ./
RUN go mod download && go mod verify

COPY main.go ./
RUN go build -ldflags="-s -w" -o log-output-writer main.go

FROM scratch

WORKDIR /app

COPY --from=builder /app/log-output-writer .

USER 65534

ENTRYPOINT ["./log-output-writer"]
