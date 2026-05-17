# ==========================================
# Stage 1: Builder
# ==========================================
FROM golang:1.25.4-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/microservice ./cmd/simple-microservice/main.go

# ==========================================
# Stage 2: Final Image
# ==========================================
FROM alpine:latest

WORKDIR /app

RUN apk --no-cache add ca-certificates tzdata && \
    addgroup -S appgroup && \
    adduser -S appuser -G appgroup

COPY --chown=appuser:appgroup --from=builder /app/bin/microservice .

USER appuser

EXPOSE 8080

ENTRYPOINT ["./microservice"]