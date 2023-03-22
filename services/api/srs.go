package api

import (
	"github.com/gin-gonic/gin"
	"monitoring-system/config"
)

func srsPage(c *gin.Context) {
	c.HTML(200, "srs.html", config.GitToken)
}
