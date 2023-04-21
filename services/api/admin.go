package api

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"monitoring-system/services/logging"
	"monitoring-system/services/middleware"
	"net/http"
	"strconv"
)

func getAdminPage(c *gin.Context) {
	c.HTML(http.StatusOK, "admin.html", nil)
}

func getRolesPage(c *gin.Context) {
	roles := middleware.GetFullInfoRoles()
	c.HTML(http.StatusOK, "role.html", gin.H{"roles": roles})
}

func getRolePage(c *gin.Context) {
	roleID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logging.Print.Error("error convert id parament in getRolePage ", err)
		c.JSON(http.StatusNotFound, middleware.GetNotFound())
		return
	}
	role := middleware.GetRole(uint(roleID))
	headRoles := middleware.GetHeadRoles()
	c.HTML(http.StatusOK, "addRole.html", gin.H{"role": role, "headRoles": headRoles})
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

func putRole(c *gin.Context) {
	roleID, err := strconv.Atoi(c.Param("id"))
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
			Meta: "error by unmarshal role",
		})
	}
	role.ID = uint(roleID)
	if err = role.Update(); err != nil {
		logging.Print.Error("error update head role ", err)
		c.JSON(http.StatusBadRequest, middleware.GetBadRequest())
		return
	}
	c.JSON(http.StatusOK, middleware.GetSuccess())
}

func deleteRole(c *gin.Context) {
	roleID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logging.Print.Error("error convert id parament in getRolePage ", err)
		c.JSON(http.StatusNotFound, middleware.GetNotFound())
		return
	}
	if err = middleware.DeleteRole(roleID); err != nil {
		logging.Print.Error("error delete role ", err)
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, middleware.GetSuccess())
}

func getHeadRolesPage(c *gin.Context) {
	headRoles := middleware.GetHeadRoles()
	c.HTML(http.StatusOK, "headRole.html", gin.H{"headRoles": headRoles})
}

func getHeadRolePage(c *gin.Context) {
	roleID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logging.Print.Error("error convert id parament in getRolePage ", err)
		c.JSON(http.StatusNotFound, middleware.GetNotFound())
		return
	}
	headRole := middleware.GetHeadRole(roleID)
	c.HTML(http.StatusOK, "addHeadRole.html", gin.H{"headRole": headRole})
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
		c.JSON(http.StatusBadRequest, middleware.GetBadRequest())
		return
	}
	c.JSON(http.StatusOK, middleware.GetSuccess())
}

func putHeadRole(c *gin.Context) {
	roleID, err := strconv.Atoi(c.Param("id"))
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
	headRole.ID = uint(roleID)
	if err = headRole.Update(); err != nil {
		logging.Print.Error("error update head role ", err)
		c.JSON(http.StatusBadRequest, middleware.GetBadRequest())
		return
	}
	c.JSON(http.StatusOK, middleware.GetSuccess())
}

func deleteHeadRole(c *gin.Context) {
	roleID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logging.Print.Error("error convert id parament in getRolePage ", err)
		c.JSON(http.StatusNotFound, middleware.GetNotFound())
		return
	}
	err = middleware.DeleteHeadRole(roleID)
	if err != nil {
		logging.Print.Error("error delete head role ", err)
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	err = middleware.UnlickHeadRole(roleID)
	if err != nil {
		logging.Print.Error("error unlink head role from roles ", err)
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, middleware.GetSuccess())
}
