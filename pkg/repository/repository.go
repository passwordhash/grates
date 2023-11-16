package repository

import (
	"github.com/jmoiron/sqlx"
	"grates/internal/entity"
)

type User interface {
	CreateUser(user entity.User) (int, error)
	GetUserByEmail(email string) (entity.User, error)
	GetAllUsers() ([]entity.User, error)
}

type Post interface {
}

type Comment interface {
}

type Repository struct {
	User *UserPostgres
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		User: NewUserPostgres(db),
	}
}
