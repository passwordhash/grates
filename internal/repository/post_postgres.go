package repository

import (
	"errors"
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

func (p *PostRepository) DeletePostById(postId int) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id=$1;", postsTable)
	res, err := p.db.Exec(query, postId)
	if err != nil {
		return err
	}

	// Проверяем удалились дли строки
	count, err := res.RowsAffected()
	if err != nil || count == 0 {
		return errors.New("no rows haven't been deleted")
	}

	return nil
}

func (p *PostRepository) UpdatePost(newPost domain.Post) error {
	//TODO implement me
	panic("implement me")
}
