package middleware

import (
	"gorm.io/gorm"
	"monitoring-system/services/logging"
)

type Issue struct {
	gorm.Model
	Name         string
	Description  string
	StatusName   string `gorm:"-:all"`
	ProjectID    uint   `json:"project_id"`
	CreatorID    uint   `json:"creator_id"`
	AssignedToID uint   `json:"assigned_to_id"`
	StatusID     uint   `json:"status"`
}

func GetIssueList(projectID uint) []Issue {
	var issues []Issue
	DB.Where("project_id = ?", projectID).Find(&issues)
	return issues
}

func GetIssue(id uint) Issue {
	issue := Issue{}
	DB.Where("id = ?", id).Find(&issue)
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
