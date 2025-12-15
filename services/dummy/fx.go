package dummy

import (
	"github.com/laciferin2024/url-shortner.go/internal/genesis"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(
		newServices,
	),
)

type in struct {
	fx.In
	*genesis.Service
}

type out struct {
	fx.Out

	Dummy Services // `name:"dummy"`
}

func newServices(i in) (o out) {
	o = out{
		Dummy: newDummy(i.Service),
	}
	return
}

func newDummy(genesis *genesis.Service) Services {
	return &service{
		genesis,
	}
}
