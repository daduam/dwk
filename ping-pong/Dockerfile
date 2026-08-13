FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.* ./
RUN go mod download && go mod verify

COPY main.go ./
RUN go build -ldflags="-s -w" -o ping-pong main.go

FROM scratch

WORKDIR /app

COPY --from=builder /app/ping-pong .

USER 65534

ENTRYPOINT ["./ping-pong"]
