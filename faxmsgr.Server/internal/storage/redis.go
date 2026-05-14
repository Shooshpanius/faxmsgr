package storage

import (
	"github.com/redis/go-redis/v9"
)

// NewRedis создаёт клиент Redis.
// addr — адрес сервера (например, "localhost:6379"),
// password — пароль (пустая строка если не используется),
// db — номер базы данных Redis.
func NewRedis(addr, password string, db int) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
}
