package middleware

import (
	"gorm.io/gorm"
	"monitoring-system/services/logging"
)

type GitRepository struct {
	gorm.Model
	ProjectID  int
	Repository string
}

func (gr GitRepository) Insert() (int64, error) {
	tx := DB.Create(&gr)
	if tx.Error != nil {
		logging.Print.Error("database error create git repository link")
		return 0, tx.Error
	}
	return tx.RowsAffected, nil
}

func (gr GitRepository) Update() (int64, error) {
	tx := DB.Model(&gr).Where("id=?", gr.ID).Update("repository", gr.Repository)
	if tx.Error != nil {
		logging.Print.Error("database error update git repository link")
		return 0, tx.Error
	}
	return tx.RowsAffected, nil
}

func GetGitRepositoryByID(id int) GitRepository {
	gr := GitRepository{}
	DB.Where("project_id = ?", id).Find(&gr)
	return gr
}
