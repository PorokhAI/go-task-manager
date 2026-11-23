
# Этап 1. Сборка приложения
FROM golang:1.27.2-bookworm AS builder

WORKDIR /app

# Сначала копируем файлы зависимостей.
# Это позволяет Docker использовать кэш при повторных сборках.
COPY go.mod go.sum ./

RUN go mod download

# Копируем исходный код проекта.
COPY . .

# Собираем API и мигратор.
# CGO_ENABLED=0 позволяет получить статические бинарные файлы.
# GOOS=linux указывает целевую операционную систему.
RUN CGO_ENABLED=0 GOOS=linux go build \
    -trimpath \
    -ldflags="-s -w" \
    -o /out/api ./cmd/api

RUN CGO_ENABLED=0 GOOS=linux go build \
    -trimpath \
    -ldflags="-s -w" \
    -o /out/migrate ./cmd/migrate

RUN CGO_ENABLED=0 GOOS=linux go build \
    -trimpath \
    -ldflags="-s -w" \
    -o /out/healthcheck ./cmd/healthcheck

# Этап 2. Финальный образ API
FROM gcr.io/distroless/static-debian12:nonroot AS api

WORKDIR /app

COPY --from=builder /out/api ./api

COPY --from=builder /out/healthcheck ./healthcheck

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
    CMD ["/app/healthcheck"]

ENTRYPOINT ["/app/api"]


# Этап 3. Финальный образ мигратора
FROM gcr.io/distroless/static-debian12:nonroot AS migrate

WORKDIR /app

COPY --from=builder /out/migrate ./migrate
COPY --from=builder /app/migrations ./migrations

ENTRYPOINT ["/app/migrate"]
