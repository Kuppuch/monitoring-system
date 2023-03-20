package public

import (
	"github.com/gin-gonic/gin"
)

func Router(r *gin.RouterGroup) {

	projectapi := r.Group("projects")
	{
		projectapi.GET("/", GetProjects)
		//projectapi.GET("/create", getProjectCreatePage)
		projectapi.POST("/", InsertProject)
		//projectapi.GET("/:id/members", getMemberPage)
		//projectapi.POST("/:id/members", insertProjectMember)

	}

	budgets := r.Group("budgets")
	{
		budgets.GET("/", GetBudget)
		budgets.POST("/", UpdateStatusBudget)
	}

	userapi := r.Group("users")
	{
		userapi.GET("/", GetUsers)
	}
}
