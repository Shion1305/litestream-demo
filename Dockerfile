FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY server/ ./server/
WORKDIR /app/server
COPY go.mod go.sum ./
RUN go mod download
RUN apk add --no-cache gcc musl-dev
RUN CGO_ENABLED=1 go build -o /app/main

# 実行ステージ
FROM alpine:3.18
RUN apk add --no-cache ca-certificates curl sqlite

# Litestream バイナリを取得（アーキテクチャ別）
RUN if [ "$(uname -m)" = "x86_64" ]; then \
      curl -fsSL https://github.com/benbjohnson/litestream/releases/download/v0.3.13/litestream-v0.3.13-linux-amd64.tar.gz | tar xz -C /usr/local/bin; \
    elif [ "$(uname -m)" = "aarch64" ]; then \
      curl -fsSL https://github.com/benbjohnson/litestream/releases/download/v0.3.13/litestream-v0.3.13-linux-arm64.tar.gz | tar xz -C /usr/local/bin; \
    else \
        echo "Unsupported architecture: $(uname -m)"; exit 1; \
    fi

WORKDIR /app
COPY --from=builder /app/main ./main
COPY ./static ./static
COPY litestream.yml .
COPY docker-entrypoint.sh .
RUN chmod +x docker-entrypoint.sh

EXPOSE 8080
ENTRYPOINT ["/app/docker-entrypoint.sh"]