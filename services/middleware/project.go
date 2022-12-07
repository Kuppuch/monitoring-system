package middleware

import (
	"fmt"
	"gorm.io/gorm"
	"monitoring-system/services/logging"
)

type Project struct {
	gorm.Model
	Name        string   `json:"name"`
	Description string   `json:"description"`
	IsPublic    bool     `json:"isPublic"`
	StatusID    uint     `json:"status"`
	Issues      []Issue  `gorm:"foreignKey:ProjectID" json:"issues"`
	Members     []Member `gorm:"foreignKey:ProjectID" json:"members"`
}

type ProjectWeb struct {
	Project
	IssuesCnt int64
}

func GetProjects() []ProjectWeb {
	var projects []ProjectWeb
	tx := DB.Table("projects as p").Select("p.*, COUNT(i.project_id) as issues_cnt").
		Joins("LEFT JOIN  issues as i ON i.project_id = p.id").
		Group("p.id").
		Order("p.id").Find(&projects)
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

func (p Project) Delete() int64 {
	tx := DB.Delete(&p)
	if tx.Error != nil {
		logging.Print.Warning(tx.Error)
	}
	return tx.RowsAffected
}
