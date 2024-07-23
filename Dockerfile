# Используем официальный образ Golang для сборки
FROM golang:1.22-alpine as builder

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем файлы проекта в рабочую директорию
COPY . .

# Устанавливаем зависимости и синхронизируем папку vendor
RUN go mod tidy
RUN go mod vendor
RUN go build -o main .

# Используем минимальный образ Debian для запуска
FROM debian:bookworm-slim

# Устанавливаем необходимые пакеты
RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

# Копируем бинарник из стадии сборки и устанавливаем права
WORKDIR /app
COPY --from=builder /app/main /app/main
COPY --from=builder /app/conf /app/conf
COPY --from=builder /app/swagger /app/swagger
RUN chmod +x /app/main

# Команда для запуска приложения
CMD ["/app/main"]
