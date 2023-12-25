package service

import (
	"context"
	"fmt"
	"sxolla-rest-api/config"
	"sxolla-rest-api/ent"
	"sxolla-rest-api/ent/enttest"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/suite"
)

type BlogEntTestSuite struct {
	suite.Suite
	entClient *ent.Client
}

func (suite *BlogEntTestSuite) SetupSuite() {
	// set StartingNumber to one
	fmt.Println(">>> From Setup Blog App Test Suite")

}

func TestBlogEntTestSuite(t *testing.T) {
	suite.Run(t, &BlogEntTestSuite{})
}

func (suite *BlogEntTestSuite) BeforeTest(suiteName, testName string) {
	fmt.Println(">>> Before Test", testName)
	// Initialize an Ent client that uses an in memory SQLite db.
	suite.entClient = enttest.Open(suite.T(), "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")

	config.SetClient(suite.entClient)
}

func (suite *BlogEntTestSuite) AfterTest(suiteName, testName string) {
	fmt.Println(">>> After Test", testName)
	// 	// cleanup resource
	suite.entClient.Close()
}

func (suite *BlogEntTestSuite) TestBlogCreate() {
	newBlog := &ent.Blog{Title: "test title", Content: "test content"}
	_, err := NewBlogOps(context.Background()).BlogCreate(newBlog)
	suite.Nil(err)
	blog, err := NewBlogOps(context.Background()).BlogGetByID(1)
	suite.Equal(newBlog.Title, blog.Title)
	suite.Equal(newBlog.Content, blog.Content)
}

func (suite *BlogEntTestSuite) TestBlogGet() {
	newBlog := &ent.Blog{Title: "test title", Content: "test content"}
	blog, err := NewBlogOps(context.Background()).BlogGetByID(1)
	suite.NotNil(err)

	_, err = NewBlogOps(context.Background()).BlogCreate(newBlog)
	blog, err = NewBlogOps(context.Background()).BlogGetByID(1)
	suite.Equal(newBlog.Title, blog.Title)
	suite.Equal(newBlog.Content, blog.Content)
}

func (suite *BlogEntTestSuite) TestBlogUpdate() {
	newBlog := &ent.Blog{Title: "test title", Content: "test content"}
	_, err := NewBlogOps(context.Background()).BlogCreate(newBlog)
	suite.Nil(err)

	newBlog.Content = "update content"
	newBlog.ID = 1
	_, err = NewBlogOps(context.Background()).BlogUpdate(newBlog)
	suite.Nil(err)
	blog, err := NewBlogOps(context.Background()).BlogGetByID(1)
	suite.Nil(err)
	suite.Equal(newBlog.Title, blog.Title)
	suite.Equal(newBlog.Content, blog.Content)
}

func (suite *BlogEntTestSuite) TestBlogDelete() {
	newBlog := &ent.Blog{Title: "test title", Content: "test content"}
	_, err := NewBlogOps(context.Background()).BlogCreate(newBlog)
	suite.Nil(err)

	deletedBlogId, err := NewBlogOps(context.Background()).BlogDelete(1)
	suite.Nil(err)
	suite.Equal(1, deletedBlogId)

}
