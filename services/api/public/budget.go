package public

import (
	"github.com/gin-gonic/gin"
	"monitoring-system/services/logging"
	"monitoring-system/services/middleware"
	"net/http"
	"strconv"
)

func GetBudget(c *gin.Context) {
	budgetID, err := strconv.Atoi(c.DefaultQuery("id", "0"))
	if err != nil {
		logging.Print.Error(err)
		c.JSON(http.StatusBadRequest, gin.Error{
			Err:  err,
			Type: 0,
			Meta: "error id format",
		})
		return
	}

	if budgetID > 0 {
		budget := middleware.GetBudgetById(budgetID)
		c.JSON(http.StatusOK, budget)
		return
	}

	budgets := middleware.GetBudgets()
	c.JSON(http.StatusOK, budgets)
}

func UpdateStatusBudget(c *gin.Context) {
	budgetID, err := strconv.Atoi(c.DefaultQuery("id", "0"))
	if err != nil {
		logging.Print.Error(err)
		c.JSON(http.StatusBadRequest, gin.Error{
			Err:  err,
			Type: 0,
			Meta: "error id format",
		})
		return
	}
	if budgetID < 1 {
		c.JSON(http.StatusOK, middleware.GetBadRequest())
		return
	}
	budget := middleware.GetBudgetById(budgetID)

	status := c.DefaultQuery("status", "0")
	if budget.StatusID == 8 {
		c.JSON(http.StatusOK, middleware.GetBadRequest())
		return
	}
	var rowAffected int64
	switch status {
	case "New":
		rowAffected, err = budget.UpdateStatus(1)
	case "Queue":
		rowAffected, err = budget.UpdateStatus(2)
	case "Work":
		rowAffected, err = budget.UpdateStatus(3)
	case "Test":
		rowAffected, err = budget.UpdateStatus(4)
	case "Return":
		rowAffected, err = budget.UpdateStatus(5)
	case "Done":
		rowAffected, err = budget.UpdateStatus(6)
	case "Reopen":
		rowAffected, err = budget.UpdateStatus(7)
	case "Close":
		rowAffected, err = budget.UpdateStatus(8)
	default:
		c.JSON(http.StatusOK, middleware.GetBadRequest())
		return
	}
	if err != nil {
		c.JSON(http.StatusOK, gin.Error{
			Err:  err,
			Type: 0,
			Meta: "error update status",
		})
		return
	}
	if rowAffected == 0 {
		c.JSON(http.StatusOK, middleware.GetBadRequest())
		return
	}
	c.JSON(http.StatusOK, middleware.GetSuccess())
}
