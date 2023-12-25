package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"sxolla-rest-api/ent"
	con "sxolla-rest-api/internal/model"
	"sxolla-rest-api/pkg/model"
	svc "sxolla-rest-api/service"

	"github.com/gin-gonic/gin"
)

// type blogHandler struct {
// 	// Auth middle
// }

func RegisterBlogRoutes(urlPrefix string, r *gin.Engine) {
	// handler := &blogHandler{}

	routes := r.Group(fmt.Sprintf("%s/posts", urlPrefix))

	routes.GET("/:id", BlogGetByID)
	routes.GET("", BlogGetAll)
	routes.POST("", BlogCreate)
	routes.PUT("/:id", BlogUpdate)
	routes.DELETE("/:id", BlogDelete)

}

func BlogGetByID(c *gin.Context) {
	blogIdStr := c.Param(con.ID)
	blogId, err := strconv.Atoi(blogIdStr)
	if err != nil {
		rspErr := model.NewBadRequest(err.Error())
		c.JSON(rspErr.Status(), gin.H{
			con.ERROR: rspErr.Message,
		})
		return
	}

	blog, err := svc.NewBlogOps(c.Request.Context()).BlogGetByID(blogId)
	if err != nil {
		rspErr := model.NewNotFound(con.BLOG, blogIdStr)
		c.JSON(rspErr.Status(), gin.H{
			con.ERROR: rspErr.Message,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		con.BLOG: blog,
	})
}

func BlogGetAll(c *gin.Context) {
	blogs, err := svc.NewBlogOps(c.Request.Context()).BlogGetAll()
	if err != nil {
		rspErr := model.NewBadRequest(err.Error())
		c.JSON(rspErr.Status(), gin.H{
			con.ERROR: rspErr.Message,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		con.BLOGS: blogs,
	})
}

func BlogCreate(c *gin.Context) {
	if c.Request.Body == nil {
		rspErr := model.NewBadRequest(string(model.ERROR_MSG_EMPTY_REQ_BODY))
		c.JSON(rspErr.Status(), gin.H{
			con.ERROR: rspErr.Message,
		})
		return
	}
	jsonData, err := io.ReadAll(c.Request.Body)
	if err != nil {
		rspErr := model.NewBadRequest(err.Error())
		c.JSON(rspErr.Status(), gin.H{
			con.ERROR: rspErr.Message,
		})
		return
	}
	var blog ent.Blog
	if err = json.Unmarshal(jsonData, &blog); err != nil {
		rspErr := model.NewBadRequest(err.Error())
		c.JSON(rspErr.Status(), gin.H{
			con.ERROR: rspErr.Message,
		})
		return
	}
	newBlog, err := svc.NewBlogOps(c.Request.Context()).BlogCreate(&blog)
	if err != nil {
		rspErr := model.NewBadRequest(err.Error())
		c.JSON(rspErr.Status(), gin.H{
			con.ERROR: rspErr.Message,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		con.BLOG: newBlog,
	})
}

func BlogUpdate(c *gin.Context) {
	if c.Request.Body == nil {
		rspErr := model.NewBadRequest(string(model.ERROR_MSG_EMPTY_REQ_BODY))
		c.JSON(rspErr.Status(), gin.H{
			con.ERROR: rspErr.Message,
		})
		return
	}
	jsonData, err := io.ReadAll(c.Request.Body)
	if err != nil {
		rspErr := model.NewBadRequest(err.Error())
		c.JSON(rspErr.Status(), gin.H{
			con.ERROR: rspErr.Message,
		})
		return
	}
	var blog ent.Blog
	if err = json.Unmarshal(jsonData, &blog); err != nil {
		rspErr := model.NewBadRequest(err.Error())
		c.JSON(rspErr.Status(), gin.H{
			con.ERROR: rspErr.Message,
		})
		return
	}

	blogIdStr := c.Param(con.ID)
	blogId, err := strconv.Atoi(blogIdStr)
	if err != nil {
		rspErr := model.NewBadRequest(err.Error())
		c.JSON(rspErr.Status(), gin.H{
			con.ERROR: rspErr.Message,
		})
		return
	}
	blog.ID = blogId

	updatedBlog, err := svc.NewBlogOps(c.Request.Context()).BlogUpdate(&blog)
	if err != nil {
		rspErr := model.NewNotFound(con.BLOG, blogIdStr)
		c.JSON(rspErr.Status(), gin.H{
			con.ERROR: rspErr.Message,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		con.BLOG: updatedBlog,
	})
}

func BlogDelete(c *gin.Context) {
	blogIdStr := c.Param(con.ID)
	blogId, err := strconv.Atoi(blogIdStr)
	if err != nil {
		rspErr := model.NewBadRequest(err.Error())
		c.JSON(rspErr.Status(), gin.H{
			con.ERROR: rspErr.Message,
		})
		return
	}
	deletedBlogId, err := svc.NewBlogOps(c.Request.Context()).BlogDelete(blogId)
	if err != nil {
		rspErr := model.NewNotFound(con.BLOG, blogIdStr)
		c.JSON(rspErr.Status(), gin.H{
			con.ERROR: rspErr.Message,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		con.BLOG_ID: deletedBlogId,
	})
}
