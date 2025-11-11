ARG GO_VERSION=1.24.0
FROM golang:${GO_VERSION}-bookworm as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

RUN go build -v -o scrapper ./cmd/scrapper/main.go

# ----------------------------
# Scrape Service
FROM debian:bookworm-slim AS scrape-cs-items
RUN apt-get update && apt-get install -y ca-certificates && \
    rm -rf /var/lib/apt/lists/* 
COPY --from=builder /app/scrapper /usr/local/bin/scrapper
CMD ["scrapper"]