FROM golang:1.23.2 AS builder

WORKDIR /app

COPY go.mod go.sum .
RUN go mod download

COPY . .

RUN go build -o /bin/http-server ./cmd/http-server \
    && go build -o /bin/migrator ./cmd/migrator


FROM debian:bookworm-slim

WORKDIR /app

COPY --from=builder /bin/http-server /bin/migrator /usr/local/bin/
COPY --from=builder /app/configs configs/
COPY --from=builder /app/migrations migrations/

COPY deploy/entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]