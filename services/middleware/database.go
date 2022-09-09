package middleware

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"monitoring-system/config"
)

var DB *gorm.DB

func Connector() error {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		config.DbHost, config.DbUser, config.DbPassword, config.DbName, config.DBPort, config.DbSslMode)
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}
	err = DB.AutoMigrate(&User{}, &Status{}, &Tracker{}, &Issue{}, &Project{})
	if err != nil {
		return err
	}
	InsertStatuses()
	return nil
}

func InsertStatuses() {
	var cas []Status
	DB.Find(&cas)
	if len(cas) == 0 {
		data := []string{
			"Новая",
			"В очереди исполнения",
			"В работе",
			"В тестировании",
			"Возвращено в работу",
			"Готова",
			"Переоткрыта",
			"Закрыта",
		}
		for _, v := range data {
			DB.Create(&Status{
				Name: v,
			})
		}
	}
}
