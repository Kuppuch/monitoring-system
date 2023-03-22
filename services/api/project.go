package api

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"monitoring-system/config"
	"monitoring-system/services/logging"
	"monitoring-system/services/middleware"
	"net/http"
	"strconv"
)

func getProjectsPage(c *gin.Context) {
	user, _ := GetUserByToken(c)
	projects := middleware.GetProjects()
	c.HTML(http.StatusOK, "projects.html", gin.H{"projects": projects, "user": user})
}

func getProjectPage(c *gin.Context) {
	projectID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logging.Print.Error(err)
		c.JSON(http.StatusBadRequest, gin.Error{
			Err:  err,
			Type: 0,
			Meta: "error id format",
		})
		return
	}

	user, _ := GetUserByToken(c)
	if projectID > 0 {
		project := middleware.GetProjectByID(projectID)
		c.HTML(http.StatusOK, "project.html", gin.H{"project": project, "user": user, "token": config.GitToken})
		return
	}
	c.JSON(http.StatusBadRequest, gin.Error{
		Err:  err,
		Type: 0,
		Meta: "error id format",
	})
}

func getProjectCreatePage(c *gin.Context) {
	user, _ := GetUserByToken(c)
	c.HTML(http.StatusOK, "addProject.html", gin.H{"user": user})
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

func getMemberPage(c *gin.Context) {
	users := middleware.GetAllUsers()
	roles := middleware.GetRoles()
	user, _ := GetUserByToken(c)
	c.HTML(http.StatusOK, "addMember.html", gin.H{"user": user, "users": users, "roles": roles})
}

func insertProjectMember(c *gin.Context) {
	type Member struct {
		User uint
		Role uint
	}

	projectIDstr := c.Param("id")
	projectID, err := strconv.Atoi(projectIDstr)
	if err != nil {
		logging.Print.Error(err)
		c.JSON(http.StatusBadRequest, gin.Error{
			Err:  err,
			Type: 0,
			Meta: "error by get project id",
		})
		return
	}
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
	member := Member{}
	err = json.Unmarshal(raw, &member)
	if err != nil {
		logging.Print.Error("error unmarshal", err)
		c.JSON(http.StatusBadRequest, gin.Error{
			Err:  err,
			Type: 0,
			Meta: "error by unmarshal issue",
		})
	}
	memberDB := middleware.Member{
		ProjectID: uint(projectID),
		UserID:    member.User,
	}
	memberDB.GetMember()
	if memberDB.ID == 0 {
		rowAffected := memberDB.InsertMember()
		if rowAffected == 0 {
			logging.Print.Error("failed insert project member", err)
			c.JSON(http.StatusInternalServerError, gin.Error{
				Err:  err,
				Type: 0,
				Meta: "database error",
			})
			return
		}
	}

	projectRole := middleware.ProjectRole{
		MemberID: memberDB.ID,
		RoleID:   member.Role,
	}
	projectRole.GetProjectRole()
	if projectRole.ID > 0 {
		logging.Print.Error("failed insert project role: member-role already exist")
		c.JSON(http.StatusInternalServerError, gin.Error{
			Err:  err,
			Type: 0,
			Meta: "user already on project with this role",
		})
		return
	}
	rowAffected := projectRole.InsertRole()
	if rowAffected == 0 {
		logging.Print.Error("failed insert project role. rowAffected =", rowAffected)
		c.JSON(http.StatusInternalServerError, gin.Error{
			Err:  err,
			Type: 0,
			Meta: "database error",
		})
		return
	}
	c.JSON(http.StatusOK, middleware.GetSuccess())
}

func getMembers(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logging.Print.Error(err)
		c.JSON(http.StatusBadRequest, gin.Error{
			Err:  err,
			Type: 0,
			Meta: "error by get project id",
		})
		return
	}
	memberView := middleware.GetMembers(id)
	c.JSON(http.StatusOK, memberView)
}
