package api

import "github.com/gin-gonic/gin"

func srsPage(c *gin.Context) {
	c.HTML(200, "srs.html", nil)
}
