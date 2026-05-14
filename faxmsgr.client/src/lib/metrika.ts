// Типы для Яндекс Метрики
type YmHitOptions = {
  callback?: () => void;
  ctx?: object;
  params?: object;
  referer?: string;
  title?: string;
};

type YmParams = Record<string, unknown>;

type YmFunction = {
  (counterId: number, action: 'init', options: YmInitOptions): void;
  (counterId: number, action: 'hit', url: string, options?: YmHitOptions): void;
  (counterId: number, action: 'reachGoal', target: string, params?: YmParams, callback?: () => void, ctx?: object): void;
  (counterId: number, action: 'setUserID', userId: string): void;
  (counterId: number, action: 'userParams', params: YmParams): void;
  (counterId: number, action: string, ...args: unknown[]): void;
  a?: unknown[][];
  l?: number;
};

type YmInitOptions = {
  clickmap?: boolean;
  trackLinks?: boolean;
  accurateTrackBounce?: boolean;
  webvisor?: boolean;
  ecommerce?: boolean | string;
  defer?: boolean;
  params?: YmParams;
};

declare global {
  interface Window {
    ym?: YmFunction;
  }
}

const COUNTER_ID = import.meta.env.VITE_FAX_YM_COUNTER_ID
  ? Number(import.meta.env.VITE_FAX_YM_COUNTER_ID)
  : null;

/**
 * Инициализирует счётчик Яндекс Метрики.
 * Вызывается один раз при запуске приложения.
 * Если переменная окружения VITE_FAX_YM_COUNTER_ID не задана — ничего не делает.
 */
export function initMetrika(options: YmInitOptions = {}): void {
  if (!COUNTER_ID) return;

  // Заглушка для вызовов до загрузки скрипта
  if (!window.ym) {
    const stub = function (...args: unknown[]) {
      (stub.a = stub.a || []).push(args);
    } as unknown as YmFunction;
    stub.l = Date.now();
    window.ym = stub;
  }

  // Загружаем тег Метрики
  const script = document.createElement('script');
  script.type = 'text/javascript';
  script.async = true;
  script.src = 'https://mc.yandex.ru/metrika/tag.js';
  const firstScript = document.getElementsByTagName('script')[0];
  firstScript.parentNode?.insertBefore(script, firstScript);

  // Добавляем noscript-пиксель
  const noscript = document.createElement('noscript');
  const img = document.createElement('img');
  img.src = `https://mc.yandex.ru/watch/${COUNTER_ID}`;
  img.style.cssText = 'position:absolute;left:-9999px';
  img.alt = '';
  noscript.appendChild(img);
  document.body.appendChild(noscript);

  window.ym(COUNTER_ID, 'init', {
    clickmap: true,
    trackLinks: true,
    accurateTrackBounce: true,
    webvisor: true,
    ...options,
  });
}

/**
 * Отправляет событие-цель в Яндекс Метрику.
 * Если счётчик не инициализирован — вызов игнорируется.
 */
export function reachGoal(target: string, params?: YmParams): void {
  if (!COUNTER_ID || !window.ym) return;
  window.ym(COUNTER_ID, 'reachGoal', target, params);
}

/**
 * Регистрирует просмотр страницы (для SPA-навигации).
 */
export function hitPage(url: string, options?: YmHitOptions): void {
  if (!COUNTER_ID || !window.ym) return;
  window.ym(COUNTER_ID, 'hit', url, options);
}
