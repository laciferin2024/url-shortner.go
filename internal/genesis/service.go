package genesis

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"time"

	"github.com/laciferin2024/url-shortner.go/utils"
	"github.com/parnurzeal/gorequest"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"go.uber.org/multierr"
)

type Service struct {
	Log  logrus.Logger //logs can be customized for every Service
	Conf *viper.Viper
}

func (s *Service) CreateContext() (ctx context.Context, cancel context.CancelFunc) {
	ctx, cancel = context.WithCancel(context.TODO())
	return
}

// ReadResponse Provide pointer
func (s *Service) ReadResponse(resp *http.Response, custom_type interface{}) (res json.RawMessage, err error) {
	res, err = utils.ReadResonse(resp)
	if err != nil {
		return
	}
	utils.ConvertJSONToGoType(res, custom_type)
	return
}

func (s *Service) MakeRequest(u url.URL) (res *http.Response, err error) {
	client := http.Client{Timeout: time.Duration(5) * time.Second}
	res, err = client.Get(u.RequestURI())
	return
}

func (s *Service) PingWebhook(url string, req interface{}) (res *http.Response, err error) {
	_, _, errs := gorequest.New().Post(url).Send(req).End()
	err = multierr.Combine(errs...)
	return
}
