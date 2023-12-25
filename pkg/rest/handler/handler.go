package handler

import (
	"os"

	"github.com/gin-gonic/gin"
)

type Config struct {
	R *gin.Engine
}

func NewHandler(c *Config) {
	urlPrefix := os.Getenv("URL_PREFIX")
	RegisterBlogRoutes(urlPrefix, c.R)
	RegisterHealthRoutes(urlPrefix, c.R)
}
