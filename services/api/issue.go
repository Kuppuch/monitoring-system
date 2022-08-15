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
	if err != nil {
		c.JSON(400, gin.Error{
			Err:  err,
			Type: 0,
			Meta: "failed parse uint",
		})
	}
	issues := middleware.GetIssueList(uint(projectID))
	c.JSON(http.StatusOK, issues)
}

func getIssueCreatePage(c *gin.Context) {
	c.HTML(http.StatusOK, "addIssue.html", nil)
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
	}
	issue.StatusID = 1
	issue.CreatorID = 1
	issue.AssignedToID = 1
	if rowAffected := issue.InsertIssue(); rowAffected == 0 {
		c.JSON(http.StatusBadRequest, middleware.GetBadRequest())
		return
	}
}
