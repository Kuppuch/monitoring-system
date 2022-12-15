package api

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"monitoring-system/services/logging"
	"monitoring-system/services/middleware"
	"net/http"
	"strconv"
)

func getIssueList(c *gin.Context) {
	projectID, err := strconv.Atoi(c.Query("project_id"))
	budgetID, err := strconv.Atoi(c.DefaultQuery("budget_id", "0"))
	if err != nil {
		c.JSON(400, gin.Error{
			Err:  err,
			Type: 0,
			Meta: "failed parse uint",
		})
		return
	}
	if budgetID > 0 {
		issues := middleware.GetBudgetIssue(budgetID)
		c.JSON(http.StatusOK, issues)
		return
	}
	issues := middleware.GetIssueList(uint(projectID))
	c.JSON(http.StatusOK, issues)
}

func getIssueByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.Error{
			Err:  err,
			Type: 0,
			Meta: "failed parse uint",
		})
		return
	}
	issue := middleware.GetIssue(uint(id))
	statuses := middleware.GetStatusList()
	user := GetUserByToken(c)
	c.HTML(http.StatusOK, "issue.html", gin.H{"issue": issue, "statuses": statuses, "user": user})
}

func getIssueCreatePage(c *gin.Context) {
	user := GetUserByToken(c)
	statuses := middleware.GetStatusList()
	trackers := middleware.GetTrackerList()
	assigned := middleware.GetAllUsers()
	c.HTML(http.StatusOK, "addIssue.html", gin.H{"user": user, "statuses": statuses, "trackers": trackers, "assigned": assigned})
}

func insertIssue(c *gin.Context) {
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
	issue := middleware.Issue{}
	err = json.Unmarshal(raw, &issue)
	if err != nil {
		logging.Print.Error("error unmarshal", err)
		c.JSON(http.StatusBadRequest, gin.Error{
			Err:  err,
			Type: 0,
			Meta: "error by unmarshal issue",
		})
		return
	}

	issue.CreatorID = GetUserByToken(c).ID

	if rowAffected := issue.InsertIssue(); rowAffected == 0 {
		c.JSON(http.StatusBadRequest, middleware.GetBadRequest())
		return
	}
}

func getIssueUserTimespent(c *gin.Context) {
	issueID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.Error{
			Err:  err,
			Type: 0,
			Meta: "issue id is not number",
		})
		return
	}
	user := GetUserByToken(c)
	if user.ID < 1 {
		c.JSON(http.StatusBadRequest, gin.Error{
			Err:  err,
			Type: 0,
			Meta: "incorrect user",
		})
		return
	}
	t := middleware.GetUserIssueTimespent(issueID, int(user.ID))
	c.JSON(http.StatusOK, t)
}

func insertIssueUserTimespent(c *gin.Context) {
	issueID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.Error{
			Err:  err,
			Type: 0,
			Meta: "issue id is not number",
		})
		return
	}
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
	t := middleware.Timespent{}
	err = json.Unmarshal(raw, &t)
	if err != nil {
		logging.Print.Error("error unmarshal", err)
		c.JSON(http.StatusBadRequest, gin.Error{
			Err:  err,
			Type: 0,
			Meta: "error by unmarshal issue",
		})
		return
	}

	user := GetUserByToken(c)
	if user.ID < 1 {
		c.JSON(http.StatusBadRequest, gin.Error{
			Err:  err,
			Type: 0,
			Meta: "incorrect user",
		})
		return
	}
	t.UserID = user.ID
	t.IssueID = uint(issueID)

	rowAffected, err := t.Insert()
	if err != nil {
		logging.Print.Error("error insert timespent", err)
		c.JSON(http.StatusInternalServerError, gin.Error{
			Err:  err,
			Type: 0,
			Meta: "error insert timespent",
		})
		return
	}
	if rowAffected == 0 {
		logging.Print.Error("error insert timespent", err)
		c.JSON(http.StatusInternalServerError, gin.Error{
			Err:  err,
			Type: 0,
			Meta: "error insert timespent",
		})
		return
	}
	c.JSON(http.StatusOK, t)
}
