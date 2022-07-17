package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func getAdminPage(c *gin.Context) {
	c.HTML(http.StatusOK, "admin.html", nil)
}
