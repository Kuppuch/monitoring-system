package middleware

import (
	"fmt"
	"gorm.io/gorm"
	"monitoring-system/services/logging"
)

type Issue struct {
	gorm.Model
	Name         string
	Description  string
	ProjectID    int         `gorm:"-:all" json:"project_id"`
	CreatorID    uint        `json:"creator_id"`
	AssignedToID uint        `json:"assigned_to_id"`
	StatusID     uint        `json:"status_id"`
	TrackerID    uint        `json:"tracker_id"`
	BudgetID     uint        `json:"budget_id"`
	Timespent    []Timespent `gorm:"foreignKey:UserID" json:"timespent,omitempty"`
}

type IssueWeb struct {
	Issue
	StatusName  string
	TrackerName string
	Creator     string
	AssignedTo  string
	BudgetName  string
}

func GetIssueList(projectID, budgetID int) []IssueWeb {
	var issues []IssueWeb
	where := fmt.Sprintf("b.project_id = %v", projectID)
	if budgetID > 0 {
		where = fmt.Sprintf("i.budget_id = %v", budgetID)
	}
	DB.Table("issues AS i").
		Select("i.*, b.name AS budget_name, t.name AS tracker_name, s.name AS status_name, " +
			"uc.name || ' ' || uc.last_name AS creator, ua.name || ' ' || ua.last_name AS assigned_to").
		Joins("INNER JOIN budgets AS b ON b.id = i.budget_id").
		Joins("INNER JOIN trackers AS t ON t.id = i.tracker_id").
		Joins("INNER JOIN statuses AS s ON s.id = i.status_id").
		Joins("INNER JOIN users AS uc ON uc.id = i.creator_id").
		Joins("INNER JOIN users AS ua ON ua.id = i.assigned_to_id").
		Where(where).Find(&issues)
	return issues
}

func GetIssue(id uint) IssueWeb {
	issue := IssueWeb{}
	DB.Table("issues").
		Select("issues.id, issues.name, issues.description, statuses.name as status_name, trackers.name as tracker_name, "+
			"issues.creator_id, issues.assigned_to_id, issues.status_id, issues.tracker_id, issues.budget_id,"+
			"u.last_name || ' ' || u.name || ' ' || u.middle_name as creator, uu.last_name || ' ' || uu.name || ' ' || uu.middle_name as assigned_to").
		Joins("inner join statuses on statuses.id = issues.status_id").
		Joins("inner join trackers on trackers.id = issues.tracker_id").
		Joins("inner join users as u on u.id = issues.creator_id").
		Joins("inner join users as uu on uu.id = issues.assigned_to_id").
		Where("issues.id = ?", id).Find(&issue)
	return issue
}

func GetUserIssues(id uint) []IssueWeb {
	var issues []IssueWeb
	DB.Table("issues").
		Select("issues.id, issues.name, issues.description, statuses.name as status_name, trackers.name as tracker_name, "+
			"issues.creator_id, issues.assigned_to_id, issues.status_id, issues.tracker_id, "+
			"u.last_name || ' ' || u.name || ' ' || u.middle_name as creator, uu.last_name || ' ' || uu.name || ' ' || uu.middle_name as assigned_to").
		Joins("inner join statuses on statuses.id = issues.status_id").
		Joins("inner join trackers on trackers.id = issues.tracker_id").
		Joins("inner join users as u on u.id = issues.creator_id").
		Joins("inner join users as uu on uu.id = issues.assigned_to_id").
		Where("issues.assigned_to_id = ?", id).Find(&issues)
	return issues
}

func (i Issue) InsertIssue() int64 {
	tx := DB.Create(&i)
	if tx.Error != nil {
		logging.Print.Warning(tx.Error)
	}
	return tx.RowsAffected
}

func StatusUpdate(issueID int, statusID int) {
	tx := DB.Table("issues").Where("id = ?", issueID).Update("status_id", statusID)
	if tx.Error != nil {
		logging.Print.Error(tx.Error)
	}
}

func GetBudgetIssue(id int) []Issue {
	var i []Issue
	DB.Where("budget_id = ?", id).Find(&i)
	return i
}

func GetBudgetIDByIssue(issueID int) int {
	var i Issue
	DB.Where("id = ?", issueID).Find(&i)
	return int(i.BudgetID)
}
