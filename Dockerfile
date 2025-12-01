ARG GO_VERSION=1.24.0
FROM golang:${GO_VERSION}-bookworm as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

RUN go build -v -o scrapper ./cmd/scrapper/main.go
RUN go build -v -o alert ./cmd/alert/main.go
RUN go build -v -o api ./cmd/api/main.go

# ----------------------------
# Scrape Service
FROM debian:bookworm-slim AS scrapper
RUN apt-get update && apt-get install -y ca-certificates && \
    rm -rf /var/lib/apt/lists/* 
COPY --from=builder /app/scrapper /usr/local/bin/scrapper
CMD ["scrapper"]

# ----------------------------
# Alert Service
FROM debian:bookworm-slim AS alert
RUN apt-get update && apt-get install -y ca-certificates && \
    rm -rf /var/lib/apt/lists/* 
COPY --from=builder /app/alert /usr/local/bin/alert
CMD ["alert"]

# ----------------------------
# API Service
FROM debian:bookworm-slim AS api
RUN apt-get update && apt-get install -y ca-certificates && \
    rm -rf /var/lib/apt/lists/* 
COPY --from=builder /app/api /usr/local/bin/api
CMD ["api"]