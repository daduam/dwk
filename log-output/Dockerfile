FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY main.go ./
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o log-output main.go

FROM scratch

COPY --from=builder /app/log-output /log-output

USER 65534

ENTRYPOINT ["/log-output"]

