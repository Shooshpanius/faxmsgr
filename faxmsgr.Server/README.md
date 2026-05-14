# faxmsgr.Server — Серверная часть FAX Messenger

REST API и real-time сервер мессенджера FAX. Стек: **Go 1.24**, **PostgreSQL 16**, **Redis 7**, **MinIO (S3)**, **WebSocket**.

## Требования

- Go 1.24+
- Docker & Docker Compose

## Быстрый старт (одной командой)

```bash
cd faxmsgr.Server
docker-compose up --build
```

Сервер будет доступен на `http://localhost:8080`.

В режиме разработки `docker-compose` автоматически подхватывает `docker-compose.override.yml`,
который пробрасывает порты внутренних сервисов на хост (PostgreSQL, Redis, MinIO) —
удобно для подключения из DBeaver, redis-cli, MinIO Console и т.д.
В production `docker-compose.override.yml` не используется.

## Запуск в режиме разработки

```bash
cd faxmsgr.Server
go run ./cmd/server
```

## Переменные окружения

| Переменная | Описание | По умолчанию |
|---|---|---|
| `FAX_DATABASE_URL` | DSN PostgreSQL (основная БД) | — |
| `FAX_ARCHIVE_DATABASE_URL` | DSN PostgreSQL (архив, Закон Яровой) | — |
| `FAX_REDIS_ADDR` | Адрес Redis (`host:port`) | `localhost:6379` |
| `FAX_REDIS_PASSWORD` | Пароль Redis | — |
| `FAX_JWT_SECRET` | Секрет для подписи JWT-токенов | — |
| `FAX_SERVER_ADDR` | Адрес и порт сервера | `:8080` |
| `FAX_S3_ENDPOINT` | Адрес MinIO/S3 | — |
| `FAX_S3_ACCESS_KEY` | Access Key для S3 | — |
| `FAX_S3_SECRET_KEY` | Secret Key для S3 | — |
| `FAX_S3_BUCKET` | Имя бакета для медиафайлов | — |

## Миграции БД

Миграции применяются через [goose](https://github.com/pressly/goose):

```bash
# Основная БД
goose -dir migrations postgres "$DATABASE_URL" up

# Архивная БД
goose -dir migrations postgres "$ARCHIVE_DATABASE_URL" up
```

## REST API

| Метод | Путь | Описание |
|---|---|---|
| POST | `/auth/request-code` | Запрос OTP на номер телефона |
| POST | `/auth/verify-code` | Проверка OTP, получение токенов |
| POST | `/auth/refresh` | Обновление токенов |
| GET | `/users/profile` | Получение профиля (🔒) |
| PUT | `/users/profile` | Обновление профиля (🔒) |
| POST | `/chats` | Создание чата (🔒) |
| GET | `/chats` | Список чатов (🔒) |
| GET | `/ws?token=...` | WebSocket-соединение (🔒) |

🔒 — требует `Authorization: Bearer <access_token>`

## WebSocket-протокол

Сообщения в формате JSON:

```json
{ "type": "MSG_SEND", "to": "<userID>", "chat_id": "<chatID>", "payload": "текст" }
{ "type": "MSG_DELIVERED", "from": "<userID>", "chat_id": "<chatID>" }
{ "type": "MSG_READ", "from": "<userID>", "chat_id": "<chatID>" }
{ "type": "STATUS_ONLINE" }
{ "type": "STATUS_OFFLINE" }
```

## Сборка Docker-образа

```bash
docker build -t faxmsgr-server .
```
