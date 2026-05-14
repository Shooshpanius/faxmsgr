import { fileURLToPath, URL } from 'node:url';

import { defineConfig } from 'vite';
import plugin from '@vitejs/plugin-react';
import { env } from 'process';

const target = env.VITE_FAX_API_URL ?? 'http://localhost:8080';

// https://vitejs.dev/config/
export default defineConfig({
    plugins: [plugin()],
    resolve: {
        alias: {
            '@': fileURLToPath(new URL('./src', import.meta.url))
        }
    },
    server: {
        proxy: {
            // Проксирование API и WebSocket на Go-сервер в dev-режиме
            '^/(auth|users|chats|ws)(/|$)': {
                target,
                secure: false,
                ws: true,
            }
        },
        port: parseInt(env.DEV_SERVER_PORT || '5173'),
    }
})
