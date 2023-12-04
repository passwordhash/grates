package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"grates/internal/domain"
)

type CommentRepository struct {
	db *sqlx.DB
}

func (c CommentRepository) Create(comment domain.CommentCreateInput) (int, error) {
	var id int

	query := fmt.Sprintf(`INSERT INTO %s (content, users_id, posts_id)`, CommentsTable)
	query += fmt.Sprintf(`VALUES ($1, $2, $3) RETURNING id;`)

	row := c.db.QueryRow(query, comment.Content, comment.UserId, comment.PostId)

	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (c CommentRepository) GetPostComments(postId int) ([]domain.Comment, error) {
	var comments []domain.Comment

	query := fmt.Sprintf(`SELECT * FROM %s WHERE posts_id=$1;`, CommentsTable)
	err := c.db.Select(&comments, query, postId)

	return comments, err
}

func (c CommentRepository) Update(userId, commentId int, newComment domain.CommentUpdateInput) error {
	query := fmt.Sprintf(`UPDATE %s SET content=$1 WHERE id=$2 AND users_id=$3;`,
		CommentsTable)
	_, err := c.db.Exec(query, newComment.Content, commentId, userId)

	return err
}

func (c CommentRepository) Delete(userId, commentId int) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE id=$1 AND users_id=$2;`, CommentsTable)
	_, err := c.db.Exec(query, commentId, userId)

	return err
}

func NewCommentRepository(db *sqlx.DB) *CommentRepository {
	return &CommentRepository{db: db}
}
