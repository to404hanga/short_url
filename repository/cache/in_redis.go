package cache

import "github.com/redis/go-redis/v9"

type InRedisCache struct {
	cmd redis.Cmdable
}
