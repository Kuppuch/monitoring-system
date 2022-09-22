package api

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"monitoring-system/services/logging"
	"monitoring-system/services/middleware"
	"net/http"
	"time"
)

type Claims struct {
	jwt.StandardClaims
	ID uint
}

type AuthParam struct {
	Login    string
	Password string
}

func getLoginPage() {

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

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 10).Unix(), //10 часов
			IssuedAt:  time.Now().Unix(),
		},
		ID: user.ID,
	})
	token, _ := t.SignedString([]byte("123"))
	c.Header("auth", token)
}
