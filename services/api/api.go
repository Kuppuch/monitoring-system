package api

import (
	"github.com/gin-gonic/gin"
	cors "github.com/rs/cors/wrapper/gin"
	"monitoring-system/services/logging"
)

func Router() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.LoadHTMLGlob("pages/*.html")
	r.Use(cors.AllowAll())
	r.GET("", getPage)
	user := r.Group("user")
	{
		user.GET("", getUser)
		user.POST("/register", insertUser)
	}

	err := r.Run(":25595")
	if err != nil {
		logging.Print.Warning(err)
	}

}

func DummyMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		AddCORSHeaders(c)
	}
}

//func HandleCORSOptionsRequest(c *gin.Context) {
//	c.Header("Access-Control-Allow-Origin", "*")
//	c.Header("Access-Control-Allow-Methods", "DELETE, GET, OPTIONS, PATCH, POST, PUT")
//	c.Header("Access-Control-Allow-Headers", "accept, accept-encoding, authorization, content-type, dnt, origin, user-agent, x-csrftoken, x-requested-with, X-Token-Bearer")
//	c.Header("Access-Control-Max-Age", "86400")
//}

func AddCORSHeaders(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", "DELETE, GET, OPTIONS, PATCH, POST, PUT")
	c.Header("Access-Control-Allow-Headers", "accept, accept-encoding, authorization, content-type, dnt, origin, user-agent, x-csrftoken, x-requested-with, x-token-bearer")
	c.Header("Access-Control-Max-Age", "86400")
}
