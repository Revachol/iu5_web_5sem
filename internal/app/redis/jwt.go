package redis

import (
	"context"
	"time"
)

const jwtPrefix = "jwt."

// getJWTKey формирует уникальный ключ Redis для хранения JWT токена.
func getJWTKey(token string) string {
	return servicePrefix + jwtPrefix + token
}

// WriteJWTToBlacklist добавляет JWT в блеклист с указанным TTL.
// Используется для мгновенной деактивации токена (например, при logout).
func (c *Client) WriteJWTToBlacklist(ctx context.Context, jwtStr string, ttl time.Duration) error {
	// Контекст с таймаутом для предотвращения подвисаний
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	return c.client.Set(ctx, getJWTKey(jwtStr), true, ttl).Err()
}

// CheckJWTInBlacklist проверяет, находится ли JWT в блеклисте.
// Если токен отсутствует, возвращается redis.Nil, что считается нормой.
func (c *Client) CheckJWTInBlacklist(ctx context.Context, jwtStr string) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	return c.client.Get(ctx, getJWTKey(jwtStr)).Err()
}
