package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"grates/internal/domain"
)

type PostRepository struct {
	db *sqlx.DB
}

func (p *PostRepository) CreatePost(post domain.Post) (int, error) {
	var id int

	query := fmt.Sprintf("INSERT INTO %s (title, content, users_id) VALUES ($1, $2, $3) RETURNING id;", postsTable)
	row := p.db.QueryRow(query, post.Title, post.Content, post.UsersId)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (p *PostRepository) GetPost(postId int) (domain.Post, error) {
	//TODO implement me
	panic("implement me")
}

func (p *PostRepository) GetUsersPosts(postId int) ([]domain.Post, error) {
	var posts []domain.Post

	query := fmt.Sprintf(`SELECT * FROM %s WHERE users_id=$1;`, postsTable)
	err := p.db.Select(&posts, query, postId)

	return posts, err
}

func NewPostPostgres(db *sqlx.DB) *PostRepository {
	return &PostRepository{db: db}
}
