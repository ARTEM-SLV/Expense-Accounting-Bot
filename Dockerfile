# Базовый образ Go
FROM golang:1.22

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем go.mod и go.sum для установки зависимостей
COPY go.mod go.sum ./
RUN go mod download

# Копируем остальные файлы проекта
COPY . .

# Сборка приложения, указываем путь к файлу main.go
RUN go build -o main ./cmd

# Указываем порт (если нужно)
EXPOSE 3000

# Запуск приложения
CMD ["./main"]