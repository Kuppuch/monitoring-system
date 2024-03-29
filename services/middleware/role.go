package middleware

import (
	"gorm.io/gorm"
)

type Role struct {
	gorm.Model
	Name         string
	HeadRoleID   int
	ProjectRoles []ProjectRole `gorm:"foreignKey:RoleID" json:"project_roles,omitempty"`
	Timespents   []Timespent   `gorm:"foreignKey:RoleID" json:"timespents,omitempty"`
}

type RoleInfo struct {
	Role
	HeadRoleName string
}

type ProjectRole struct {
	gorm.Model
	MemberID uint
	RoleID   uint
}

func GetRoles() []Role {
	var roles []Role
	DB.Find(&roles)
	return roles
}

func GetFullInfoRoles() []RoleInfo {
	var rolesInfo []RoleInfo
	DB.Table("roles").Select("roles.id, roles.name, roles.head_role_id, hr.name as head_role_name").
		Joins("LEFT JOIN head_roles AS hr ON hr.id = roles.head_role_id").
		Order("roles.id").
		Find(&rolesInfo)
	return rolesInfo
}

func GetRole(roleID uint) Role {
	role := Role{}
	DB.Where("id = ?", roleID).Find(&role)
	return role
}

func DeleteRole(roleID int) error {
	tx := DB.Where("id = ?", roleID).Delete(&Role{})
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

func (r *ProjectRole) GetProjectRole() {
	DB.Where("member_id = ? AND role_id = ?", r.MemberID, r.RoleID).Find(r)
}

func GetProjectRoles(memberID uint) []ProjectRole {
	var projectRoles []ProjectRole
	DB.Where("member_id = ?", memberID).Find(&projectRoles)
	return projectRoles
}

func (r *Role) Insert() error {
	tx := DB.Create(r)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

func (r *Role) Update() error {
	tx := DB.Save(r)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

func (r *ProjectRole) InsertProjectRole() error {
	tx := DB.Create(r)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

func UnlickHeadRole(headRole int) error {
	tx := DB.Model(&Role{}).Select("head_role_id").
		Where("head_role_id = ?", headRole).Updates(map[string]interface{}{"head_role_id": nil})
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}
