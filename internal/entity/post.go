package entity

import "time"

type Post struct {
	Id      int       `json:"id" db:"id" binding:"required"`
	Title   string    `json:"title"`
	Content string    `json:"content" binding:"required"`
	UsersId string    `json:"users-id" db:"users_id" binding:"required"`
	Date    time.Time `json:"date" binding:"required"`
}

type Comment struct {
	Id      int       `json:"id" db:"id" binding:"required"`
	Content string    `json:"content" binding:"required"`
	UsersId string    `json:"users-id" db:"users_id" binding:"required"`
	PostsId string    `json:"posts-id" db:"posts_id" binding:"required"`
	Date    time.Time `json:"date" binding:"required"`
}

type Like struct {
	Id      int       `json:"id" db:"id"`
	UsersId string    `json:"users-id" db:"users_id" binding:"required"`
	PostsId string    `json:"posts-id" db:"posts_id" binding:"required"`
	Date    time.Time `json:"date" binding:"required"`
}
