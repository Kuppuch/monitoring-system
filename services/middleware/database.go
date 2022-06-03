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
	err = DB.AutoMigrate(&User{})
	if err != nil {
		return err
	}
	return nil
}
