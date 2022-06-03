package api

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"monitoring-system/services/logging"
	"monitoring-system/services/middleware"
	"net/http"
	"strconv"
)

func getUser(c *gin.Context) {
	uid, err := strconv.Atoi(c.Query("uid"))
	if err != nil {
		c.JSON(400, middleware.HttpStatus{
			Code:   400,
			Status: "error",
		})
	}
	if uid == 0 {
		c.JSON(http.StatusOK, middleware.GetAll())
		return
	}
	user := middleware.User{}
	user.ID = uint(uid)
	err = user.GetUser()
	if err != nil {
		logging.Print.Warning(err)
	}
	user.Password = ""
	c.JSON(http.StatusOK, user)
}

func insertUser(c *gin.Context) {
	raw, err := c.GetRawData()
	if err != nil {
		logging.Print.Warning(err)
	}
	user := middleware.User{}
	err = json.Unmarshal(raw, &user)
	if err != nil {
		logging.Print.Warning(err)
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		logging.Print.Warning(err)
	}
	user.Password = string(hash)
	rowAffected := user.InsertUser()
	if rowAffected > 0 {
		c.JSON(http.StatusOK, middleware.HttpStatus{
			Code:   http.StatusOK,
			Status: "success",
		})
	} else {
		c.JSON(http.StatusBadRequest, middleware.HttpStatus{
			Code:   http.StatusBadRequest,
			Status: "error",
		})
	}

}
