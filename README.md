# FAX Messenger

Мессенджер с веб-клиентом и мобильным приложением.

## Структура репозитория

| Папка | Описание |
|---|---|
| `faxmsgr.Server/` | Серверная часть — Go 1.24, REST API + WebSocket, PostgreSQL, Redis |
| `faxmsgr.client/` | Веб-клиент — React 19, TypeScript, Vite |
| `faxmsgr.FlutterMobile/` | Мобильное приложение — Flutter |

## Быстрый старт

### Сервер + веб-клиент

```bash
# Из корня репозитория
docker-compose up --build
```

Веб-клиент будет доступен на `http://localhost:80`.

> **Примечание:** `faxmsgr.Server/docker-compose.yml` предназначен для локальной разработки — запуск из директории `faxmsgr.Server/`.
> Корневой `docker-compose.yml` используется для production-деплоя из корня репозитория.

### Мобильное приложение

```bash
cd faxmsgr.FlutterMobile
flutter run
```

## Документация

- [`faxmsgr.Server/README.md`](faxmsgr.Server/README.md) — серверная часть
- [`faxmsgr.client/README.md`](faxmsgr.client/README.md) — веб-клиент
- [`faxmsgr.FlutterMobile/README.md`](faxmsgr.FlutterMobile/README.md) — мобильное приложение
- [`.github/copilot-instructions.md`](.github/copilot-instructions.md) — соглашения и архитектура для Copilot
