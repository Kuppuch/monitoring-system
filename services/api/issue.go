package api

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"monitoring-system/services/api/socket"
	"monitoring-system/services/logging"
	"monitoring-system/services/middleware"
	"net/http"
	"regexp"
	"strconv"
	"strings"
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

	var issues []middleware.IssueWeb
	if budgetID > 0 {
		issues = middleware.GetIssueList(projectID, budgetID)
	} else {
		issues = middleware.GetIssueList(projectID, 0)
	}

	for i, v := range issues {
		timespent := middleware.GetIssueTimespent(int(v.ID))
		var counter float32 = 0
		for _, t := range timespent {
			counter += t.Spent
		}
		hour := int(counter)
		minute := (counter - float32(hour)) * 60
		issues[i].TimespentData = fmt.Sprintf("%vч. %vмин.", hour, minute)
	}
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
	user, _ := GetUserByToken(c)
	issue := middleware.GetIssue(uint(id))
	budget := middleware.GetBudget(int(issue.BudgetID))
	project := middleware.GetProjectByID(budget.ProjectID)

	issue.BudgetName = budget.Name
	issue.ProjectName = project.Name

	member := middleware.Member{
		ProjectID: project.ID,
		UserID:    user.ID,
	}
	member.GetMember()
	userProjectRolesID := middleware.GetProjectRoles(member.ID)
	var userProjectRoles []middleware.Role
	for _, v := range userProjectRolesID {
		userProjectRoles = append(userProjectRoles, middleware.GetRole(v.RoleID))
	}
	statuses := middleware.GetStatusList()

	timespent := middleware.GetIssueTimespent(id)
	var timespentCounter float32 = 0.0

	type Timespents struct {
		TimespentID uint
		UserID      uint
		Spent       float32
		User        string
		Role        string
	}

	timespents := []Timespents{}
	for _, v := range timespent {
		timespentCounter += v.Spent

		timespentAuthor := middleware.User{}
		timespentAuthor.ID = v.UserID
		_ = timespentAuthor.GetUser()
		role := middleware.GetRole(v.RoleID)
		timespents = append(timespents, Timespents{
			TimespentID: v.ID,
			UserID:      timespentAuthor.ID,
			User:        timespentAuthor.LastName + " " + timespentAuthor.Name,
			Spent:       v.Spent,
			Role:        role.Name,
		})
	}

	issue.ProjectID = int(project.ID)

	c.HTML(
		http.StatusOK,
		"issue.html",
		gin.H{"issue": issue,
			"statuses":     statuses,
			"user":         user,
			"projectRoles": userProjectRoles,
			"timespent":    timespentCounter,
			"timespents":   timespents,
		})
}

func getIssueCreatePage(c *gin.Context) {
	projectID, err := strconv.Atoi(c.DefaultQuery("project_id", "0"))
	if err != nil {
		logging.Print.Error("error ", err)
		c.JSON(http.StatusBadRequest, middleware.GetBadRequest())
		return
	}
	var budgets []middleware.Budget
	if projectID > 0 {
		project := middleware.GetProjectByID(projectID)
		if project.ID == 0 {
			c.JSON(http.StatusBadRequest, middleware.GetBadRequest())
			return
		}
		budgets = middleware.GetProjectBudgets(int(project.ID))
	}

	user, _ := GetUserByToken(c)
	statuses := middleware.GetStatusList()
	trackers := middleware.GetTrackerList()
	assigned := middleware.GetAllUsers()

	c.HTML(http.StatusOK, "addIssue.html", gin.H{
		"user":     user,
		"statuses": statuses,
		"trackers": trackers,
		"assigned": assigned,
		"budgets":  budgets,
	})
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
	m := map[string]interface{}{}
	err = json.Unmarshal(raw, &m)
	if err != nil {
		logging.Print.Error("error unmarshal", err)
		c.JSON(http.StatusBadRequest, gin.Error{
			Err:  err,
			Type: 0,
			Meta: "error by unmarshal issue",
		})
		return
	}

	user, err := GetUserByToken(c)
	issue.CreatorID = user.ID
	project := middleware.GetProjectByID(issue.ProjectID)
	if project.ID == 0 {
		c.JSON(http.StatusBadRequest, middleware.GetBadRequest())
		return
	}

	if rowAffected := issue.InsertIssue(); rowAffected == 0 {
		c.JSON(http.StatusBadRequest, middleware.GetBadRequest())
		return
	}
	if issue.CreatorID != issue.AssignedToID {
		notification := middleware.Notification{
			CreatorID:    issue.CreatorID,
			AssignedToID: issue.AssignedToID,
			Content:      "На вас назначена новая задача",
			View:         false,
			Source:       "issue/" + strconv.Itoa(int(issue.ID)),
		}
		notification.Insert()
		// Отправляем уведомление в канал уведомлений
		nByte, _ := json.Marshal(&notification)
		socket.BigChannel <- nByte
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
	user, err := GetUserByToken(c)
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

// saveIssue - сохраняет изменения в задаче, в том числе трудосписания (ворклоги, timespent)
func saveIssue(c *gin.Context) {
	user, err := GetUserByToken(c)
	if user.ID < 1 {
		c.JSON(http.StatusBadRequest, gin.Error{
			Err:  err,
			Type: 0,
			Meta: "incorrect user",
		})
		return
	}
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

	m := map[string]string{}
	err = json.Unmarshal(raw, &m)
	if err != nil {
		logging.Print.Error("error unmarshal", err)
		c.JSON(http.StatusBadRequest, gin.Error{
			Err:  err,
			Type: 0,
			Meta: "error by unmarshal issue",
		})
		return
	}
	t, statusID, err := parseTimespent(m, user, issueID)
	if err != nil {
		logging.Print.Error("error parse timespent ", err)
		c.JSON(http.StatusBadRequest, gin.Error{
			Err:  err,
			Type: 0,
			Meta: err.Error(),
		})
		return
	}
	middleware.StatusUpdate(issueID, statusID)

	if t.Spent != 0 {
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
	}

	timespent := middleware.GetIssueTimespent(issueID)
	timespentMap := map[string]float32{
		"timespent": 0.0,
	}
	for _, v := range timespent {
		timespentMap["timespent"] += v.Spent
	}
	c.JSON(http.StatusOK, timespentMap)
}

func parseTimespent(m map[string]string, user middleware.User, issueID int) (middleware.Timespent, int, error) {
	spentStr, ok := m["spent_str"]
	if spentStr == "" {
		spentStr = "0h 0m"
	}
	statusStr, ok := m["status"]
	role_idStr, ok := m["role_id"]
	if !ok {
		return middleware.Timespent{}, 0, errors.New("error parse timespent")
	}
	roleID, err := strconv.Atoi(role_idStr)
	statusID, err := strconv.Atoi(statusStr)
	if err != nil {
		return middleware.Timespent{}, 0, err
	}
	hourStr := ""
	minuteStr := ""
	hourDelimiter := 0
	for i := 0; i < len(spentStr); i++ {
		if spentStr[i] == 'h' {
			hourStr = spentStr[:i]
			hourDelimiter = i + 1
		}
		if spentStr[i] == 'm' {
			minuteStr = strings.TrimSpace(spentStr[hourDelimiter:i])
		}
	}

	if len(spentStr) > 0 && len(hourStr) < 1 && len(minuteStr) < 1 {
		if matched, _ := regexp.MatchString(`\d\s\d\d`, spentStr); matched {
			hourStr = spentStr[0:1]
			minuteStr = spentStr[2:4]
		}
		if matched, _ := regexp.MatchString(`^\d{1,2}$`, spentStr); matched {
			if len(spentStr) < 2 {
				hourStr = spentStr
			} else if len(spentStr) < 3 {
				minuteStr = spentStr
			} else {
				return middleware.Timespent{}, 0, errors.New("минуты должны быть не больше 60")
			}
		}
	}

	if len(hourStr) < 1 {
		hourStr = "0"
	}
	if len(minuteStr) < 1 {
		minuteStr = "0"
	}
	hour, err := strconv.Atoi(hourStr)
	minute, err := strconv.Atoi(minuteStr)
	if err != nil {
		return middleware.Timespent{}, 0, err
	}
	if minute%15 != 0 {
		return middleware.Timespent{}, 0, errors.New("минуты должны быть кратны 15")
	}
	if minute > 59 {
		hour = minute / 60
		minute = minute % 60
	}
	minuteFloat := float32(minute) / float32(60)
	t := middleware.Timespent{
		IssueID:  uint(issueID),
		UserID:   user.ID,
		RoleID:   uint(roleID),
		Spent:    float32(hour) + minuteFloat,
		SpentStr: spentStr,
	}
	return t, statusID, nil
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

	user, err := GetUserByToken(c)
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

func myIssuesPage(c *gin.Context) {
	user, _ := GetUserByToken(c)
	issues := middleware.GetUserIssues(user.ID)
	c.HTML(http.StatusOK, "myIssues.html", gin.H{"user": user, "issues": issues})
}
