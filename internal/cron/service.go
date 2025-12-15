package cron

import (
	"github.com/gocraft/work"
	"github.com/gomodule/redigo/redis"
	"github.com/laciferin2024/url-shortner.go/internal/genesis"
)

type service struct {
	*genesis.Service
	redisPool *redis.Pool
	enqueuer  *work.Enqueuer

	pool *work.WorkerPool
}
