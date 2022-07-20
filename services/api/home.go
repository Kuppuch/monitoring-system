package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func getMainPage(c *gin.Context) {
	c.HTML(http.StatusOK, "home.html", nil)
}
