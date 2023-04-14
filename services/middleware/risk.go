package middleware

import (
	"gorm.io/gorm"
	"math"
	"time"
)

type Risk struct {
	gorm.Model
	ProbabilityOccurrence float32 // Вероятность возникновения
	Influence             float32 // Влияние на проект
	Level                 float32 // Уровнь риска
	Description           string  // Описание риска
	Impact                string  // Воздействие на проект
	Solution              string  // Вариант решения
}

func GetRiskList() []Risk {
	var risks []Risk
	DB.Order("id desc").Find(&risks)
	return risks
}

func GetBannerInfo() (int, float64) {
	type riskInfo struct {
		Count     int
		UpdatedAt time.Time
	}
	ri := riskInfo{}
	DB.Raw("SELECT count(*), updated_at  FROM risks GROUP BY updated_at ORDER BY updated_at desc").Scan(&ri)
	timeBetween := math.Round(time.Now().Sub(ri.UpdatedAt).Minutes())
	return ri.Count, timeBetween
}

func (r *Risk) Insert() error {
	tx := DB.Create(r)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}
