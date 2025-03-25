ARG GO_VERSION=1.24.0
FROM golang:${GO_VERSION}-bookworm as builder

WORKDIR /usr/src/app
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .
RUN go build -v -o /bin/scrape-cs-items ./cmd/scrape-cs-items
RUN go build -o /bin/alerts ./cmd/alerts

# ----------------------------
# Scrape Service
FROM debian:bookworm-slim AS scrape-cs-items
RUN apt-get update && apt-get install -y ca-certificates && \
    rm -rf /var/lib/apt/lists/*
COPY --from=builder /bin/scrape-cs-items /usr/local/bin
CMD ["scrape-cs-items"]

# ----------------------------
# Alerts Service
FROM debian:bookworm-slim AS alerts
RUN apt-get update && apt-get install -y ca-certificates && \
    rm -rf /var/lib/apt/lists/*
COPY --from=builder /bin/alerts /usr/local/bin
CMD ["alerts"]