package api

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"monitoring-system/services/logging"
	"monitoring-system/services/middleware"
	"net/http"
)

func getRiskPage(c *gin.Context) {
	user, err := GetUserByToken(c)
	if err != nil {
		logging.Print.Error("error getting user from token ", err)
		c.JSON(http.StatusBadRequest, middleware.GetBadRequest())
		return
	}
	count, timeBetween := middleware.GetBannerInfo()
	c.HTML(http.StatusOK, "risk.html", gin.H{"user": user, "countRisk": count, "timeBetween": timeBetween})
}

func getRiskListPage(c *gin.Context) {
	risks := middleware.GetRiskList()
	user, err := GetUserByToken(c)
	if err != nil {
		logging.Print.Error(err)
		c.JSON(http.StatusBadRequest, middleware.GetBadRequest())
		return
	}
	c.HTML(http.StatusOK, "riskList.html", gin.H{"user": user, "risks": risks})
}

func getRiskCreatePage(c *gin.Context) {
	user, err := GetUserByToken(c)
	if err != nil {
		logging.Print.Error(err)
		c.JSON(http.StatusBadRequest, middleware.GetBadRequest())
		return
	}
	c.HTML(http.StatusOK, "riskCreate.html", gin.H{"user": user})
}

func InserRisk(c *gin.Context) {
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
	risk := middleware.Risk{}
	err = json.Unmarshal(raw, &risk)
	if err != nil {
		logging.Print.Error("error unmarshal", err)
		c.JSON(http.StatusBadRequest, gin.Error{
			Err:  err,
			Type: 0,
			Meta: "error by unmarshal risk",
		})
		return
	}
	if err = risk.Insert(); err != nil {
		logging.Print.Error(err)
		c.JSON(http.StatusBadRequest, gin.Error{
			Err:  err,
			Type: 0,
			Meta: "error insert risk into DB",
		})
		return
	}
	c.JSON(http.StatusOK, middleware.GetSuccess())
}
