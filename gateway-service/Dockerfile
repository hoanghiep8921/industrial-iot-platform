# ============================================================
# Gateway Service - Dockerfile
# ============================================================
FROM golang:1.22-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /gateway ./cmd/gateway/

FROM alpine:3.19
RUN apk add --no-cache ca-certificates tzdata
ENV TZ=Asia/Ho_Chi_Minh

COPY --from=builder /gateway /usr/local/bin/gateway

EXPOSE 8081
ENTRYPOINT ["gateway"]
