package domain

import (
	"grates/pkg/repository"
	"time"
)

type Post struct {
	Id      int       `json:"id" db:"id" binding:"required" example:"732436"`
	Title   string    `json:"title" db:"title" example:"Post title"`
	Content string    `json:"content" db:"content" binding:"required" example:"Occaecat quis officia pariatur non aliquip culpa id elit amet sit occaecat ex sunt ullamco duis reprehenderit in esse. Culpa minim nulla pariatur voluptate ea proident dolor velit eu do labore ut."`
	UsersId int       `json:"users-id" db:"users_id" binding:"required" example:"6296"`
	Date    time.Time `json:"date" db:"date" binding:"required" example:"2021-01-01T00:00:00Z"`

	Comments   []Comment `json:"comments" db:"-"`
	LikesCount int       `json:"likes-count" db:"likes_count"`
}

type PostUpdateInput struct {
	Title   string `json:"title,omitempty"`
	Content string `json:"content,omitempty"`
	repository.DBifyable
}

// DBifyFields возвращает поля, которые нужно обновить в БД.
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

type CommentCreateInput struct {
	Content string `json:"content"`
	UserId  int    `json:"-"`
	PostId  int    `json:"-"`
}

type CommentUpdateInput struct {
	Content string `json:"content,omitempty"`
}

type Like struct {
	Id      int       `json:"id" db:"id"`
	UsersId int       `json:"users-id" db:"users_id" binding:"required"`
	PostsId int       `json:"posts-id" db:"posts_id" binding:"required"`
	Date    time.Time `json:"date" binding:"required"`
}
