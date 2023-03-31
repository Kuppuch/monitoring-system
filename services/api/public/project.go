package public

import (
	"encoding/json"
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

func InsertProject(c *gin.Context) {
	raw, err := c.GetRawData()
	if err != nil {
		logging.Print.Error(err)
		c.JSON(http.StatusBadRequest, gin.Error{
			Err:  err,
			Type: 0,
			Meta: "error by get raw data",
		})
		return
	}

	type PublicProject struct {
		middleware.Project
		Budget string `json:"budget"`
		ExtID  int    `json:"ext_id"`
	}

	pproject := PublicProject{}
	err = json.Unmarshal(raw, &pproject)
	if err != nil {
		logging.Print.Error("error unmarshal", err)
		c.JSON(http.StatusBadRequest, gin.Error{
			Err:  err,
			Type: 0,
			Meta: "error by unmarshal project",
		})
	}
	if len(pproject.Project.Name) < 1 {
		c.JSON(http.StatusBadRequest, middleware.GetBadRequest())
		return
	}
	project := pproject.Project
	project.StatusID = 1
	if rowAffected := project.InsertProject(); rowAffected == 0 {
		c.JSON(http.StatusBadRequest, middleware.GetBadRequest())
		return
	}
	budget := middleware.Budget{}
	budget.Name = pproject.Budget
	budget.ProjectID = int(project.ID)
	budget.ExtID = pproject.ExtID
	rowAffected := budget.Insert()
	if rowAffected == 0 {
		c.JSON(http.StatusBadRequest, middleware.GetBadRequest())
		return
	}
	c.JSON(http.StatusOK, middleware.GetSuccess())
}

func updateStatus(c *gin.Context) {
	projectID, err := strconv.Atoi(c.DefaultQuery("id", "0"))
	status := c.DefaultQuery("status", "")
	if err != nil {
		logging.Print.Error(err)
		c.JSON(http.StatusBadRequest, gin.Error{
			Err:  err,
			Type: 0,
			Meta: "error id format",
		})
		return
	}

	if projectID < 1 {
		c.JSON(http.StatusOK, middleware.GetBadRequest())
		return
	}

	var rowAffected int64
	project := middleware.GetProjectByID(projectID)
	switch status {
	case "New":
		rowAffected, err = project.UpdateStatus(1)
	case "Queue":
		rowAffected, err = project.UpdateStatus(2)
	case "Work":
		rowAffected, err = project.UpdateStatus(3)
	case "Test":
		rowAffected, err = project.UpdateStatus(4)
	case "Return":
		rowAffected, err = project.UpdateStatus(5)
	case "Done":
		rowAffected, err = project.UpdateStatus(6)
	case "Reopen":
		rowAffected, err = project.UpdateStatus(7)
	case "Close":
		rowAffected, err = project.UpdateStatus(8)
	default:
		c.JSON(http.StatusOK, middleware.GetBadRequest())
		return
	}
	if err != nil {
		c.JSON(http.StatusOK, gin.Error{
			Err:  err,
			Type: 0,
			Meta: "error update status",
		})
		return
	}
	if rowAffected == 0 {
		c.JSON(http.StatusOK, middleware.GetBadRequest())
		return
	}
	c.JSON(http.StatusOK, middleware.GetSuccess())
}
