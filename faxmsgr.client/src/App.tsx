import { useMemo } from 'react';
import './App.css';
import { LAST_PR_DATE, LAST_PR_NUMBER } from './version';

const ALL_TAGLINES = [
    'Новые мессенджеры падают. Факс — никогда.',
    'В 1990-х факс не зависал и не просил обновиться.',
    'Технологии ушли вперёд. Надёжность осталась в прошлом. Мы её вернули.',
    'Пока у них очередное «критическое обновление», у нас — просто работает.',
    'Факс не знал, что такое «потеря соединения». Мы тоже.',
    'Старая школа. Без багов, без утечек данных, без «что-то пошло не так».',
    'Прогресс — это хорошо. Но надёжность факса мы взяли с собой.',
    'Новые технологии сложны и ненадёжны. Мы сделали наоборот.',
    'Раньше сообщение доходило всегда. Мы вспомнили, как это делалось.',
    'Современные мессенджеры забыли главное — доставлять сообщения. Мы не забыли.',
    'Факс работал в любую погоду. FAX messenger — тоже.',
    'Инновация — это не усложнить. Это вернуть то, что работало.',
];

function pickRandom<T>(arr: T[], count: number): T[] {
    const shuffled = [...arr].sort(() => Math.random() - 0.5);
    return shuffled.slice(0, count);
}

function App() {
    const taglines = useMemo(() => pickRandom(ALL_TAGLINES, 3), []);

    return (
        <div className="landing">
            <header className="landing-header">
                <div className="wip-badge">🚧 В разработке — дата выпуска пока неизвестна</div>
                <div className="version-tag">PR #{LAST_PR_NUMBER} · {LAST_PR_DATE}</div>
            </header>

            <main>
                <section className="hero">
                    <div className="hero-logo">
                        <span className="logo-fax">FAX</span>
                        <span className="logo-messenger">messenger</span>
                    </div>
                    <h1 className="hero-tagline">
                        Факс не умер.<br />Он просто стал быстрее.
                    </h1>
                    <p className="hero-sub">
                        ФАКСтически мгновенная доставка сообщений — без бумаги, без шума,
                        без лишнего.
                    </p>
                    <a href="#features" className="hero-cta">Узнать больше</a>
                </section>

                <section className="taglines">
                    {taglines.map((line) => (
                        <blockquote key={line}>"{line}"</blockquote>
                    ))}
                </section>

                <section className="features" id="features">
                    <h2>Что внутри</h2>
                    <div className="features-grid">
                        <div className="feature-card">
                            <span className="feature-icon">🎟️</span>
                            <h3>Только по инвайту</h3>
                            <p>
                                Регистрация — по приглашению. Без спама, без случайных людей.
                                Только те, кого вы позвали сами.
                            </p>
                        </div>
                        <div className="feature-card">
                            <span className="feature-icon">💬</span>
                            <h3>Текстовые сообщения</h3>
                            <p>
                                На старте — обмен текстом в реальном времени. Быстро,
                                надёжно и без задержек.
                            </p>
                        </div>
                        <div className="feature-card feature-card--soon">
                            <span className="feature-icon">🖼️</span>
                            <h3>Картинки и медиа</h3>
                            <p>
                                Фото, файлы и прочее — придут позже. Факс тоже не сразу
                                научился передавать цвет.
                            </p>
                            <span className="soon-label">Скоро</span>
                        </div>
                        <div className="feature-card feature-card--soon">
                            <span className="feature-icon">🔔</span>
                            <h3>Уведомления</h3>
                            <p>
                                Push-уведомления и мобильное приложение — в планах.
                                Следите за обновлениями.
                            </p>
                            <span className="soon-label">Скоро</span>
                        </div>
                    </div>
                </section>

                <section className="status-section">
                    <h2>Текущий статус</h2>
                    <p>
                        Проект активно разрабатывается. Многое ещё впереди — интерфейс,
                        мобильное приложение, медиафайлы. Дата запуска пока не объявлена,
                        но работа идёт.
                    </p>
                    <p className="status-hint">Следи за обновлениями — факс придёт, когда будет готов.</p>
                </section>
            </main>

            <footer className="landing-footer">
                <span>FAX messenger</span>
                <span className="footer-sep">·</span>
                <span>PR #{LAST_PR_NUMBER}</span>
                <span className="footer-sep">·</span>
                <span>{LAST_PR_DATE}</span>
            </footer>
        </div>
    );
}

export default App;