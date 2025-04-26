# Используем официальный Go-образ
FROM golang:1.21 AS builder

WORKDIR /app

# Копируем исходники
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Собираем сервер
RUN go build -o cityserver ./cmd/cityserver

# Финальный образ
FROM debian:bullseye-slim

WORKDIR /app

# Копируем бинарник и конфиг
COPY --from=builder /app/cityserver .
COPY config ./config

# Указываем порт (если нужно)
EXPOSE 8080

CMD ["./cityserver"]
