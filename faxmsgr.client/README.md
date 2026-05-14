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

Клиент запустится на `http://localhost:5173` и будет проксировать API-запросы на Go-сервер.

## Сборка для продакшена

```bash
npm run build
```

Артефакты попадают в `dist/` и раздаются сервером через Nginx.

## Переменные окружения

Создайте файл `.env.local` в корне клиента:

| Переменная | Описание | Обязательная |
|---|---|---|
| `VITE_FAX_API_URL` | URL Go-сервера для проксирования в dev-режиме (по умолчанию `http://localhost:8080`) | Нет |
| `VITE_FAX_YM_COUNTER_ID` | Номер счётчика Яндекс Метрики | Нет |

Если `VITE_FAX_YM_COUNTER_ID` не задан, Яндекс Метрика не подключается.

## Яндекс Метрика

Для отслеживания событий используйте хелперы из `src/lib/metrika.ts`:

```ts
import { reachGoal, hitPage } from './lib/metrika'

// Отправить цель
reachGoal('registration')

// Зарегистрировать переход (SPA)
hitPage(window.location.pathname)
```

## Линтинг

```bash
npm run lint
```

