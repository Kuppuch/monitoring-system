package api

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"monitoring-system/services/logging"
	"monitoring-system/services/middleware"
	"net/http"
)

func updateTimespent(c *gin.Context) {
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
	m := map[string]float64{}
	err = json.Unmarshal(raw, &m)
	if err != nil || m["value"] == 0 {
		logging.Print.Error("error unmarshal", err)
		c.JSON(http.StatusBadRequest, gin.Error{
			Err:  err,
			Type: 0,
			Meta: "ошибка в типе данных. Попробуйте ввести число",
		})
		return
	}
	timespent := middleware.GetTimespentByID(int(m["timespentID"]))
	if timespent.ID < 1 {
		logging.Print.Error("bad timespent id", err)
		c.JSON(http.StatusBadRequest, gin.Error{
			Err:  err,
			Type: 0,
			Meta: "трудозатрат с таким id не существует",
		})
		return
	}

	newSpent := float32(m["value"])
	timespent.Spent = newSpent

	if err = timespent.Update(); err != nil {
		logging.Print.Error("failed update timespent", err)
		c.JSON(http.StatusInternalServerError, gin.Error{
			Err:  err,
			Type: 0,
			Meta: "Не удалось обновить трудозатраты",
		})
		return
	}
	c.JSON(http.StatusOK, timespent)
}
