package cache

import (
	"github.com/gomodule/redigo/redis"
	"github.com/laciferin2024/url-shortner.go/enums"
	"github.com/laciferin2024/url-shortner.go/internal/genesis"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(
		newServices,
	))

type in struct {
	fx.In

	*genesis.Service
}

type out struct {
	fx.Out
	Services
}

func newServices(i in) (o out) {
	cache := &cache{
		i.Service,
	}
	o = out{
		Services: newRedis(cache),
	}
	return
}

func newRedis(cache *cache) Services {
	conf := cache.Conf
	return &redisStore{
		pool: &redis.Pool{
			MaxActive: 50,
			MaxIdle:   5,
			Wait:      true,
			Dial: func() (redis.Conn, error) {
				opts := []redis.DialOption{}

				if pass := conf.GetString(enums.REDIS_MASTER_PASSWORD); pass != "" {
					opts = append(opts, redis.DialPassword(pass))
				}

				return redis.Dial("tcp",
					conf.GetString(enums.REDIS_SERVER),
					opts...)
			},
		},
		cache: cache,
	}
}
