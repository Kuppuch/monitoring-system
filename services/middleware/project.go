package middleware

import (
	"fmt"
	"gorm.io/gorm"
	"monitoring-system/services/logging"
)

type Project struct {
	gorm.Model
	Name        string  `json:"name"`
	Description string  `json:"description"`
	IsPublic    bool    `json:"isPublic"`
	StatusID    uint    `json:"status"`
	Issues      []Issue `gorm:"foreignKey:ProjectID" json:"issues"`
}

type ProjectWeb struct {
	Project
	IssuesCnt int64
}

func GetProjects() []Project {
	var projects []Project
	tx := DB.Find(&projects)
	if tx.Error != nil {
		logging.Print.Warning(tx.Error)
	}
	return projects
}

func (p ProjectWeb) GetIssueCount() {
	//var issues []Issue
	DB.Where("project_id = ?", p.ID).Find([]Issue{}).Count(&p.IssuesCnt)
	fmt.Println(p.IssuesCnt)
}

func GetProjectByID(id int) *Project {
	var project Project
	tx := DB.Where("id = ?", id).Find(&project)
	if tx.Error != nil {
		logging.Print.Warning(tx.Error, "id: ", id)
	}
	return &project
}

func (p Project) InsertProject() int64 {
	tx := DB.Create(&p)
	if tx.Error != nil {
		logging.Print.Warning(tx.Error)
	}
	return tx.RowsAffected
}
