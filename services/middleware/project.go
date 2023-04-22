package middleware

import (
	"fmt"
	"gorm.io/gorm"
	"monitoring-system/services/logging"
	"time"
)

type Project struct {
	gorm.Model
	Name            string          `json:"name"`
	Description     string          `json:"description"`
	IsPublic        bool            `json:"isPublic"`
	StatusID        uint            `json:"status"`
	PlanStart       time.Time       `json:"planStart"`
	PlanFinish      time.Time       `json:"planFinish"`
	Members         []Member        `gorm:"foreignKey:ProjectID" json:"members"`
	Budgets         []Budget        `gorm:"foreignKey:ProjectID" json:"budgets"`
	GitRepositories []GitRepository `gorm:"foreignKey:ProjectID" json:"gitRepositories"`
	//Issues          []Issue         `gorm:"foreignKey:ProjectID" json:"issues"`
}

type ProjectWeb struct {
	Project
	IssuesCnt int64
	Updated   int64
}

type ProjectTimespent struct {
	RoleID    int
	Sort      int
	Color     string
	Timespent float32
	ProjectID int
}

func GetProjects(userID uint) []ProjectWeb {
	var projects []ProjectWeb
	tx := DB.Table("projects as p").Select("p.*, COUNT(i.budget_id) as issues_cnt").
		Joins("LEFT JOIN  budgets as b ON b.project_id = p.id").
		Joins("LEFT JOIN  issues as i ON i.budget_id = b.id").
		Joins("INNER JOIN members AS m ON m.project_id = p.id").
		Where("m.user_id = ?", userID).
		Group("p.id").
		Order("p.id").Find(&projects)
	if tx.Error != nil {
		logging.Print.Warning(tx.Error)
	}
	for i, v := range projects {
		projects[i].Updated = int64(time.Now().Sub(v.UpdatedAt).Seconds())
	}
	return projects
}

func GetAllProjects() []ProjectWeb {
	var projects []ProjectWeb
	tx := DB.Table("projects as p").Select("p.*, COUNT(i.budget_id) as issues_cnt").
		Joins("LEFT JOIN  budgets as b ON b.project_id = p.id").
		Joins("LEFT JOIN  issues as i ON i.budget_id = b.id").
		Group("p.id").
		Order("p.id").Find(&projects)
	if tx.Error != nil {
		logging.Print.Warning(tx.Error)
	}
	for i, v := range projects {
		projects[i].Updated = int64(time.Now().Sub(v.UpdatedAt).Seconds())
	}
	return projects
}

func (p ProjectWeb) GetIssueCount() {
	DB.Where("project_id = ?", p.ID).Find([]Issue{}).Count(&p.IssuesCnt)
	fmt.Println(p.IssuesCnt)
}

func GetProjectByID(id int) *Project {
	var project Project
	tx := DB.Where("id = ?", id).Find(&project)
	if tx.Error != nil {
		logging.Print.Warning(tx.Error, "id: ", id)
	}
	return &project
}

func GetProjectBoundByID(id int) (time.Time, time.Time) {
	type TempStruct struct {
		Start time.Time
		End   time.Time
	}
	tempStruct := TempStruct{}
	DB.Raw(`SELECT p.id, MIN(b.start_at) AS start, MAX(b.end_ad) AS end
				  FROM budgets AS b
				 INNER JOIN projects AS p ON p.id = b.project_id
				 WHERE p.id = ?
				 GROUP BY p.id`, id).Scan(&tempStruct)
	return tempStruct.Start, tempStruct.End
}

func (p *Project) InsertProject() int64 {
	tx := DB.Create(&p)
	if tx.Error != nil {
		logging.Print.Warning(tx.Error)
	}
	return tx.RowsAffected
}

func (p Project) Delete() int64 {
	tx := DB.Delete(&p)
	if tx.Error != nil {
		logging.Print.Warning(tx.Error)
	}
	return tx.RowsAffected
}

func (p *Project) UpdateStatus(statusID int) (int64, error) {
	tx := DB.Model(&Project{}).Where("status_id = ?", p.ID).Update("status_id", statusID)
	if tx.Error != nil {
		return 0, tx.Error
	}
	return tx.RowsAffected, nil
}

func GetProjectTimespent() []ProjectTimespent {
	var projectTimespent []ProjectTimespent
	DB.Raw(`SELECT hr.id AS role_id, hr.sort, hr.color, SUM(t.spent) AS timespent, p.id AS project_id
				  FROM timespents AS t
				 INNER JOIN issues AS i ON i.id = t.issue_id
				 INNER JOIN budgets AS b ON b.id = i.budget_id
				 INNER JOIN projects AS p ON p.id = b.project_id
				 INNER JOIN roles AS r ON r.id = t.role_id
				 INNER JOIN head_roles AS hr ON hr.id = r.head_role_id
				 GROUP BY hr.id, p.id`).
		Find(&projectTimespent)
	return projectTimespent
}
