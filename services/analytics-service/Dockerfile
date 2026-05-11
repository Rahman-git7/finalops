FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /analytics-service ./cmd/server

FROM scratch
COPY --from=builder /analytics-service /analytics-service
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
EXPOSE 8004
ENTRYPOINT ["/analytics-service"]
