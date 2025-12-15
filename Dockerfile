# Build stage
FROM golang:1.25.5-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o url-shortner ./cmd/...

# Final stage
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/url-shortner .
COPY --from=builder /app/deployments ./deployments

ENV PORT=80

EXPOSE 80

ENTRYPOINT ["./url-shortner"]