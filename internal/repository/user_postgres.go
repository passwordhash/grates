package repository

import (
	"fmt"
	"grates/internal/domain"
)

// CreateUser при успешном срабатывании, возвращает id созданного пользователя
func (r *UserRepository) CreateUser(user domain.User) (int, error) {
	var id int

	query := fmt.Sprintf(`INSERT INTO %s (name, surname, email, password_hash)
						VALUES ($1, $2, $3, $4)
						RETURNING id;`, UsersTable)
	// TRY
	row := r.db.QueryRow(query, user.Name, user.Surname, user.Email, user.Password)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

// GetUser возвращает domain.User, если пользователь с такой почтой и паролем сущетсвует
func (r *UserRepository) GetUser(email, password string) (domain.User, error) {
	var user domain.User

	query := fmt.Sprintf("SELECT * FROM %s WHERE email=$1 AND password_hash=$2", UsersTable)
	err := r.db.Get(&user, query, email, password)

	return user, err
}

func (r *UserRepository) GetUserById(id int) (domain.User, error) {
	var user domain.User

	query := fmt.Sprintf("SELECT * FROM %s WHERE id=$1;", UsersTable)
	err := r.db.Get(&user, query, id)

	return user, err
}

// GetUserByEmail : возвращет пользователя
func (r *UserRepository) GetUserByEmail(email string) (domain.User, error) {
	var user domain.User

	query := fmt.Sprintf("SELECT * FROM %s WHERE email=$1;", UsersTable)
	err := r.db.Get(&user, query, email)

	return user, err
}

func (r *UserRepository) GetAllUsers() ([]domain.User, error) {
	var users []domain.User

	query := fmt.Sprintf("SELECT * FROM %s;",
		UsersTable)
	err := r.db.Select(&users, query)

	return users, err
}
