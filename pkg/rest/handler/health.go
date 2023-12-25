package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterHealthRoutes(urlPrefix string, r *gin.Engine) {

	routes := r.Group(fmt.Sprintf("%s/healthcheck", urlPrefix))

	routes.GET("", HealthCheckHandler)

}

func HealthCheckHandler(c *gin.Context) {

	c.JSON(http.StatusOK, gin.H{
		"status": "UP",
	})
}
