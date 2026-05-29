# Stage 1: Build Vue.js frontend
FROM node:20-alpine AS vue-builder
WORKDIR /web
COPY web/package.json ./
RUN npm install
COPY web/ .
RUN npm run build

# Stage 2: Build Go server
FROM golang:1.24-alpine AS go-builder
WORKDIR /app
COPY server/go.* ./
RUN go mod download
COPY server/ .
COPY --from=vue-builder /web/dist ./web-dist
    RUN go mod tidy && CGO_ENABLED=0 go build -ldflags="-s -w" -o server .

# Stage 3: Runtime
FROM alpine:3.20
RUN adduser -D appuser && apk add --no-cache ca-certificates tzdata
COPY --from=go-builder /app/server /app/server
RUN mkdir -p /app/data /app/backups /app/certs && chown -R appuser:appuser /app
USER appuser
WORKDIR /app
EXPOSE 8443
VOLUME ["/app/data", "/app/backups", "/app/certs"]
CMD ["/app/server"]
