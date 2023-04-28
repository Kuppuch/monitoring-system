package middleware

import (
	"fmt"
	"time"
)

type ProjectTimespentInfo struct {
	Spent          float32
	EstimatedHours int
	IssueId        int
	RoleName       string
	IssueName      string
	BudgetID       int
	BudgetName     string
}

func GetProjectTimespentInfo(projectID int, budgetID int, dateStart time.Time, dateEnd time.Time) []ProjectTimespentInfo {
	var projectTimespentsInfo []ProjectTimespentInfo
	whereStr := ""
	if dateStart.Year() > 2000 && dateEnd.Year() > 2000 {
		whereStr = fmt.Sprintf("AND i.updated_at BETWEEN '%v-%v-%v 00:00:00' AND '%v-%v-%v 23:59:59'",
			dateStart.Year(),
			int(dateStart.Month()),
			dateStart.Day(),
			dateEnd.Year(),
			int(dateEnd.Month()),
			dateEnd.Day())
	} else if dateStart.Year() > 2000 {
		whereStr = fmt.Sprintf("AND i.updated_at > '%v-%v-%v 00:00:00'",
			dateStart.Year(),
			int(dateStart.Month()),
			dateStart.Day())
	} else if dateEnd.Year() > 2000 {
		whereStr = fmt.Sprintf("AND i.updated_at < '%v-%v-%v 00:00:00'",
			dateEnd.Year(),
			int(dateEnd.Month()),
			dateEnd.Day())
	}

	DB.Table("timespents AS t").Select("SUM(t.spent) AS spent, i.estimated_hours, i.id as issue_id, i.name AS issue_name, b.id AS budget_id, b.name AS budget_name").
		Joins("INNER JOIN issues AS i ON i.id = t.issue_id").
		Joins("INNER JOIN budgets AS b ON b.id = i.budget_id").
		Joins("INNER JOIN projects AS p ON p.id = b.project_id").
		Joins("INNER JOIN roles AS r ON r.id = t.role_id").
		Where("p.id = ? AND b.id = ? "+whereStr, projectID, budgetID).
		Group("i.id, i.name, b.id, b.name, p.id, p.name").
		Order("i.id").Find(&projectTimespentsInfo)
	//DB.Raw(`SELECT SUM(t.spent) AS spent, i.estimated_hours, r.name AS role_name, i.name AS issue_name, b.id AS budget_id, b.name AS budget_name
	//			FROM timespents AS t
	//					 INNER JOIN issues AS i ON i.id = t.issue_id
	//					 INNER JOIN budgets AS b ON b.id = i.budget_id
	//					 INNER JOIN projects AS p ON p.id = b.project_id
	//					 INNER JOIN roles AS r ON r.id = t.role_id
	//			 WHERE p.id = ? AND b.id = ?
	//			GROUP BY r.name, i.id, i.name, b.id, b.name, p.id, p.name
	//			ORDER BY i.id`, projectID, budgetID).Find(&projectTimespentsInfo)
	return projectTimespentsInfo
}
