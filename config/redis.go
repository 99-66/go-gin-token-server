package config

import (
	"context"
	"errors"
	"fmt"
	"github.com/caarlos0/env"
	"github.com/go-redis/redis/v8"
)

var CTX = context.Background()
var REDIS *redis.Client

type RedisConfig struct {
	Host string `env:"REDIS_HOST"`
	Port int `env:"REDIS_PORT"`
	DB int `env:"REDIS_DB"`
	Password string `env:"REDIS_PASSWORD"`
}

// InitRedis Redis connection 을 생성하여 반환한다
func InitRedis() (*redis.Client, error){
	redisConfig := RedisConfig{}
	if err := env.Parse(&redisConfig); err != nil {
		return nil, errors.New("could not load redis configuration")
	}
	dsn := fmt.Sprintf("%s:%d", redisConfig.Host, redisConfig.Port)

	r := redis.NewClient(&redis.Options{
		Addr: dsn,
		Password: redisConfig.Password,
		DB: redisConfig.DB,
	})

	return r, nil
}