# ========= BUILD STAGE =============
FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o checkly-api ./cmd/api

# ========= FINAL STAGE =============
FROM alpine:latest
WORKDIR /
COPY --from=builder /app/checkly-api ./
CMD ["./checkly-api"]