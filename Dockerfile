# ── Build stage ──────────────────────────────────────────────────────────────
FROM golang:1.22-alpine AS builder

WORKDIR /app
COPY go.mod ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o crimes-loader ./cmd/loader

# ── Runtime stage ─────────────────────────────────────────────────────────────
FROM alpine:3.19

WORKDIR /app
COPY --from=builder /app/crimes-loader .

# Directorio para montar el dataset desde fuera del contenedor
RUN mkdir -p /app/data /app/output

# Variables de entorno con defaults configurables
ENV WORKERS=8
ENV CHUNK_SIZE=5000
ENV VERBOSE=true
ENV DATA_FILE=/app/data/crimes.csv
ENV OUTPUT_DIR=/app/output

ENTRYPOINT ["sh", "-c", \
  "./crimes-loader -file=$DATA_FILE -workers=$WORKERS -chunk=$CHUNK_SIZE -output=$OUTPUT_DIR"]
