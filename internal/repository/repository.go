package repository

import (
	"github.com/jmoiron/sqlx"
	"grates/internal/domain"
)

type User interface {
	CreateUser(user domain.User) (int, error)
	GetUser(email string, password string) (domain.User, error)
	GetUserByEmail(email string) (domain.User, error)
	GetAllUsers() ([]domain.User, error)
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
