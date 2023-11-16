package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"grates/internal/entity"
)

type UserPostgres struct {
	db *sqlx.DB
}

func NewUserPostgres(db *sqlx.DB) *UserPostgres {
	return &UserPostgres{db: db}
}

func (r *UserPostgres) GetAllUsers() ([]entity.User, error) {
	var users []entity.User

	query := fmt.Sprintf("SELECT * FROM %s;",
		usersTable)
	err := r.db.Select(&users, query)

	return users, err
}
