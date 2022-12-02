package public

import (
	"github.com/gin-gonic/gin"
	"monitoring-system/services/logging"
	"monitoring-system/services/middleware"
	"net/http"
	"strconv"
)

func GetProjects(c *gin.Context) {
	projectID, err := strconv.Atoi(c.DefaultQuery("id", "0"))
	if err != nil {
		logging.Print.Error(err)
		c.JSON(http.StatusBadRequest, gin.Error{
			Err:  err,
			Type: 0,
			Meta: "error id format",
		})
		return
	}

	if projectID > 0 {
		project := middleware.GetProjectByID(projectID)
		c.JSON(http.StatusOK, project)
		return
	}

	projects := middleware.GetProjects()
	c.JSON(http.StatusOK, projects)
}
