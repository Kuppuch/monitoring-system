package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	cors "github.com/rs/cors/wrapper/gin"
	"html/template"
	"monitoring-system/services/api/public"
	"monitoring-system/services/logging"
	"monitoring-system/services/middleware"
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
	r.Static("/favicon.ico", "./pages/img/favicon.ico")

	r.Use(cors.AllowAll())
	r.Use(AuthRequired)

	r.POST("upload", uploadProfileImg)

	r.GET("/", getMainPage)

	r.GET("/login", getLoginPage)
	r.POST("/login", login)

	r.GET("/submit", getSubmitPage)
	r.POST("/submit", submit)

	adm := r.Group("admin")
	{
		adm.GET("", getAdminPage)

		adm.GET("/roles", getRolesPage)
		adm.GET("/roles/:id", getRolePage)
		adm.GET("/roles/create", getRoleCreatePage)
		adm.POST("/roles/create", createRole)
		adm.PUT("/roles/:id", putRole)
		adm.DELETE("/roles/:id", deleteRole)

		adm.GET("/roles/head", getHeadRolesPage)
		adm.GET("/roles/head/:id", getHeadRolePage)
		adm.GET("/roles/head/create", getHeadRoleCreatePage)
		adm.POST("/roles/head/create", createHeadRole)
		adm.PUT("/roles/head/:id", putHeadRole)
		adm.DELETE("/roles/head/:id", deleteHeadRole)
	}

	user := r.Group("users")
	{
		uadm := user.Group("adm")
		{
			uadm.GET("/:id", getUserAdm)
			uadm.POST("/register", insertUser)
		}

		user.GET("/reg", getPage)
		user.GET("", getUser)
		user.PATCH("/upload", uploadProfileImg)
		user.PATCH("/update", updateProfile)
		user.GET("/:id", GetProfilePhoto)
	}

	project := r.Group("projects")
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
		project.GET("/timespent", getProjectTimespent)
		project.GET("/:id/info", getProjectInfo)
		project.POST("/:id/info/budget", getProjectBudgetInfo)

	}

	issue := r.Group("issues")
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
		budget.GET("/:id", getBudget)
		// внешнее api
		budget.POST("", insertBudget)
		// конец внешнего api
		budget.GET("/create", getBudgetCreatePage)
		budget.POST("/create", insertBudget)
		budget.GET("/:id/timespent", getBudgetTimespent)
	}

	notification := r.Group("notification")
	{
		notification.GET("/", getNotificationPage)
		notification.GET("/json", getNotifications)
		notification.GET("/socket", socketFunc)
		notification.GET("/send", sendMessage)
		notification.GET("/read", setReadNotification)
	}

	risk := r.Group("risks")
	{
		risk.GET("", getRiskPage)
		risk.GET("/list", getRiskListPage)
		//risk.GET("/list", getExpPage)
		risk.GET("/create", getRiskCreatePage)
		risk.POST("/create", InserRisk)
	}

	timespent := r.Group("timespents")
	{
		timespent.PATCH("/save", updateTimespent)
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
	if c.Request.URL.Path == "/login" || c.Request.URL.Path == "/favicon.ico" || c.Request.URL.Path == "/notification/socketFunc" {
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
