package redis

import (
	"github.com/go-redis/redis"
)

type RedisStore struct {
	DB *redis.Client
}
