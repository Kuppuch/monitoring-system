package middleware

import (
	"gorm.io/gorm"
	"monitoring-system/services/logging"
)

type Issue struct {
	gorm.Model
	Name         string
	Description  string
	ProjectID    uint        `json:"project_id"`
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
}

func GetIssueList(projectID uint) []Issue {
	var issues []Issue
	DB.Where("project_id = ?", projectID).Find(&issues)
	return issues
}

func GetIssue(id uint) IssueWeb {
	issue := IssueWeb{}
	DB.Table("issues").
		Select("issues.id, issues.name, issues.description, statuses.name as status_name, trackers.name as tracker_name, projects.id, "+
			"issues.project_id, issues.creator_id, issues.assigned_to_id, issues.status_id, issues.tracker_id, "+
			"u.last_name || ' ' || u.name || ' ' || u.middle_name as creator, uu.last_name || ' ' || uu.name || ' ' || uu.middle_name as assigned_to").
		Joins("inner join statuses on statuses.id = issues.status_id").
		Joins("inner join trackers on trackers.id = issues.tracker_id").
		Joins("inner join projects on projects.id = issues.project_id").
		Joins("inner join users as u on u.id = issues.creator_id").
		Joins("inner join users as uu on uu.id = issues.assigned_to_id").
		Where("issues.id = ?", id).Find(&issue)
	return issue
}

func (i Issue) InsertIssue() int64 {
	tx := DB.Create(&i)
	if tx.Error != nil {
		logging.Print.Warning(tx.Error)
	}
	return tx.RowsAffected
}

func StatusUpdate(statusID uint) {
	tx := DB.Where("id = ?").Update("StatusID", statusID)
	if tx.Error != nil {

	}
}

func GetBudgetIssue(id int) []Issue {
	var i []Issue
	DB.Where("budget_id = ?", id).Find(&i)
	return i
}
