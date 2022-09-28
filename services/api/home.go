package api

import (
	"encoding/base64"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"monitoring-system/services/middleware"
	"net/http"
	"strconv"
)

func getMainPage(c *gin.Context) {
	bts, err := ioutil.ReadFile("./pages/img/logo.png")
	if err != nil {
		log.Fatal(err)
	}

	var base64Encoding string
	mimeType := http.DetectContentType(bts)
	switch mimeType {
	case "image/jpeg":
		base64Encoding += "data:image/jpeg;base64,"
	case "image/png":
		base64Encoding += "data:image/png;base64,"
	}
	base64Encoding += base64.StdEncoding.EncodeToString(bts)

	token, _ := c.Cookie("auth")
	userID, err := middleware.CheckToken(token)
	if err != nil || userID == 0 {
		c.JSON(http.StatusInternalServerError, gin.Error{
			Err:  err,
			Type: 0,
			Meta: "incorrect token for user with id = " + strconv.Itoa(int(userID)),
		})
		return
	}
	user := middleware.User{}
	user.ID = userID
	err = user.GetUser()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.Error{
			Err:  err,
			Type: 0,
			Meta: err.Error(),
		})
		return
	}

	c.HTML(http.StatusOK, "index.html", gin.H{
		"user": user,
	})
}
