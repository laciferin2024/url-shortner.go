package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/laciferin2024/url-shortner.go/entities"
)

func (h *AppHandler) Home(c *gin.Context) {
	var (
		res = entities.GenericResponse{}
		//err error
	)

	//defer h.handleResponse(c, &res, &err)

	res.Success = true
	res.Message = "App service is up and running"
	c.JSON(http.StatusOK, res)
}

type urlReq struct {
	Url string `json:"url" binding:"required"`
}

func (h *AppHandler) ShortenUrl(c *gin.Context) {

	var u1 urlReq

	if err := c.ShouldBind(&u1); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"error": "request does not match the accepted format",
		})
		return

	}

	shortUrl, err := h.appServices.ShortenUrl(c.Request.Context(), u1.Url)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	host := c.Request.Host
	path := c.Request.URL.RequestURI()

	c.JSON(http.StatusOK, gin.H{
		"short-url": fmt.Sprintf("http://%s%s/%s", host, path, shortUrl),
		"Url":       u1.Url,
	})

}

func (h *AppHandler) RetrieveUrl(c *gin.Context) {

	shortUrl := c.Param("surl")

	origUrl, err := h.appServices.RetrieveOriginalUrl(shortUrl)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return

	}

	c.Redirect(http.StatusPermanentRedirect, origUrl)
}

func (h *AppHandler) ListUrls(c *gin.Context) {
	page := 1
	pageSize := 10

	// Parse query params (simple implementation)
	// In a real app, use binding or strconv
	// For now, hardcoded defaults or simple logic if needed
	// But let's try to bind if possible, or just use defaults

	urls, err := h.appServices.ListUrls(c.Request.Context(), pageSize, (page-1)*pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, urls)
}
