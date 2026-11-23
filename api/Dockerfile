# Build stage
FROM golang:1.27-alpine AS builder
WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/server

# Run stage
FROM alpine:3.24
RUN apk --no-cache add ca-certificates \
    && adduser -D -H -u 10001 app
WORKDIR /app
COPY --from=builder /app/server .
USER app
EXPOSE 8080
CMD ["./server"]
