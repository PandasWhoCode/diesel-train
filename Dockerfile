# syntax=docker/dockerfile:1

# ── Build stage ───────────────────────────────────────────────────────────────
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Copy module files first for layer-cached dependency download.
COPY go.mod go.sum ./
RUN go mod download

COPY main.go .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w" -trimpath -o diesel .

# ── Final stage — distroless (no shell, no package manager, minimal surface) ──
FROM gcr.io/distroless/static-debian12:nonroot

COPY --from=builder /app/diesel /diesel

ENTRYPOINT ["/diesel"]
