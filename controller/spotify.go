package controller

import (
	"github.com/gin-gonic/gin"
	"golang-songs/service"
	"net/http"
)

func Login(c *gin.Context) {
	c.Redirect(http.StatusTemporaryRedirect, service.GetRedirectURL())
}