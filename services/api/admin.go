package api

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"monitoring-system/services/logging"
	"monitoring-system/services/middleware"
	"net/http"
)

func getAdminPage(c *gin.Context) {
	c.HTML(http.StatusOK, "admin.html", nil)
}

func getRolePage(c *gin.Context) {
	roles := middleware.GetFullInfoRoles()
	c.HTML(http.StatusOK, "role.html", gin.H{"roles": roles})
}

func getHeadRolePage(c *gin.Context) {
	headRoles := middleware.GetHeadRoles()
	c.HTML(http.StatusOK, "headRole.html", gin.H{"headRoles": headRoles})
}

func getHeadRoleCreatePage(c *gin.Context) {
	c.HTML(http.StatusOK, "addHeadRole.html", nil)
}

func createHeadRole(c *gin.Context) {
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
	headRole := middleware.HeadRole{}
	err = json.Unmarshal(raw, &headRole)
	if err != nil {
		logging.Print.Error("error unmarshal", err)
		c.JSON(http.StatusBadRequest, gin.Error{
			Err:  err,
			Type: 0,
			Meta: "error by unmarshal headRole",
		})
	}
	if err = headRole.Insert(); err != nil {
		logging.Print.Error("error create head role ", err)
		c.JSON(http.StatusOK, middleware.GetBadRequest())
		return
	}
	c.JSON(http.StatusOK, middleware.GetSuccess())
}

func getRoleCreatePage(c *gin.Context) {
	c.HTML(http.StatusOK, "addRole.html", gin.H{"headRoles": middleware.GetHeadRoles()})
}

func createRole(c *gin.Context) {
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

	role := middleware.Role{}
	err = json.Unmarshal(raw, &role)
	if err != nil {
		logging.Print.Error("error unmarshal", err)
		c.JSON(http.StatusBadRequest, gin.Error{
			Err:  err,
			Type: 0,
			Meta: "error by unmarshal headRole",
		})
		return
	}
	if err = role.Insert(); err != nil {
		logging.Print.Error("error create role ", err)
		c.JSON(http.StatusOK, middleware.GetBadRequest())
		return
	}
	c.JSON(http.StatusOK, middleware.GetSuccess())
}
