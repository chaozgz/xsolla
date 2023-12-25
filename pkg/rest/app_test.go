package rest

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sxolla-rest-api/config"
	"sxolla-rest-api/ent"
	"sxolla-rest-api/ent/enttest"
	"sxolla-rest-api/pkg/rest/handler"
	"sxolla-rest-api/service"
	"testing"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/suite"
)

type BlogRestAppTestSuite struct {
	suite.Suite
	entClient  *ent.Client
	blogSample *ent.Blog
}

func (suite *BlogRestAppTestSuite) SetupSuite() {
	// set StartingNumber to one
	fmt.Println(">>> From Setup Blog App Test Suite")

	suite.blogSample = &ent.Blog{Title: "test title", Content: "test content"}

}

func (suite *BlogRestAppTestSuite) BeforeTest(suiteName, testName string) {
	fmt.Println(">>> Before Test", testName)
	// Initialize an Ent client that uses an in memory SQLite db.
	suite.entClient = enttest.Open(suite.T(), "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")

	config.SetClient(suite.entClient)
}

func (suite *BlogRestAppTestSuite) AfterTest(suiteName, testName string) {
	fmt.Println(">>> After Test", testName)
	// 	// cleanup resource
	suite.entClient.Close()
}

func TestBlogRestAppTestSuite(t *testing.T) {
	suite.Run(t, &BlogRestAppTestSuite{})
}

func (suite *BlogRestAppTestSuite) TestBlogGetAll() {
	router := gin.Default()
	handler.RegisterBlogRoutes("api/v1", router)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/posts", nil)
	req = req.WithContext(context.Background())
	router.ServeHTTP(w, req)
	suite.Equal(http.StatusOK, w.Code)
	suite.Equal(`{"blogs":[]}`, w.Body.String())
	// create
	err := suite.createBlog()
	suite.Nil(err)
	// read
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/v1/posts", nil)
	req = req.WithContext(context.Background())
	router.ServeHTTP(w, req)
	suite.Equal(http.StatusOK, w.Code)
	blogsRsp, err := convertJsonToBlogsRsp(w.Body.Bytes())
	suite.Nil(err)
	suite.Equal(1, len(blogsRsp.Blogs))
	suite.Equal(suite.blogSample.Title, blogsRsp.Blogs[0].Title)
	suite.Equal(suite.blogSample.Content, blogsRsp.Blogs[0].Content)
}

func (suite *BlogRestAppTestSuite) TestBlogGetByID() {
	router := gin.Default()
	handler.RegisterBlogRoutes("api/v1", router)
	// create
	err := suite.createBlog()
	suite.Nil(err)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/posts/1", nil)
	req = req.WithContext(context.Background())
	router.ServeHTTP(w, req)
	suite.Equal(http.StatusOK, w.Code)

	blogRsp, err := convertJsonToBlogRsp(w.Body.Bytes())
	suite.Nil(err)
	suite.Equal(suite.blogSample.Title, blogRsp.Blog.Title)
	suite.Equal(suite.blogSample.Content, blogRsp.Blog.Content)
}

func (suite *BlogRestAppTestSuite) TestBlogGetByID_NOT_FOUND() {
	router := gin.Default()
	handler.RegisterBlogRoutes("api/v1", router)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/posts/1", nil)
	req = req.WithContext(context.Background())
	router.ServeHTTP(w, req)
	suite.Equal(http.StatusNotFound, w.Code)
}

func (suite *BlogRestAppTestSuite) TestBlogUpdate() {
	router := gin.Default()
	handler.RegisterBlogRoutes("api/v1", router)

	// create
	err := suite.createBlog()
	suite.Nil(err)
	// update
	newBlog := &ent.Blog{Title: "new title", Content: "new content"}

	blogByte, _ := json.Marshal(newBlog)
	bodyReader := bytes.NewReader(blogByte)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/api/v1/posts/1", bodyReader)
	req = req.WithContext(context.Background())
	router.ServeHTTP(w, req)
	suite.Equal(http.StatusOK, w.Code)
	// read
	result, err := suite.readBlogById(1)
	suite.Nil(err)
	suite.Equal(newBlog.Title, result.Title)
	suite.Equal(newBlog.Content, result.Content)
}

func (suite *BlogRestAppTestSuite) TestBlogUpdate_REQ_BODY_EMPTY() {
	router := gin.Default()
	handler.RegisterBlogRoutes("api/v1", router)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/api/v1/posts/1", nil)
	req = req.WithContext(context.Background())
	router.ServeHTTP(w, req)
	suite.Equal(http.StatusBadRequest, w.Code)
	suite.Equal(`{"error":"Bad request. Reason: request body can not be empty"}`, w.Body.String())
}

func (suite *BlogRestAppTestSuite) TestBlogUpdate_NOT_FOUND() {
	router := gin.Default()
	handler.RegisterBlogRoutes("api/v1", router)
	// update
	newBlog := &ent.Blog{Title: "new title", Content: "new content"}

	blogByte, _ := json.Marshal(newBlog)
	bodyReader := bytes.NewReader(blogByte)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/api/v1/posts/1", bodyReader)
	req = req.WithContext(context.Background())
	router.ServeHTTP(w, req)
	suite.Equal(http.StatusNotFound, w.Code)
}

func (suite *BlogRestAppTestSuite) TestBlogDelete() {
	router := gin.Default()
	handler.RegisterBlogRoutes("api/v1", router)

	// create
	err := suite.createBlog()
	suite.Nil(err)
	// delete
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/api/v1/posts/1", nil)
	req = req.WithContext(context.Background())
	router.ServeHTTP(w, req)
	suite.Equal(http.StatusOK, w.Code)
	// read
	_, err = suite.readBlogById(1)
	suite.NotEmpty(err)
}

func (suite *BlogRestAppTestSuite) TestBlogDelete_NOT_FOUND() {
	router := gin.Default()
	handler.RegisterBlogRoutes("api/v1", router)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/api/v1/posts/1", nil)
	req = req.WithContext(context.Background())
	router.ServeHTTP(w, req)
	suite.Equal(http.StatusNotFound, w.Code)
}

func (suite *BlogRestAppTestSuite) TestHealthCheck() {
	router := gin.Default()
	handler.RegisterHealthRoutes("api/v1", router)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/healthcheck", nil)
	req = req.WithContext(context.Background())
	router.ServeHTTP(w, req)
	suite.Equal(http.StatusOK, w.Code)
	suite.Equal(`{"status":"UP"}`, w.Body.String())
}

func (suite *BlogRestAppTestSuite) createBlog() error {
	_, err := service.NewBlogOps(context.Background()).BlogCreate(suite.blogSample)
	if err != nil {
		return err
	}
	return nil
}

func (suite *BlogRestAppTestSuite) readBlogById(id int) (*ent.Blog, error) {
	blog, err := service.NewBlogOps(context.Background()).BlogGetByID(id)
	if err != nil {
		return nil, err
	}
	return blog, nil
}

type BlogsRsp struct {
	Blogs []ent.Blog `json:"blogs"`
}

type BlogRsp struct {
	Blog ent.Blog `json:"blog"`
}

func convertJsonToBlogsRsp(body []byte) (*BlogsRsp, error) {
	var blogsRsp *BlogsRsp
	err := json.Unmarshal(body, &blogsRsp)
	if err != nil {
		return nil, err
	}
	return blogsRsp, nil
}

func convertJsonToBlogRsp(body []byte) (*BlogRsp, error) {
	var blogRsp BlogRsp
	err := json.Unmarshal(body, &blogRsp)
	if err != nil {
		return nil, err
	}
	return &blogRsp, nil
}
