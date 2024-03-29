package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"monitoring-system/config"
	"monitoring-system/services/logging"
	"monitoring-system/services/middleware"
	"net/http"
	"sort"
	"strconv"
	"time"
)

func getProjectsPage(c *gin.Context) {
	user, _ := GetUserByToken(c)
	var projects []middleware.ProjectWeb
	if user.Admin {
		projects = middleware.GetAllProjects()
	} else {
		projects = middleware.GetProjects(user.ID)
	}
	for i, project := range projects {
		if len([]rune(project.Name)) > 15 {
			projects[i].Name = string([]rune(project.Name)[0:15]) + "..."
		}
	}

	c.HTML(http.StatusOK, "projects.html", gin.H{"projects": projects, "user": user})
}

func getProjectPage(c *gin.Context) {
	projectID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logging.Print.Error(err)
		c.JSON(http.StatusBadRequest, gin.Error{
			Err:  err,
			Type: 0,
			Meta: "error id format",
		})
		return
	}

	user, _ := GetUserByToken(c)
	if projectID > 0 {
		project := middleware.GetProjectByID(projectID)
		gitRepository := middleware.GetGitRepositoryByID(projectID)
		c.HTML(
			http.StatusOK,
			"project.html",
			gin.H{"project": project, "user": user, "token": config.GitToken, "git_repository": gitRepository.Repository})
		return
	}
	c.JSON(http.StatusBadRequest, gin.Error{
		Err:  err,
		Type: 0,
		Meta: "error id format",
	})
}

func getProjectCreatePage(c *gin.Context) {
	user, _ := GetUserByToken(c)
	c.HTML(http.StatusOK, "addProject.html", gin.H{"user": user})
}

func insertProject(c *gin.Context) {
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
	project := middleware.Project{}
	err = json.Unmarshal(raw, &project)
	m := map[string]interface{}{}
	err = json.Unmarshal(raw, &m)
	if err != nil {
		logging.Print.Error("error unmarshal", err)
		c.JSON(http.StatusBadRequest, gin.Error{
			Err:  err,
			Type: 0,
			Meta: "error by unmarshal project",
		})
		return
	}
	if len(project.Name) < 1 {
		c.JSON(http.StatusBadRequest, middleware.GetBadRequest())
		return
	}
	project.StatusID = 1
	dateStart, err := time.Parse("2006-01-02T15:04:05.000Z", m["dateStart"].(string))
	dateEnd, err := time.Parse("2006-01-02T15:04:05.000Z", m["dateEnd"].(string))
	if dateEnd.Sub(dateStart).Hours() < 8 {
		logging.Print.Error("trying to create a project less than 8 hours long")
		c.JSON(http.StatusBadRequest, gin.Error{
			Err:  errors.New("trying to create a project less than 8 hours long"),
			Type: 0,
			Meta: "Нельзя создавать проект длительностью менее 1 дня",
		})
		return
	}
	project.PlanStart = dateStart
	project.PlanFinish = dateEnd
	if rowAffected := project.InsertProject(); rowAffected == 0 {
		c.JSON(http.StatusBadRequest, middleware.GetBadRequest())
		return
	}
	budget := middleware.Budget{
		Name:      "Основной",
		ExtID:     0,
		ProjectID: int(project.ID),
		StartAt:   dateStart,
		EndAd:     dateEnd,
		StatusID:  3,
	}
	if rowAffected := budget.Insert(); rowAffected == 0 {
		c.JSON(http.StatusBadRequest, middleware.GetBadRequest())
		return
	}
	c.JSON(http.StatusOK, middleware.GetSuccess())
}

func getMemberPage(c *gin.Context) {
	users := middleware.GetAllUsers()
	roles := middleware.GetRoles()
	user, _ := GetUserByToken(c)
	c.HTML(http.StatusOK, "addMember.html", gin.H{"user": user, "users": users, "roles": roles})
}

func insertProjectMember(c *gin.Context) {
	type Member struct {
		User uint
		Role uint
	}

	projectIDstr := c.Param("id")
	projectID, err := strconv.Atoi(projectIDstr)
	if err != nil {
		logging.Print.Error(err)
		c.JSON(http.StatusBadRequest, gin.Error{
			Err:  err,
			Type: 0,
			Meta: "error by get project id",
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
	member := Member{}
	err = json.Unmarshal(raw, &member)
	if err != nil {
		logging.Print.Error("error unmarshal", err)
		c.JSON(http.StatusBadRequest, gin.Error{
			Err:  err,
			Type: 0,
			Meta: "error by unmarshal issue",
		})
	}
	memberDB := middleware.Member{
		ProjectID: uint(projectID),
		UserID:    member.User,
	}
	memberDB.GetMember()
	if memberDB.ID == 0 {
		rowAffected := memberDB.InsertMember()
		if rowAffected == 0 {
			logging.Print.Error("failed insert project member", err)
			c.JSON(http.StatusInternalServerError, gin.Error{
				Err:  err,
				Type: 0,
				Meta: "database error",
			})
			return
		}
	}

	projectRole := middleware.ProjectRole{
		MemberID: memberDB.ID,
		RoleID:   member.Role,
	}
	projectRole.GetProjectRole()
	if projectRole.ID > 0 {
		logging.Print.Error("failed insert project role: member-role already exist")
		c.JSON(http.StatusInternalServerError, gin.Error{
			Err:  err,
			Type: 0,
			Meta: "user already on project with this role",
		})
		return
	}
	rowAffected := projectRole.InsertProjectRole()
	if rowAffected != nil {
		logging.Print.Error("failed insert project role. rowAffected =", rowAffected)
		c.JSON(http.StatusInternalServerError, gin.Error{
			Err:  err,
			Type: 0,
			Meta: "database error",
		})
		return
	}
	c.JSON(http.StatusOK, middleware.GetSuccess())
}

func getMembers(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logging.Print.Error(err)
		c.JSON(http.StatusBadRequest, gin.Error{
			Err:  err,
			Type: 0,
			Meta: "error by get project id",
		})
		return
	}
	memberView := middleware.GetMembers(id)
	c.JSON(http.StatusOK, memberView)
}

func linkGitRepository(c *gin.Context) {
	projectIDstr := c.Param("id")
	projectID, err := strconv.Atoi(projectIDstr)
	if err != nil {
		logging.Print.Error(err)
		c.JSON(http.StatusBadRequest, gin.Error{
			Err:  err,
			Type: 0,
			Meta: "error by get project id",
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

	gr := middleware.GitRepository{}
	err = json.Unmarshal(raw, &gr)
	if err != nil {
		logging.Print.Error("error unmarshal", err)
		c.JSON(http.StatusBadRequest, gin.Error{
			Err:  err,
			Type: 0,
			Meta: "error by unmarshal issue",
		})
	}
	gr.ProjectID = projectID
	grFromDB := middleware.GetGitRepositoryByID(gr.ProjectID)

	if grFromDB.Repository == gr.Repository {
		c.JSON(http.StatusOK, middleware.GetSuccess())
		return
	}
	if grFromDB.Repository != gr.Repository && grFromDB.ProjectID == gr.ProjectID {
		grFromDB.Repository = gr.Repository
		rowAffected, err := grFromDB.Update()
		if err != nil || rowAffected == 0 {
			logging.Print.Error(fmt.Sprintf("error update git repository for project %v ", projectID), err)
			c.JSON(http.StatusBadRequest, gin.Error{
				Err:  err,
				Type: 0,
				Meta: "error update git repository",
			})
		}
		c.JSON(http.StatusOK, middleware.GetSuccess())
		return
	}

	rowAffected, err := gr.Insert()
	if err != nil || rowAffected == 0 {
		logging.Print.Error(fmt.Sprintf("error save git repository for project %v ", projectID), err)
		c.JSON(http.StatusBadRequest, gin.Error{
			Err:  err,
			Type: 0,
			Meta: "error save git repository",
		})
	}
	c.JSON(http.StatusOK, middleware.GetSuccess())
}

func GetActualGitRepository(c *gin.Context) {
	projectIDstr := c.Param("id")
	projectID, err := strconv.Atoi(projectIDstr)
	if err != nil {
		logging.Print.Error(err)
		c.JSON(http.StatusBadRequest, gin.Error{
			Err:  err,
			Type: 0,
			Meta: "error by get project id",
		})
		return
	}
	gr := middleware.GetGitRepositoryByID(projectID)
	m := map[string]interface{}{
		"Repository": gr.Repository,
	}
	b, err := json.Marshal(m)
	if err != nil {
		panic(err)
	}
	c.JSON(http.StatusOK, string(b))
}

func getProjectTimespent(c *gin.Context) {
	type RoleTimespent struct {
		RoleID    int `json:"role_id"`
		Sort      int
		Color     string  `json:"color"`
		Timespent float32 `json:"timespent"`
	}
	type ProjectInfo struct {
		ProjectStart   time.Time       `json:"project_start"`
		ProjectEnd     time.Time       `json:"project_end"`
		RoleTimespents []RoleTimespent `json:"role_timespents"`
	}

	projectTimespent := middleware.GetProjectTimespent()
	buildRoleTimespent := make(map[int][]RoleTimespent)
	for _, v := range projectTimespent {
		buildRoleTimespent[v.ProjectID] = append(buildRoleTimespent[v.ProjectID], RoleTimespent{
			RoleID:    v.RoleID,
			Sort:      v.Sort,
			Color:     v.Color,
			Timespent: v.Timespent,
		})
	}

	buildTimespent := make(map[int]ProjectInfo)
	for i, v := range buildRoleTimespent {
		sort.SliceStable(v, func(ii, jj int) bool {
			return v[ii].RoleID < v[jj].RoleID
		})
		start, end := middleware.GetProjectBoundByID(i)
		buildTimespent[i] = ProjectInfo{
			ProjectStart:   start,
			ProjectEnd:     end,
			RoleTimespents: v,
		}
	}
	for _, v := range buildTimespent {
		sort.SliceStable(v.RoleTimespents, func(i, j int) bool {
			return v.RoleTimespents[i].Sort < v.RoleTimespents[j].Sort
		})
	}
	c.JSON(http.StatusOK, buildTimespent)
}

func getProjectInfo(c *gin.Context) {
	projectID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logging.Print.Error(err)
		c.JSON(http.StatusBadRequest, gin.Error{
			Err:  err,
			Type: 0,
			Meta: "projectID не является числом",
		})
		return
	}
	mainBudget := middleware.GetMainProjectBudgetByProjectID(projectID)
	pti := middleware.GetProjectTimespentInfo(projectID, int(mainBudget.ID), time.Time{}, time.Time{})
	project := middleware.GetProjectByID(projectID)
	user, _ := GetUserByToken(c)
	budgets := middleware.GetProjectBudgets(projectID)

	c.HTML(http.StatusOK, "infoProject.html", gin.H{"user": user, "pti": pti, "project": project, "budgets": budgets})
}

func getProjectBudgetInfo(c *gin.Context) {
	projectID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logging.Print.Error(err)
		c.JSON(http.StatusBadRequest, gin.Error{
			Err:  err,
			Type: 0,
			Meta: "projectID не является числом",
		})
		return
	}
	raw, err := c.GetRawData()
	var m = make(map[string]interface{})
	err = json.Unmarshal(raw, &m)
	budgetID := int(m["budgetId"].(float64))
	if err != nil {
		logging.Print.Error(err)
		c.JSON(http.StatusBadRequest, gin.Error{
			Err:  err,
			Type: 0,
			Meta: "budgetID не является числом",
		})
		return
	}
	budget := middleware.Budget{}
	if budgetID == 0 {
		budget = middleware.GetMainProjectBudgetByProjectID(projectID)
	} else {
		budget = middleware.GetBudget(budgetID)
	}
	dateStartStr := m["dateStart"].(string)
	dateEndStr := m["dateEnd"].(string)
	dateStart, err := time.Parse("2006-01-02", dateStartStr)
	dateEnd, err := time.Parse("2006-01-02", dateEndStr)
	if err != nil && (len(dateStartStr) > 0 && len(dateEndStr) > 0) {
		logging.Print.Error(err)
		c.JSON(http.StatusBadRequest, gin.Error{
			Err:  err,
			Type: 0,
			Meta: "дата начала или конца не является датой",
		})
		return
	}
	pti := middleware.GetProjectTimespentInfo(projectID, int(budget.ID), dateStart, dateEnd)
	c.JSON(http.StatusOK, pti)
}
