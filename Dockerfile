# Укажите базовый образ
FROM golang:1.22-alpine AS build

# Создайте рабочую директорию
WORKDIR /app

# Скопируйте go.mod и go.sum, а затем установите зависимости
COPY go.mod go.sum ./
RUN go mod download

# Скопируйте остальной исходный код
COPY . .

# Соберите бинарный файл
RUN go build -o main .

# Используйте минимальный образ для конечного контейнера
FROM alpine:latest

# Создайте директорию для приложения
WORKDIR /root/

# Скопируйте собранный бинарный файл из предыдущего этапа
COPY --from=build /app/main .
COPY --from=build /app/conf /root/conf

# Экспортируйте переменные окружения
ENV DB_DRIVER=postgres
ENV DB_USER=postgres
ENV DB_PASSWORD=1001
ENV DB_HOST=db
ENV DB_PORT=5432
ENV DB_NAME=postgres

# Запустите приложение
CMD ["./main"]
