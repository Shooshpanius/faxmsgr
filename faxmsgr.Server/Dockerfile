# Сборка Go-бинарника
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o server ./cmd/server

# Минимальный образ для запуска
FROM gcr.io/distroless/static-debian12
COPY --from=builder /app/server /server
COPY --from=builder /app/migrations /migrations
EXPOSE 8080
ENTRYPOINT ["/server"]
