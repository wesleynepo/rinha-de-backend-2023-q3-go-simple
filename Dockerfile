FROM golang:1.20-buster as builder

WORKDIR /app

COPY go.* ./
RUN go mod download

COPY . ./

RUN go build -ldflags='-s -X main.buildTime=${current_time}' -o=./bin/api ./cmd/api 

FROM debian:buster-slim
RUN set -x && apt-get update && DEBIAN_FRONTEND=noninteractive apt-get install -y \
    ca-certificates && \
    rm -rf /var/lib/apt/lists/*

COPY --from=builder /app/bin/api /app/api

CMD ["/app/api"]
