package domain

import (
	"time"
)

type Post struct {
	Id      int       `json:"id" db:"id" binding:"required"`
	Title   string    `json:"title" db:"title"`
	Content string    `json:"content" db:"content" binding:"required"`
	UsersId int       `json:"users-id" db:"users_id" binding:"required"`
	Date    time.Time `json:"date" db:"date" binding:"required"`
}

type PostUpdateInput struct {
	Title   string `json:"title,omitempty"`
	Content string `json:"content,omitempty"`
}

func (p *PostUpdateInput) DBifyFields() map[string]string {
	return map[string]string{
		"Title":   "title",
		"Content": "content",
	}
}

type Comment struct {
	Id      int       `json:"id" db:"id" binding:"required"`
	Content string    `json:"content" binding:"required"`
	UsersId int       `json:"users-id" db:"users_id" binding:"required"`
	PostsId int       `json:"posts-id" db:"posts_id" binding:"required"`
	Date    time.Time `json:"date" binding:"required"`
}

type Like struct {
	Id      int       `json:"id" db:"id"`
	UsersId int       `json:"users-id" db:"users_id" binding:"required"`
	PostsId int       `json:"posts-id" db:"posts_id" binding:"required"`
	Date    time.Time `json:"date" binding:"required"`
}
