package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type HttpStatus struct {
	Code   int
	Status string
}

func GetSuccess() map[string]interface{} {
	return gin.H{"Code": http.StatusOK, "Status": "success"}
}

func GetBadRequest() map[string]interface{} {
	return gin.H{"Code": http.StatusBadRequest, "Status": "error"}
}

func GetNotFound() map[string]interface{} {
	return gin.H{"Code": http.StatusNotFound, "Status": "error"}
}
