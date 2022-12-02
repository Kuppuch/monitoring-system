package public

import (
	"github.com/gin-gonic/gin"
	"monitoring-system/services/middleware"
	"net/http"
)

func GetUsers(c *gin.Context) {
	users := middleware.GetAllUsers()
	for i := range users {
		users[i].Password = ""
	}
	c.JSON(http.StatusOK, users)
}
