package api

import (
	"github.com/gin-gonic/gin"
	cors "github.com/rs/cors/wrapper/gin"
	"html/template"
	"monitoring-system/services/logging"
	"strings"
)

func Router() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.SetFuncMap(template.FuncMap{
		"upper": strings.ToUpper,
	})
	r.LoadHTMLGlob("pages/**/*")
	r.Static("/css", "pages/css")
	r.Use(cors.AllowAll())

	r.GET("/", getMainPage)
	r.GET("/admin", getAdminPage)

	user := r.Group("user")
	{
		user.GET("/reg", getPage)
		user.GET("", getUser)
		user.POST("/register", insertUser)
	}

	project := r.Group("project")
	{
		project.GET("/", getProjectsPage)
		project.GET("/create", getProjectCreatePage)
		project.POST("/create", insertProject)
	}

	issue := r.Group("issue")
	{
		issue.GET("", getIssueList)
		issue.POST("/create", insertIssue)
	}

	err := r.Run(":25595")
	if err != nil {
		logging.Print.Warning(err)
	}

}
