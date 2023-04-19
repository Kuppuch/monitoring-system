package middleware

import "gorm.io/gorm"

type HeadRole struct {
	gorm.Model
	Name string
}
