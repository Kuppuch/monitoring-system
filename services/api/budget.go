package api

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"monitoring-system/services/logging"
	"monitoring-system/services/middleware"
	"net/http"
	"strconv"
)

func getBudgets(c *gin.Context) {
	projectID, err := strconv.Atoi(c.DefaultQuery("project_id", "0"))
	if err != nil {
		c.JSON(400, gin.Error{
			Err:  err,
			Type: 0,
			Meta: "failed parse uint",
		})
		return
	}
	if projectID > 0 {
		budgets := middleware.GetProjectBudgets(projectID)
		c.JSON(http.StatusOK, budgets)
		return
	}
	user, _ := GetUserByToken(c)
	budgetsView := middleware.GetBudgetsWithProject()
	c.HTML(http.StatusOK, "budgets.html", gin.H{"budgets": budgetsView, "user": user})
}

func getBudget(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.Error{
			Err:  err,
			Type: 0,
			Meta: "failed parse int",
		})
		return
	}
	budget := middleware.GetBudget(id)
	user, _ := GetUserByToken(c)
	c.HTML(http.StatusOK, "budget.html", gin.H{"budget": budget, "user": user})
}

func insertBudget(c *gin.Context) {
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
	b := middleware.Budget{}
	err = json.Unmarshal(raw, &b)
	if err != nil {
		logging.Print.Error("error unmarshal", err)
		c.JSON(http.StatusBadRequest, gin.Error{
			Err:  err,
			Type: 0,
			Meta: "error by unmarshal issue",
		})
		return
	}
	if b.Name == "" {
		c.JSON(http.StatusBadRequest, gin.Error{
			Err:  err,
			Type: 0,
			Meta: "Название не может быть пустым",
		})
		return
	}
	if project := middleware.GetProjectByID(b.ProjectID); project.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.Error{
			Err:  err,
			Type: 0,
			Meta: "проект не существует",
		})
		return
	}

	b.StatusID = 1
	rawAffected := b.Insert()
	if err != nil || rawAffected == 0 {
		c.JSON(http.StatusInternalServerError, gin.Error{
			Err:  err,
			Type: 0,
			Meta: "database error create budget",
		})
		return
	}
	c.JSON(http.StatusOK, middleware.GetSuccess())
}

func getBudgetCreatePage(c *gin.Context) {
	projectID, err := strconv.Atoi(c.DefaultQuery("project_id", "0"))
	if err != nil {
		logging.Print.Error("error ", err)
		c.JSON(http.StatusBadRequest, middleware.GetBadRequest())
		return
	}
	project := middleware.GetProjectByID(projectID)
	if project.ID == 0 {
		c.JSON(http.StatusBadRequest, middleware.GetBadRequest())
		return
	}
	user, _ := GetUserByToken(c)
	c.HTML(http.StatusOK, "addBudget.html", gin.H{"user": user})
}

func getBudgetTimespent(c *gin.Context) {
	budgetId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.Error{
			Err:  err,
			Type: 0,
			Meta: "failed parse uint",
		})
		return
	}
	timespents := middleware.GetBudgetTimespent(budgetId)
	c.JSON(http.StatusOK, timespents)
}
