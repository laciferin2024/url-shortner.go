package cache

import (
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/laciferin2024/url-shortner.go/internal/genesis"
)

type cache struct {
	*genesis.Service
}

type redisStore struct {
	pool              *redis.Pool
	defaultExpiration time.Duration
	*cache
}
