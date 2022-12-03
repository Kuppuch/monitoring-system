package public

import (
	"github.com/gin-gonic/gin"
)

func Router(r *gin.RouterGroup) {

	projectapi := r.Group("projects")
	{
		projectapi.GET("/", GetProjects)
		//projectapi.GET("/create", getProjectCreatePage)
		//projectapi.POST("/create", insertProject)
		//projectapi.GET("/:id/members", getMemberPage)
		//projectapi.POST("/:id/members", insertProjectMember)

	}

	userapi := r.Group("users")
	{
		userapi.GET("/", GetUsers)
	}
}
