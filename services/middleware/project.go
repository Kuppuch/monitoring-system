package middleware

import (
	"gorm.io/gorm"
	"monitoring-system/services/logging"
)

type Project struct {
	gorm.Model
	Name        string
	Description string
	IsPublic    bool `json:"isPublic"`
	Status      uint
}

func GetProjects() []Project {
	var projects []Project
	tx := DB.Find(&projects)
	if tx.Error != nil {
		logging.Print.Warning(tx.Error)
	}
	return projects
}

func (p Project) InsertProject() int64 {
	tx := DB.Create(&p)
	if tx.Error != nil {
		logging.Print.Warning(tx.Error)
	}
	return tx.RowsAffected
}
