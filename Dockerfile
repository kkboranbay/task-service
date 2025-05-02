FROM golang:1.24.2-alpine AS builder

# Установка зависимостей
RUN apk add --no-cache gcc musl-dev git

# Установка рабочей директории
WORKDIR /app

# Копирование файлов go.mod и go.sum
COPY go.mod go.sum ./

# Скачивание зависимостей
RUN go mod download

# Копирование исходного кода
COPY . .

# Сборка приложения
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o task-service ./cmd/api

# Финальный образ
FROM alpine:3.19

# Установка зависимостей
RUN apk --no-cache add ca-certificates tzdata

# Создание непривилегированного пользователя
RUN adduser -D -g '' appuser

# Копирование бинарного файла из builder
COPY --from=builder /app/task-service /usr/local/bin/

# Создание директории для миграций и копирование их
COPY --from=builder /app/migrations /migrations

# Настройка рабочей директории
WORKDIR /usr/local/bin

# Переключение на непривилегированного пользователя
USER appuser

# Определение команды запуска
CMD ["task-service"]