package redis

import (
	"context"
	"fmt"
	"github.com/Revachol/iu5_web_5sem/internal/app/config"
	"github.com/go-redis/redis/v8"
	//"strconv"
	"time"
)

const servicePrefix = "stars-catalog."

// Client представляет Redis-клиент с конфигурацией и подключением.
type Client struct {
	cfg    config.RedisConfig
	client *redis.Client
}

// New создаёт новый Redis-клиент и проверяет соединение через PING.
// Возвращает ошибку, если соединение не удалось установить.
func New(ctx context.Context, cfg config.RedisConfig) (*Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:        cfg.Host + ":" + "6379",
		Username:    cfg.User,
		Password:    cfg.Password,
		DB:          0,
		DialTimeout: cfg.DialTimeout,
		ReadTimeout: cfg.ReadTimeout,
	})

	// Контекст с таймаутом для PING
	pingCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	if _, err := rdb.Ping(pingCtx).Result(); err != nil {
		return nil, fmt.Errorf("cannot ping redis: %w", err)
	}

	return &Client{
		cfg:    cfg,
		client: rdb,
	}, nil
}

// Close закрывает соединение с Redis.
// Возвращает ошибку, если закрытие не удалось.
func (c *Client) Close() error {
	return c.client.Close()
}
