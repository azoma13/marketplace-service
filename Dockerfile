FROM golang:alpine AS modules
COPY go.mod go.sum /modules/
WORKDIR /modules
RUN go mod download

FROM golang:alpine AS builder
COPY --from=modules /go/pkg /go/pkg
COPY . /app
WORKDIR /app
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -tags migrate -o /bin/app ./cmd/app

FROM scratch
COPY --from=builder /app/config /config
COPY --from=builder /app/migrations /migrations
COPY --from=builder /app/.env /.env
COPY --from=builder /bin/app /app
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
CMD ["/app"]