FROM golang:1.22-alpine AS builder

WORKDIR /app

RUN apk add --no-cache ca-certificates

COPY go.mod ./
RUN go mod download

COPY cmd ./cmd
COPY internal ./internal

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/trainer ./cmd/trainer

FROM alpine:3.22

WORKDIR /app

RUN apk add --no-cache ca-certificates && update-ca-certificates

COPY --from=builder /out/trainer /usr/local/bin/trainer

VOLUME ["/data", "/output"]

ENTRYPOINT ["trainer"]
CMD ["-data", "/output/clean_records_sample.json", "-out", "/output/train_results.json"]