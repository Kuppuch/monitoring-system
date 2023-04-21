package middleware

import "gorm.io/gorm"

type HeadRole struct {
	gorm.Model
	Name string
	Sort int
}

func (hr *HeadRole) Insert() error {
	tx := DB.Create(hr)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

func GetHeadRoles() []HeadRole {
	var headRoles []HeadRole
	DB.Find(&headRoles)
	return headRoles
}
