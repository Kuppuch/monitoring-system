package api

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"monitoring-system/services/middleware"
	"net/http"
	"strconv"
)

func getSubmitPage(c *gin.Context) {
	c.HTML(http.StatusOK, "submit.html", nil)
}

func submit(c *gin.Context) {
	user := GetUserByToken(c)
	raw, err := c.GetRawData()
	if err != nil {
		fmt.Println(err)
		return
	}
	type Submit struct {
		Action   string `json:"action"`
		Entity   string `json:"entity"`
		Id       string `json:"id"`
		Password string `json:"password"`
	}
	s := Submit{}
	err = json.Unmarshal(raw, &s)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.Error{
			Err:  err,
			Type: 0,
			Meta: "incorrect data",
		})
		return
	}
	if user.Password != fmt.Sprintf("%x", sha256.Sum256([]byte(s.Password))) {
		fmt.Println("неверный пароль")
		c.JSON(http.StatusBadRequest, gin.Error{
			Err:  err,
			Type: 0,
			Meta: "incorrect password",
		})
		return
	}
	id, err := strconv.Atoi(s.Id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.Error{
			Err:  err,
			Type: 0,
			Meta: "incorrect id",
		})
	}
	if DeleteProject(id) {
		c.JSON(http.StatusOK, nil)
		return
	}
	c.JSON(http.StatusBadRequest, gin.Error{
		Err:  err,
		Type: 0,
		Meta: "incorrect id",
	})
}

func DeleteProject(id int) bool {
	project := middleware.Project{
		Model: gorm.Model{ID: uint(id)},
	}
	rowAffected := project.Delete()
	if rowAffected > 0 {
		return true
	}
	return false
}
