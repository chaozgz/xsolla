package rest

import (
	"sxolla-rest-api/pkg/rest/handler"

	"github.com/gin-gonic/gin"
)

func inject() (*gin.Engine, error) {
	// initialize gin.Engine
	router := gin.Default()
	router.Use(CORSMiddleware())

	handler.NewHandler(&handler.Config{
		R: router,
	})

	return router, nil
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
