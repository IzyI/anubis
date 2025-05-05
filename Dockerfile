FROM golang:1.23-alpine AS builder

# Устанавливаем зависимости
RUN apk add --no-cache git

# Рабочая директория
WORKDIR /app

# Копируем go.mod и go.sum (если есть)
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходный код
COPY . .

# Собираем приложение
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/my_app ./cmd/server/main.go

# Финальный образ (минимальный)
FROM alpine:latest

# Копируем бинарник из builder
COPY --from=builder /app/my_app /my_app

# Открываем порт (если приложение слушает 8080)
# EXPOSE 1010

# Запускаем приложение
CMD ["/my_app"]