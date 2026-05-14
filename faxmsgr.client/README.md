# faxmsgr.client — Веб-клиент FAX Messenger

Веб-клиент мессенджера FAX. Стек: **React 19**, **TypeScript**, **Vite**.

## Требования

- Node.js 20+
- npm 10+

## Разработка

```bash
npm install
npm run dev
```

Клиент запустится на `https://localhost:59551` и будет проксировать API-запросы на сервер ASP.NET Core.

## Сборка для продакшена

```bash
npm run build
```

Артефакты попадают в `dist/` и раздаются сервером через Nginx.

## Линтинг

```bash
npm run lint
```

