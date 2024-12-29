package redis

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/s21platform/school-service/internal/config"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Repository struct {
	conn *redis.Client
}

func New(cfg *config.Config) *Repository {
	redisPort := cfg.Cache.Port
	redisHost := cfg.Cache.Host
	redisAddr := fmt.Sprintf("%s:%s", redisHost, redisPort)

	rdb := redis.NewClient(&redis.Options{
		Addr:         redisAddr,
		Password:     "",
		DB:           0,
		MinIdleConns: 2,
	})

	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		log.Fatal(err)
	}

	return &Repository{conn: rdb}
}

func (r *Repository) Get(ctx context.Context) (string, error) {

	keys, err := r.conn.Keys(ctx, "*").Result()
	if err != nil {
		fmt.Println("Ошибка получения ключей:", err)
		return "", err
	}

	if len(keys) > 0 {
		randomKey := keys[rand.Intn(len(keys))]

		val, err := r.conn.Get(ctx, randomKey).Result()
		if errors.Is(err, redis.Nil) {
			log.Printf("Ключ %s не найден\n", randomKey)
			return "", err
		} else if err != nil {
			log.Printf("Ошибка при получении ключа %s: %v\n", randomKey, err)
			return "", err
		}
		return val, nil
	} else {
		return "", status.Errorf(codes.Unknown, "В Redis нет ключей")
	}
}

func (r *Repository) Set(ctx context.Context, key string, value string, expiration time.Duration) error {

	err := r.conn.Set(ctx, key, value, expiration).Err()
	if err != nil {
		log.Printf("Ошибка при установке ключа %s: %v", key, err)
		return err
	}

	return nil
}
