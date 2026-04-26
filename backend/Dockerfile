FROM golang:latest AS builder
RUN apt-get update && apt-get install -y git && rm -rf /var/lib/apt/lists/*
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o /app/main cmd/app/main.go

FROM alpine
RUN apk add --no-cache ca-certificates
COPY --from=builder /app/main /app/main
COPY config/config.yml /app/config/config.yml
COPY migrations/ /app/migrations/
WORKDIR /app

ARG COMMIT_SHA=unknown
ARG BUILD_DATE=unknown
LABEL org.opencontainers.image.source="https://github.com/Qulive-edu/ArchitectureCoursework" \
      org.opencontainers.image.revision=${COMMIT_SHA} \
      org.opencontainers.image.created=${BUILD_DATE}

CMD ["./main"]