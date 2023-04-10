package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	cors "github.com/rs/cors/wrapper/gin"
	"html/template"
	"monitoring-system/services/api/public"
	"monitoring-system/services/logging"
	"monitoring-system/services/middleware"
	"net/http"
	"strings"
)

func Router() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.SetFuncMap(template.FuncMap{
		"upper": strings.ToUpper,
	})
	r.LoadHTMLGlob("pages/**/*.html")
	r.Static("/css", "./pages/css")
	r.Static("/js", "./pages/js")
	r.Static("/img", "./pages/img")
	r.Static("/photo", "./lib/users")

	r.Use(cors.AllowAll())
	//r.Use(AuthRequired)

	r.POST("upload", uploadProfileImg)

	r.GET("/", getMainPage)
	r.GET("/admin", getAdminPage)

	r.GET("/login", getLoginPage)
	r.POST("/login", login)

	r.GET("/submit", getSubmitPage)
	r.POST("/submit", submit)

	user := r.Group("user")
	{
		adm := user.Group("adm")
		{
			adm.GET("/:id", getUserAdm)
			adm.POST("/register", insertUser)
		}

		user.GET("/reg", getPage)
		user.GET("", getUser)
		user.PATCH("/upload", uploadProfileImg)
		user.GET("/:id", GetProfilePhoto)
	}

	project := r.Group("project")
	{
		project.GET("/", getProjectsPage)
		project.GET("/:id", getProjectPage)
		project.GET("/create", getProjectCreatePage)
		project.POST("/create", insertProject)
		project.GET("/:id/members", getMemberPage)
		project.POST("/:id/members", insertProjectMember)
		project.GET("/:id/members/list", getMembers)
		project.POST("/:id/link_rep", linkGitRepository)
		project.GET("/:id/link_rep", GetActualGitRepository)

	}

	issue := r.Group("issue")
	{
		issue.GET("", getIssueList)
		issue.GET("/:id", getIssueByID)
		issue.GET("/create", getIssueCreatePage)
		issue.POST("/create", insertIssue)
		issue.GET("/:id/timespent", getIssueUserTimespent)
		issue.POST("/:id/timespent", insertIssueUserTimespent)
		issue.GET("/my_issue", myIssuesPage)
		issue.POST("/save/:id", saveIssue)
	}

	budget := r.Group("budgets")
	{
		budget.GET("", getBudgets)
		budget.POST("", insertBudget)
		budget.GET("/:id/timespent", getBudgetTimespent)
	}

	notification := r.Group("notification")
	{
		notification.GET("/", func(c *gin.Context) {
			c.HTML(http.StatusOK, "notification.html", nil)
		})
		notification.GET("/socket", socket)
	}

	setting := r.Group("setting")
	{
		setting.GET("/", func(c *gin.Context) {
			c.HTML(http.StatusOK, "settings.html", nil)
		})
	}

	r.GET("/srs", srsPage)

	api := r.Group("api")
	public.Router(api)

	err := r.Run(":25595")
	if err != nil {
		logging.Print.Warning(err)
	}

}

func AuthRequired(c *gin.Context) {
	if c.Request.URL.Path == "/login" || c.Request.URL.Path == "/favicon.ico" {
		return
	}
	token, _ := c.Cookie("auth")
	if len(token) < 1 {
		c.Redirect(302, "/login")
	} else {
		_, err := middleware.CheckToken(token)
		if err != nil {
			fmt.Println(err)
			c.Redirect(302, "/login")
		}
	}
	c.Next()
}
