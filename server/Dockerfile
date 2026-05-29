FROM node:20-alpine AS web-builder
WORKDIR /build
COPY web/package.json web/vite.config.js web/index.html ./
COPY web/src ./src
COPY web/public ./public
RUN npm install && npm run build

FROM golang:1.23-alpine AS server-builder
WORKDIR /build
COPY server/go.mod server/go.sum ./
RUN go mod download
COPY server/ .
RUN CGO_ENABLED=0 GOOS=linux go build -o /consumable-tracker .

FROM postgres:16-alpine
RUN apk add --no-cache ca-certificates tzdata openssl

ENV POSTGRES_DB=consumable_db
ENV POSTGRES_USER=admin

COPY --from=server-builder /consumable-tracker /usr/local/bin/
COPY --from=web-builder /build/dist /app/web-dist
COPY database/init.sql /docker-entrypoint-initdb.d/
COPY scripts/docker-cmd.sh /usr/local/bin/
RUN chmod +x /usr/local/bin/docker-cmd.sh

RUN mkdir -p /backups /certs && chown -R postgres:postgres /backups /certs

EXPOSE 8443

HEALTHCHECK --interval=30s --timeout=5s --start-period=15s --retries=3 \
    CMD pg_isready -U admin -d consumable_db -q && wget -q --spider http://localhost:8443/ || exit 1

CMD ["docker-cmd.sh"]
