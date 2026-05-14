# faxmsgr.Server — Серверная часть FAX Messenger

REST API и real-time сервер мессенджера FAX. Стек: **ASP.NET Core 10**, **C#**, **SignalR**, **PostgreSQL + EF Core**.

## Требования

- .NET 10 SDK
- PostgreSQL 15+

## Разработка

```bash
dotnet run
```

В режиме разработки автоматически запускается веб-клиент через SpaProxy.  
OpenAPI (Scalar) доступен по адресу `https://localhost:<port>/openapi/v1.json`.

## Переменные окружения

| Переменная | Описание |
|---|---|
| `ConnectionStrings__Default` | Строка подключения к PostgreSQL |
| `Email__SmtpHost` | SMTP-хост для отправки писем |

Для локальной разработки используйте `dotnet user-secrets`.

## Миграции БД

```bash
dotnet ef migrations add <MigrationName>
dotnet ef database update
```

## Сборка Docker-образа

```bash
docker build -t faxmsgr-server .
```
