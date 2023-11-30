package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"grates/internal/domain"
	"grates/pkg/repository"
)

type CommentRepository struct {
	db *sqlx.DB
}

func (c CommentRepository) Create(comment domain.CommentCreateInput) (int, error) {
	var id int

	query := fmt.Sprintf(`INSERT INTO %s (content, users_id, posts_id)`, repository.CommentsTable)
	query += fmt.Sprintf(`VALUES ($1, $2, $3) RETURNING id;`)

	row := c.db.QueryRow(query, comment.Content, comment.UserId, comment.PostId)

	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (c CommentRepository) GetPostComments(postId int) ([]domain.Comment, error) {
	var comments []domain.Comment

	query := fmt.Sprintf(`SELECT * FROM %s WHERE posts_id=$1;`, repository.CommentsTable)
	err := c.db.Select(&comments, query, postId)

	return comments, err
}

func (c CommentRepository) Update(id int, newComment domain.CommentCreateInput) error {
	//TODO implement me
	panic("implement me")
}

func (c CommentRepository) Delete(id int) error {
	//TODO implement me
	panic("implement me")
}

func NewCommentRepository(db *sqlx.DB) *CommentRepository {
	return &CommentRepository{db: db}
}
