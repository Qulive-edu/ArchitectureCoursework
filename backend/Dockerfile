FROM golang:1.24-alpine AS builder
RUN apk add --no-cache git
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /app/main cmd/app/main.go

FROM alpine
RUN apk add --no-cache ca-certificates
COPY --from=builder /app/main /app/main
COPY config/config.yml /app/config/config.yml
WORKDIR /app
CMD ["./main"]
