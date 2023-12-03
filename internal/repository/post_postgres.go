package repository

import (
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"grates/internal/domain"
	"grates/pkg/repository"
	"reflect"
)

type PostRepository struct {
	db *sqlx.DB
}

func NewPostPostgres(db *sqlx.DB) *PostRepository {
	return &PostRepository{db: db}
}

func (p *PostRepository) Create(post domain.Post) (int, error) {
	var id int

	query := fmt.Sprintf("INSERT INTO %s (title, content, users_id) VALUES ($1, $2, $3) RETURNING id;", repository.PostsTable)
	row := p.db.QueryRow(query, post.Title, post.Content, post.UsersId)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (p *PostRepository) Get(postId int) (domain.Post, error) {
	var post domain.Post

	query := fmt.Sprintf("SELECT * FROM %s WHERE id=$1;", repository.PostsTable)
	err := p.db.Get(&post, query, postId)

	return post, err
}

func (p *PostRepository) GetUsersPosts(userId int) ([]domain.Post, error) {
	var posts []domain.Post

	query := fmt.Sprintf(`SELECT * FROM %s WHERE users_id=$1`, repository.PostsTable)
	err := p.db.Select(&posts, query, userId)
	if err != nil {
		return nil, err
	}

	return posts, err
}

func (p *PostRepository) Update(postId int, input domain.PostUpdateInput) error {
	fieldsDb := input.DBifyFields()
	t := reflect.TypeOf(input)
	v := reflect.ValueOf(input)
	argId := 1
	args := make([]interface{}, 0)
	querySet := ""

	// Проходимся по всем полям структуры
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)
		if !value.IsZero() {
			querySet += fmt.Sprintf("%s=$%x, ", fieldsDb[field.Name], argId)
			args = append(args, value.Interface())
			argId += 1
		}
	}
	querySet = querySet[0 : len(querySet)-2]
	args = append(args, postId)

	query := fmt.Sprintf("UPDATE %s SET %s WHERE id=$%x", repository.PostsTable, querySet, argId)
	res, err := p.db.Exec(query, args...)
	if err != nil {
		return err
	}

	// Проверяем изменились ли строки
	count, err := res.RowsAffected()
	if err != nil || count == 0 {
		return errors.New("no rows haven't been updated")
	}

	return nil
}

func (p *PostRepository) Delete(postId int) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id=$1;", repository.PostsTable)
	res, err := p.db.Exec(query, postId)
	if err != nil {
		return err
	}

	// Проверяем удалились ли строки
	count, err := res.RowsAffected()
	if err != nil || count == 0 {
		return errors.New("no rows haven't been deleted")
	}

	return nil
}
