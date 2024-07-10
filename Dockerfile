# Используем официальный образ Golang как базовый
FROM golang:1.22-alpine AS build

# Устанавливаем рабочую директорию внутри контейнера
WORKDIR /app

# Копируем исходный код
COPY . .

# Собираем бинарный файл
RUN go build -o testhub-spec-uni

# Используем минимальный образ для конечного контейнера
FROM debian:bookworm-slim

# Создаем рабочую директорию
WORKDIR /app

# Копируем бинарный файл из предыдущего этапа
COPY --from=build /app/testhub-spec-uni .

# Определяем порт, который будет использоваться
EXPOSE 8080

# Устанавливаем переменные окружения
ENV DB_DRIVER=postgres
ENV DB_USER=postgres
ENV DB_PASSWORD=1001
ENV DB_HOST=db
ENV DB_PORT=5432
ENV DB_NAME=postgres

# Запускаем приложение
CMD ["./testhub-spec-uni"]
