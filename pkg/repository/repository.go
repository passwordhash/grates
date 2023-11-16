package repository

import (
	"github.com/jmoiron/sqlx"
	"grates/internal/entity"
)

type Authorization interface {
	CreateUser(user entity.User) (int, error)
}

type User interface {
	GetAllUsers() ([]entity.User, error)
}

type Post interface {
}

type Comment interface {
}

type Repository struct {
	Authorization *AuthPostgres
	User          *UserPostgres
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
		User:          NewUserPostgres(db),
	}
}
