package handlers

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/laciferin2024/url-shortner.go/entities"
	"github.com/laciferin2024/url-shortner.go/er"
	"github.com/laciferin2024/url-shortner.go/services/app"
	"github.com/laciferin2024/url-shortner.go/services/dummy"
	"github.com/laciferin2024/url-shortner.go/utils"
	"github.com/sirupsen/logrus"

	"github.com/spf13/viper"
)

type Handler struct {
	conf *viper.Viper
	log  *logrus.Logger
}

type DummyHandler struct {
	Handler

	dummyServices dummy.Services
}

type HomeHandler struct {
	Handler
}

type AppHandler struct {
	Handler

	appServices app.Services
}

func Bind(c *gin.Context, req interface{}) (statusCode int, err error) {
	if c.Request.Method == "GET" {
		queries := c.Request.URL.Query()

		if err = utils.ConvertMapToAny(queries, req); err != nil {
			c.Error(err).SetType(er.ValidationError)
			statusCode = http.StatusUnprocessableEntity
		}
		return
	}
	if err = c.ShouldBindBodyWith(req, binding.JSON); err != nil {
		//if err = c.ShouldBind(&req); err != nil {
		c.Error(err).SetType(er.ValidationError)
		statusCode = http.StatusUnprocessableEntity
	}
	return
}

func (h *Handler) createContext() (ctx context.Context, cancel context.CancelFunc) {
	ctx, cancel = context.WithCancel(context.TODO())
	return
}

func (h *Handler) getSystemHost(c *gin.Context) (u string) {
	scheme := "http"
	if c.Request.TLS != nil {
		scheme = "https"
	}
	//c.Request.Proto //can be used as well to differentiate http versions
	u = fmt.Sprintf("%s://%s", scheme, c.Request.Host)
	return
}

func (h *Handler) error(c *gin.Context, err error, errType gin.ErrorType) {
	if err != nil {
		c.Error(err).SetType(errType)
	}
}

// addErrors use with defer
func (h *Handler) addErrors(c *gin.Context, err error) {
	if err != nil {
		c.Error(err)
	}
}

func (h *Handler) addError(c *gin.Context, err error_) {
	c.Error(gin.Error{
		errors.New(err.Message),
		gin.ErrorType(err.StatusCode),
		nil,
	})
}

func (h *Handler) handleResponse(c *gin.Context, res *entities.GenericResponse, finalError *error) {

	if res == nil && finalError != nil {
		res = &entities.GenericResponse{}
	} else {
		return
	}

	err := c.Errors.Last() //not working while using backup feature
	if err != nil && *finalError != nil {
		if res.StatusCode == 0 {
			res.StatusCode = http.StatusServiceUnavailable
		}
		res.Data = err
		res.Message = err.Error()
		res.Success = false
		c.AbortWithStatusJSON(res.StatusCode, res)
		return
	}
	if res.StatusCode == 0 {
		res.StatusCode = http.StatusCreated
	}
	c.JSON(res.StatusCode, res)
}

func (h *Handler) handleError(c *gin.Context) {

	err := c.Errors.Last() //not working while using backup feature

	if err == nil {
		return
	}

	if err.Type == 0 {
		err.Type = http.StatusNotFound
	}

	res := entities.GenericResponse{
		Message:    err.Error(),
		Data:       err,
		Success:    false,
		StatusCode: int(err.Type),
	}

	c.AbortWithStatusJSON(res.StatusCode, res)
	return
}

type error_ struct {
	StatusCode int
	Message    string
}

// use statuscode 0 for default
func newError(message string, statusCode int) error_ {
	return error_{
		statusCode,
		message,
	}
}
