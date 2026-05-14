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
cd faxmsgr.Server
docker-compose up --build
```

Сервер будет доступен на `http://localhost:8080`.

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
