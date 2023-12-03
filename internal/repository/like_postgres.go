package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
)

type LikeRepository struct {
	db *sqlx.DB
}

func NewLikeRepository(db *sqlx.DB) *LikeRepository {
	return &LikeRepository{db: db}
}

func (r *LikeRepository) LikePost(userId, postId int) error {
	query := fmt.Sprintf("INSERT INTO likes (user_id, post_id) VALUES ($1, $2)", LikesPostsTable)

	_, err := r.db.Exec(query, userId, postId)

	return err
}

func (r *LikeRepository) UnlikePost(userId, postId int) error {
	query := fmt.Sprintf("DELETE FROM likes WHERE user_id=$1 AND post_id=$2", LikesPostsTable)

	_, err := r.db.Exec(query, userId, postId)

	return err
}
