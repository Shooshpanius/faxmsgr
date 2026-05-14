// Точка входа сервера FAX Messenger.
// Инициализирует подключения к БД, Redis, S3, запускает HTTP-сервер.
package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"

	"faxmsgr/server/internal/archive"
	"faxmsgr/server/internal/auth"
	"faxmsgr/server/internal/chat"
	"faxmsgr/server/internal/middleware"
	"faxmsgr/server/internal/storage"
	"faxmsgr/server/internal/user"
	"faxmsgr/server/internal/ws"
)

func main() {
	// Инициализация структурированного логгера
	logger, err := zap.NewProduction()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to init logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync() //nolint:errcheck

	// Подключение к PostgreSQL
	pgDSN := mustEnv("FAX_DATABASE_URL")
	pool, err := storage.NewPostgres(context.Background(), pgDSN)
	if err != nil {
		logger.Fatal("failed to connect to postgres", zap.Error(err))
	}
	defer pool.Close()

	// Подключение к PostgreSQL (архив — Закон Яровой)
	archiveDSN := mustEnv("FAX_ARCHIVE_DATABASE_URL")
	archivePool, err := storage.NewPostgres(context.Background(), archiveDSN)
	if err != nil {
		logger.Fatal("failed to connect to archive postgres", zap.Error(err))
	}
	defer archivePool.Close()

	// Подключение к Redis
	redisAddr := getEnv("FAX_REDIS_ADDR", "localhost:6379")
	rdb := storage.NewRedis(redisAddr, os.Getenv("FAX_REDIS_PASSWORD"), 0)

	// JWT-секрет
	jwtSecret := mustEnv("FAX_JWT_SECRET")

	// WebSocket-хаб
	hub := ws.NewHub(logger)
	go hub.Run()

	// Сервис архивирования (Закон Яровой)
	_ = archive.NewService(archivePool, logger)

	// Сервисы
	authSvc := auth.NewService(pool, rdb, jwtSecret, logger)
	userSvc := user.NewService(pool, logger)
	chatSvc := chat.NewService(pool, logger)

	// Роутер
	r := chi.NewRouter()
	r.Use(chimiddleware.RequestID)
	r.Use(chimiddleware.RealIP)
	r.Use(chimiddleware.Recoverer)

	// Публичные маршруты
	r.Post("/auth/request-code", auth.MakeRequestCodeHandler(authSvc))
	r.Post("/auth/verify-code", auth.MakeVerifyCodeHandler(authSvc))
	r.Post("/auth/refresh", auth.MakeRefreshHandler(authSvc))

	// Защищённые маршруты
	r.Group(func(r chi.Router) {
		r.Use(middleware.JWT(jwtSecret, logger))
		r.Get("/users/profile", user.MakeGetProfileHandler(userSvc))
		r.Put("/users/profile", user.MakePutProfileHandler(userSvc))
		r.Post("/chats", chat.MakeCreateChatHandler(chatSvc))
		r.Get("/chats", chat.MakeListChatsHandler(chatSvc))
	})

	// WebSocket
	r.Get("/ws", ws.MakeHandler(hub, jwtSecret, logger))

	addr := getEnv("FAX_SERVER_ADDR", ":8080")
	srv := &http.Server{
		Addr:         addr,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Запуск с graceful shutdown
	go func() {
		logger.Info("server starting", zap.String("addr", addr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("server error", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("shutting down server")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("server forced to shutdown", zap.Error(err))
	}
}

// mustEnv возвращает значение переменной окружения или завершает процесс.
func mustEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		fmt.Fprintf(os.Stderr, "required env var %s is not set\n", key)
		os.Exit(1)
	}
	return v
}

// getEnv возвращает значение переменной окружения или значение по умолчанию.
func getEnv(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}
