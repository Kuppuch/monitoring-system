package middleware

type ProjectTimespentInfo struct {
	Spent          float32
	EstimatedHours int
	RoleName       string
	IssueName      string
	BudgetID       int
	BudgetName     string
}

func GetProjectTimespentInfo(projectID int, budgetID int) []ProjectTimespentInfo {
	var projectTimespentsInfo []ProjectTimespentInfo
	DB.Raw(`SELECT SUM(t.spent) AS spent, i.estimated_hours, r.name AS role_name, i.name AS issue_name, b.id AS budget_id, b.name AS budget_name
				FROM timespents AS t
						 INNER JOIN issues AS i ON i.id = t.issue_id
						 INNER JOIN budgets AS b ON b.id = i.budget_id
						 INNER JOIN projects AS p ON p.id = b.project_id
						 INNER JOIN roles AS r ON r.id = t.role_id
				 WHERE p.id = ? AND b.id = ?
				GROUP BY r.name, i.id, i.name, b.id, b.name, p.id, p.name
				ORDER BY i.id`, projectID, budgetID).Find(&projectTimespentsInfo)
	return projectTimespentsInfo
}
