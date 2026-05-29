FROM node:20-alpine AS web-builder
WORKDIR /build
COPY web/package.json web/vite.config.js web/index.html ./
COPY web/src ./src
RUN npm install && npm run build

FROM golang:1.23-alpine AS server-builder
RUN apk add --no-cache gcc musl-dev
WORKDIR /build
COPY server/go.mod server/go.sum ./
RUN go mod download
COPY server/ .
RUN CGO_ENABLED=0 GOOS=linux go build -o /consumable-tracker .

FROM alpine:3.20
RUN apk add --no-cache ca-certificates tzdata openssl postgresql-client

RUN adduser -D -u 1000 appuser

COPY --from=server-builder /consumable-tracker /usr/local/bin/
COPY --from=web-builder /build/dist /app/web-dist

RUN mkdir -p /backups /certs && chown -R appuser:appuser /backups /certs /app

USER appuser
WORKDIR /app

EXPOSE 8443

HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8443/api/health || exit 1

ENTRYPOINT ["/usr/local/bin/consumable-tracker"]
