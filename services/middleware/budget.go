package middleware

import (
	"gorm.io/gorm"
	"monitoring-system/services/logging"
	"time"
)

type Budget struct {
	gorm.Model
	Name      string    `json:"name"`
	ExtID     int       `json:"ext_id"`
	ProjectID int       `json:"project_id"`
	StartAt   time.Time `json:"start_at"`
	EndAd     time.Time `json:"end_ad"`
	Issues    []Issue   `gorm:"foreignKey:BudgetID" json:"timespent,omitempty"`
}

func (b Budget) Insert() (int64, error) {
	tx := DB.Create(&b)
	if tx.Error != nil {
		logging.Print.Error("database error create budget")
		return 0, tx.Error
	}
	return tx.RowsAffected, nil
}

func GetBudget(id int) Budget {
	b := Budget{}
	DB.Where("id = ?", id).Find(&b)
	return b
}

func GetBudgets() []Budget {
	var b []Budget
	DB.Find(&b)
	return b
}

func GetProjectBudgets(projectID int) []Budget {
	var b []Budget
	DB.Where("project_id = ?", projectID).Find(&b)
	return b
}