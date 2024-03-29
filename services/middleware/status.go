package middleware

import "gorm.io/gorm"

type Status struct {
	gorm.Model
	Name     string    `json:"name"`
	Issues   []Issue   `gorm:"foreignKey:StatusID" json:"issues"`
	Projects []Project `gorm:"foreignKey:StatusID" json:"projects"`
	Budgets  []Budget  `gorm:"foreignKey:StatusID" json:"budgets"`
}

func GetStatusList() []Status {
	var statuses []Status
	DB.Find(&statuses)
	return statuses
}
