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

ARG COMMIT_SHA=unknown
ARG BUILD_DATE=unknown
LABEL org.opencontainers.image.source="https://github.com/Qulive-edu/ArchitectureCoursework" \
      org.opencontainers.image.revision=${COMMIT_SHA} \
      org.opencontainers.image.created=${BUILD_DATE}

WORKDIR /app
CMD ["./main"]