ARG GO_VERSION=1.24.0
FROM golang:${GO_VERSION}-bookworm as builder

WORKDIR /usr/src/app
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .
RUN go build -v -o /scrape-cs-items 

FROM debian:bookworm

COPY --from=builder /scrape-cs-items /usr/local/bin
CMD ["scrape-cs-items"]