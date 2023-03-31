package middleware

import (
	"gorm.io/gorm"
	"time"
)

type Budget struct {
	gorm.Model
	Name      string    `json:"name"`
	ExtID     int       `json:"ext_id"`
	ProjectID int       `json:"project_id"`
	StartAt   time.Time `json:"start_at"`
	EndAd     time.Time `json:"end_ad"`
	StatusID  int       `json:"status_id"`
	Issues    []Issue   `gorm:"foreignKey:BudgetID" json:"timespent,omitempty"`
}

func (b Budget) Insert() int64 {
	b.StatusID = 3 // в работе
	tx := DB.Create(&b)
	//if tx.Error != nil {
	//	logging.Print.Error("database error create budget")
	//	return 0, tx.Error
	//}
	return tx.RowsAffected
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

func GetBudgetById(id int) Budget {
	var b Budget
	DB.Where("id = ?", id).Find(&b)
	return b
}

func GetMainProjectBudgetByProjectID(projectID int) Budget {
	budget := Budget{}
	DB.Where("name = 'Основной' AND project_id = ?", projectID).Find(&budget)
	return budget
}

func (b Budget) UpdateStatus(statusId int) (int64, error) {
	tx := DB.Model(Budget{}).Where("id = ?", b.ID).Update("status_id", statusId)
	if tx.Error != nil {
		return 0, tx.Error
	}
	return tx.RowsAffected, nil
}
