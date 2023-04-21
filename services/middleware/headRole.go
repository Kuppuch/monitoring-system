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

func GetHeadRole(headRoleID int) HeadRole {
	var headRole HeadRole
	DB.Where("id = ?", headRoleID).Find(&headRole)
	return headRole
}

func (hr *HeadRole) Update() error {
	tx := DB.Save(hr)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

func DeleteHeadRole(headRoleID int) error {
	tx := DB.Where("id = ?", headRoleID).Delete(&HeadRole{})
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}
