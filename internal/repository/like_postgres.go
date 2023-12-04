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

func (r *LikeRepository) GetPostLikesCount(postId int) (int, error) {
	var count int

	query := fmt.Sprintf("SELECT COUNT(*) FROM %s JOIN %s ON posts.id = likes_posts.posts_id WHERE posts.id=$1", PostsTable, LikesPostsTable)
	err := r.db.Get(&count, query, postId)

	return count, err
}

// GetUsersPostLikesCount возвращает список лайков поста от пользователя
func (r *LikeRepository) GetUsersPostLikesCount(userId, postId int) (int, error) {
	var count int

	query := fmt.Sprintf(`
		SELECT COUNT(*) FROM %s 
		INNER JOIN %s ON posts.id = likes_posts.posts_id 
		WHERE posts.id=$1 AND likes_posts.users_id=$2;
	`, PostsTable, LikesPostsTable)
	err := r.db.Get(&count, query, postId, userId)

	return count, err
}

func (r *LikeRepository) LikePost(userId, postId int) error {
	query := fmt.Sprintf("INSERT INTO %s (users_id, posts_id) VALUES ($1, $2)", LikesPostsTable)

	_, err := r.db.Exec(query, userId, postId)

	return err
}

func (r *LikeRepository) UnlikePost(userId, postId int) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE users_id=$1 AND posts_id=$2", LikesPostsTable)

	_, err := r.db.Exec(query, userId, postId)

	return err
}
