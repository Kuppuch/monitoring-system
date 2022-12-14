package api

import (
	"github.com/gin-gonic/gin"
	"strconv"
)

func getBudgetTimespent(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.Error{
			Err:  err,
			Type: 0,
			Meta: "failed parse uint",
		})
		return
	}
	_ = id

}
