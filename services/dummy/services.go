package dummy

import "github.com/laciferin2024/url-shortner.go/entities"

type Services interface {
	Dummy(dumb1 *entities.Dummy)
}

func (s *service) Dummy(dumb1 *entities.Dummy) {
	dumb1.Dummy = "dumber"
	return
}
