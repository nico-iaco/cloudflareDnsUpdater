# file: Dockerfile
# ---------- Builder ----------
FROM golang:1.24-alpine AS builder

WORKDIR /src

# Dipendenze di base
RUN apk add --no-cache git ca-certificates

# Pre-carica i mod
COPY go.mod ./
# Se presente, sblocca anche go.sum
# COPY go.sum ./
RUN go mod download

# Copia il resto del sorgente
COPY . .

# Build statica, piccola e riproducibile
ENV CGO_ENABLED=0
RUN go build -trimpath -ldflags="-s -w" -o /out/app .

# ---------- Runtime ----------
FROM alpine:3.20

# Certificati per HTTPS e utente non root
RUN apk add --no-cache ca-certificates \
 && addgroup -S app && adduser -S app -G app

WORKDIR /app
COPY --from=builder /out/app /app/app

USER app
ENTRYPOINT ["/app/app"]