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
  (counterId: number, action: 'hit', url: string, options?: YmHitOptions): void;
  (counterId: number, action: 'reachGoal', target: string, params?: YmParams, callback?: () => void, ctx?: object): void;
  (counterId: number, action: 'setUserID', userId: string): void;
  (counterId: number, action: 'userParams', params: YmParams): void;
  (counterId: number, action: string, ...args: unknown[]): void;
  a?: unknown[][];
  l?: number;
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
