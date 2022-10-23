package api

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"monitoring-system/services/logging"
	"monitoring-system/services/middleware"
	"net/http"
	"strconv"
)

func getProjectsPage(c *gin.Context) {
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

	user := GetUserByToken(c)
	if projectID > 0 {
		project := middleware.GetProjectByID(projectID)
		c.HTML(http.StatusOK, "project.html", gin.H{"project": project, "user": user})
		return
	}

	projects := middleware.GetProjects()
	c.HTML(http.StatusOK, "projects.html", gin.H{"projects": projects, "user": user})
}

func getProjectCreatePage(c *gin.Context) {
	user := GetUserByToken(c)
	c.HTML(http.StatusOK, "addProject.html", gin.H{"user": user})
}

func getMemberPage(c *gin.Context) {
	users := middleware.GetAllUsers()
	roles := middleware.GetRoles()
	user := GetUserByToken(c)
	c.HTML(http.StatusOK, "addMember.html", gin.H{"user": user, "users": users, "roles": roles})
}

func insertProject(c *gin.Context) {
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
	project := middleware.Project{}
	err = json.Unmarshal(raw, &project)
	if err != nil {
		logging.Print.Error("error unmarshal", err)
		c.JSON(http.StatusBadRequest, gin.Error{
			Err:  err,
			Type: 0,
			Meta: "error by unmarshal project",
		})
	}
	if len(project.Name) < 1 {
		c.JSON(http.StatusBadRequest, middleware.GetBadRequest())
		return
	}
	project.StatusID = 1
	if rowAffected := project.InsertProject(); rowAffected == 0 {
		c.JSON(http.StatusBadRequest, middleware.GetBadRequest())
		return
	}
	c.JSON(http.StatusOK, middleware.GetSuccess())
}
