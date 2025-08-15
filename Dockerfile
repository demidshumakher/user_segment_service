# Build stage
FROM golang:1.25 AS builder

WORKDIR /app

# Кэшируем зависимости
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходники
COPY . .

# Сборка статического бинарника
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server ./main.go

# Run stage
FROM alpine:3.20

# Добавим сертификаты
RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

# Непривилегированный пользователь
RUN adduser -D -g '' appuser
USER appuser

COPY --from=builder /app/server /app/server

# Значение по умолчанию; переопределяется APP_PORT из .env
ENV APP_PORT=8080

EXPOSE 8080

CMD ["/app/server"]
