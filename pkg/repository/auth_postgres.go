package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"grates/internal/entity"
)

type AuthPostgres struct {
	db *sqlx.DB
}

func NewAuthPostgres(db *sqlx.DB) *AuthPostgres {
	return &AuthPostgres{db: db}
}

// CreateUser : при успешном срабатывании, возвращает id созданного пользователя
func (r *AuthPostgres) CreateUser(user entity.User) (int, error) {
	var id int

	query := fmt.Sprintf(`INSERT INTO %s (name, surname, email, password_hash)
						VALUES ($1, $2, $3, $4)
						RETURNING id;`, usersTable)
	// TRY
	row := r.db.QueryRow(query, user.Name, user.Surname, user.Email, user.Password)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}
