package cron

import (
	"github.com/gocraft/work"
	"github.com/gomodule/redigo/redis"
	"github.com/laciferin2024/url-shortner.go/internal/cache"
	"github.com/laciferin2024/url-shortner.go/internal/genesis"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Invoke(initCron),
)

type in struct {
	fx.In
	Service       *genesis.Service
	CacheServices cache.Services
}

type out struct {
	fx.Out
}

func initCron(i in) (o out) {

	redisPool := i.CacheServices.GetPool().(*redis.Pool)

	cron := &service{
		i.Service,
		redisPool,
		work.NewEnqueuer("gocraft", redisPool),
		work.NewWorkerPool(service{}, 10, "gocraft", redisPool),
	}

	go cron.init()

	return
}
