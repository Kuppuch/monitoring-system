package api

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"monitoring-system/services/logging"
	"monitoring-system/services/middleware"
	"net/http"
)

type Claims struct {
	jwt.StandardClaims
	ID uint
}

type AuthParam struct {
	Login    string
	Password string
}

func getLoginPage(c *gin.Context) {
	c.HTML(http.StatusOK, "auth.html", nil)
}

func login(c *gin.Context) {
	raw, err := c.GetRawData()
	if err != nil {
		logging.Print.Warning(err)
	}
	auth := AuthParam{}
	err = json.Unmarshal(raw, &auth)
	if err != nil {
		logging.Print.Warning(err)
		c.JSON(http.StatusInternalServerError, gin.Error{
			Err:  err,
			Type: 0,
			Meta: err.Error(),
		})
		return
	}

	user := middleware.User{Email: auth.Login}
	_ = user.GetUserByEmail()
	hash := fmt.Sprintf("%x", sha256.Sum256([]byte(auth.Password)))
	if err != nil {
		logging.Print.Warning(err)
		c.JSON(http.StatusInternalServerError, gin.Error{
			Err:  err,
			Type: 0,
			Meta: "failed create password hash",
		})
		return
	}
	if user.ID == 0 || hash != user.Password {
		c.JSON(http.StatusBadRequest, gin.Error{
			Err:  errors.New("user not found"),
			Type: 0,
			Meta: "user not found",
		})
		return
	}

	token := middleware.GetToken(user.ID)

	authInstance := middleware.Auth{
		UserID: user.ID,
		Token:  token,
	}
	err = authInstance.InsertAuth()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.Error{
			Err:  errors.New("auth err"),
			Type: 0,
			Meta: "auth err",
		})
		return
	}

	c.Header("Authorization", token)
	c.JSON(http.StatusOK, gin.H{"Authorization": token})
}

func GetUserByToken(c *gin.Context) middleware.User {
	token, _ := c.Cookie("auth")
	uid, _ := middleware.CheckToken(token)
	user := middleware.User{}
	user.ID = uid
	_ = user.GetUser()
	return user
}
