# Стадия сборки
FROM golang:1.24 AS builder

WORKDIR /app

# 1. Копируем только файлы зависимостей
COPY go.mod go.sum ./

# 2. Скачиваем зависимости
RUN go mod download

# 3. Теперь копируем весь код проекта
COPY . .

# 4. Сборка бинарника
RUN go build -o cityserver ./cmd/

# Финальный минимальный образ
FROM debian:bullseye-slim

WORKDIR /app

COPY --from=builder /app/cityserver .
COPY config ./config

EXPOSE 8080

CMD ["./cityserver"]
