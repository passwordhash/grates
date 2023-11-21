package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"grates/internal/domain"
)

type UserPostgres struct {
	db *sqlx.DB
}

func NewUserPostgres(db *sqlx.DB) *UserPostgres {
	return &UserPostgres{db: db}
}

// CreateUser при успешном срабатывании, возвращает id созданного пользователя
func (r *UserPostgres) CreateUser(user domain.User) (int, error) {
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

// GetUser возвращает domain.User, если пользователь с такой почтой и паролем сущетсвует
func (r *UserPostgres) GetUser(email, password string) (domain.User, error) {
	var user domain.User

	query := fmt.Sprintf("SELECT * FROM %s WHERE email=$1 AND password_hash=$2", usersTable)
	err := r.db.Get(&user, query, email, password)

	return user, err
}

// GetUserByEmail : возвращет пользователя
func (r *UserPostgres) GetUserByEmail(email string) (domain.User, error) {
	var user domain.User

	query := fmt.Sprintf("SELECT * FROM %s WHERE email=$1;", usersTable)
	err := r.db.Get(&user, query, email)

	return user, err
}

func (r *UserPostgres) GetAllUsers() ([]domain.User, error) {
	var users []domain.User

	query := fmt.Sprintf("SELECT * FROM %s;",
		usersTable)
	err := r.db.Select(&users, query)

	return users, err
}
