package api

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"monitoring-system/services/logging"
	"monitoring-system/services/middleware"
	"net/http"
	"os"
	"strconv"
)

func getPage(c *gin.Context) {
	c.HTML(http.StatusOK, "register.html", nil)
}

func getUser(c *gin.Context) {
	uid, err := strconv.Atoi(c.DefaultQuery("uid", "0"))
	if err != nil {
		c.JSON(400, gin.Error{
			Err:  err,
			Type: 0,
			Meta: "failed parse int",
		})
	}
	if uid == 0 {
		c.HTML(http.StatusOK, "users.html", gin.H{"users": middleware.GetAllUsers()})
		//c.JSON(http.StatusOK, middleware.GetAllUsers())
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
	c.HTML(http.StatusOK, "users.html", gin.H{"users": user})
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
		c.JSON(http.StatusInternalServerError, gin.Error{
			Err:  err,
			Type: 0,
			Meta: "perhaps raw is null",
		})
		return
	}
	//TODO переделать на длину
	if user.Email == "" || user.Password == "" || user.LastName == "" || user.Name == "" {
		c.JSON(http.StatusBadRequest, gin.Error{
			Err:  errors.New("null field"),
			Type: 0,
			Meta: "not enough data to register",
		})
		return
	}
	user.Password = fmt.Sprintf("%x", sha256.Sum256([]byte(user.Password)))
	if err != nil {
		logging.Print.Warning(err)
		c.JSON(http.StatusInternalServerError, gin.Error{
			Err:  err,
			Type: 0,
			Meta: "failed create password hash",
		})
		return
	}
	dbUser := middleware.User{}
	dbUser.Email = user.Email
	err = dbUser.GetUserByEmail()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.Error{
			Err:  err,
			Type: 0,
			Meta: "failed check existing user",
		})
		logging.Print.Error(err)
		return
	}
	if dbUser.ID > 0 {
		c.JSON(http.StatusBadRequest, gin.Error{
			Err:  errors.New("registration failed"),
			Type: 0,
			Meta: "user with given email address already exists",
		})
		return
	}

	rowAffected := user.InsertUser()
	if rowAffected > 0 {
		err = createUserRep(user.ID)
		if err != nil {
			logging.Print.Error("Unable to create user directory: ", err)
		}
		c.JSON(http.StatusOK, middleware.GetSuccess())
	} else {
		c.JSON(http.StatusBadRequest, middleware.GetBadRequest())
	}

}

func createUserRep(userID uint) error {
	err := os.Mkdir("lib/users/"+strconv.Itoa(int(userID)), 0777)
	if err != nil {
		return err
	}
	return nil
}
