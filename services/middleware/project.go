package middleware

import (
	"fmt"
	"gorm.io/gorm"
	"monitoring-system/services/logging"
)

type Project struct {
	gorm.Model
	Name        string `json:"name"`
	Description string `json:"description"`
	IsPublic    bool   `json:"isPublic"`
	StatusID    uint   `json:"status"`
	//Issues      []Issue  `gorm:"foreignKey:ProjectID" json:"issues"`
	Members []Member `gorm:"foreignKey:ProjectID" json:"members"`
	Budgets []Budget `gorm:"foreignKey:ProjectID" json:"budgets"`
}

type ProjectWeb struct {
	Project
	IssuesCnt int64
}

func GetProjects() []ProjectWeb {
	var projects []ProjectWeb
	tx := DB.Table("projects as p").Select("p.*, COUNT(i.budget_id) as issues_cnt").
		Joins("LEFT JOIN  budgets as b ON b.project_id = p.id").
		Joins("LEFT JOIN  issues as i ON i.budget_id = b.id").
		Group("p.id").
		Order("p.id").Find(&projects)
	if tx.Error != nil {
		logging.Print.Warning(tx.Error)
	}
	return projects
}

func (p ProjectWeb) GetIssueCount() {
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

func (p *Project) InsertProject() int64 {
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

func (p *Project) UpdateStatus(statusID int) (int64, error) {
	tx := DB.Model(&Project{}).Where("status_id = ?", p.ID).Update("status_id", statusID)
	if tx.Error != nil {
		return 0, tx.Error
	}
	return tx.RowsAffected, nil
}
