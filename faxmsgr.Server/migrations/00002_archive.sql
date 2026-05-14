-- +goose Up
-- Таблицы архивной БД (Закон Яровой)
-- Все сообщения и метаданные хранятся 3 года

CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Архив сообщений
CREATE TABLE IF NOT EXISTS archived_messages (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    message_id  UUID NOT NULL UNIQUE,
    sender_id   UUID NOT NULL,
    chat_id     UUID NOT NULL,
    body        TEXT,
    media_url   TEXT,
    sender_ip   TEXT,
    archived_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    expires_at  TIMESTAMPTZ NOT NULL
);

-- Партиционирование по дате архивирования (для эффективного удаления устаревших записей)
CREATE INDEX IF NOT EXISTS idx_archived_messages_expires ON archived_messages(expires_at);

-- Аудит-журнал входов
CREATE TABLE IF NOT EXISTS audit_logins (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id      UUID NOT NULL,
    phone        TEXT NOT NULL,
    ip           TEXT NOT NULL,
    logged_in_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    expires_at   TIMESTAMPTZ NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_audit_logins_expires ON audit_logins(expires_at);

-- +goose Down
DROP TABLE IF EXISTS audit_logins;
DROP TABLE IF EXISTS archived_messages;
