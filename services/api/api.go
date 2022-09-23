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
	//html := template.Must(template.ParseFiles("./pages/img/logo.png"))
	//r.SetHTMLTemplate(html)
	r.SetFuncMap(template.FuncMap{
		"upper": strings.ToUpper,
	})
	r.LoadHTMLGlob("pages/**/*.html")
	r.Static("/css", "./pages/css")
	r.Static("/js", "./pages/js")
	r.Use(cors.AllowAll())

	r.GET("/", getMainPage)
	r.GET("/admin", getAdminPage)

	r.GET("/login", getLoginPage)
	r.POST("/login", login)

	user := r.Group("user")
	{
		user.GET("/reg", getPage)
		user.GET("", getUser)
		user.POST("/register", insertUser)
		user.PATCH("/upload_img", uploadProfileImg)
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
		issue.GET("/:id", getIssueByID)
		issue.GET("/create", getIssueCreatePage)
		issue.POST("/create", insertIssue)
	}

	err := r.Run(":25595")
	if err != nil {
		logging.Print.Warning(err)
	}

}
