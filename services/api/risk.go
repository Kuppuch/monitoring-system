package api

import (
	"github.com/gin-gonic/gin"
	"monitoring-system/services/logging"
	"monitoring-system/services/middleware"
	"net/http"
)

func getRiskPage(c *gin.Context) {
	user, err := GetUserByToken(c)
	if err != nil {
		logging.Print.Error("error getting user from token ", err)
		c.JSON(http.StatusBadRequest, middleware.GetBadRequest())
		return
	}
	c.HTML(http.StatusOK, "risk.html", gin.H{"user": user})
}
