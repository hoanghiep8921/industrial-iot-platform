# ============================================================
# Alarm Service - Dockerfile
# ============================================================
FROM golang:alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /alarm ./cmd/alarm/

FROM alpine:3.19
RUN apk add --no-cache ca-certificates tzdata
ENV TZ=Asia/Ho_Chi_Minh

COPY --from=builder /alarm /usr/local/bin/alarm
COPY migrations/ /migrations/

EXPOSE 8086
ENTRYPOINT ["alarm"]
