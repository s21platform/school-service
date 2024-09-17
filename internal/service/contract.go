package service

import (
	"context"
	"time"
)

type RedisR interface {
	Get(ctx context.Context) (string, error)
	Set(ctx context.Context, key string, value string, expiration time.Duration) error
}
