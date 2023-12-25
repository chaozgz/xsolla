package service

import (
	"context"
	"sxolla-rest-api/config"
	"sxolla-rest-api/ent"
	"time"
)

type BlogOps struct {
	ctx    context.Context
	client *ent.Client
}

func NewBlogOps(ctx context.Context) *BlogOps {
	return &BlogOps{
		ctx:    ctx,
		client: config.GetClient(),
	}
}

func (r *BlogOps) BlogGetAll() ([]*ent.Blog, error) {
	blogs, err := r.client.Blog.Query().All(r.ctx)
	if err != nil {
		return nil, err
	}

	return blogs, nil
}

func (r *BlogOps) BlogPagination(offset, limit int) ([]*ent.Blog, error) {
	blogs, err := r.client.Blog.Query().Limit(limit).Offset(offset).All(r.ctx)
	if err != nil {
		return nil, err
	}

	return blogs, nil
}

func (r *BlogOps) BlogGetByID(id int) (*ent.Blog, error) {

	blog, err := r.client.Blog.Get(r.ctx, id)
	if err != nil {
		return nil, err
	}

	return blog, nil
}

func (r *BlogOps) BlogCreate(newBlog *ent.Blog) (*ent.Blog, error) {

	newCreatedBlog, err := r.client.Blog.Create().
		SetTitle(newBlog.Title).
		SetContent(newBlog.Content).
		Save(r.ctx)

	if err != nil {
		return nil, err
	}

	return newCreatedBlog, nil
}

func (r *BlogOps) BlogUpdate(blog *ent.Blog) (*ent.Blog, error) {

	updatedUser, err := r.client.Blog.UpdateOneID(blog.ID).
		SetTitle(blog.Title).
		SetContent(blog.Content).
		SetUpdatedAt(time.Now()).
		Save(r.ctx)
	if err != nil {
		return nil, err
	}

	return updatedUser, nil
}

func (r *BlogOps) BlogDelete(id int) (int, error) {
	err := r.client.Blog.
		DeleteOneID(id).
		Exec(r.ctx)

	if err != nil {
		return 0, err
	}

	return id, nil
}
