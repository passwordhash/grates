package repository

import (
	"errors"
	"fmt"
	"grates/internal/domain"
	"reflect"
)

// CreateUser при успешном срабатывании, возвращает id созданного пользователя
func (r *UserRepository) CreateUser(user domain.UserSignUpInput) (int, error) {
	var id int

	query := fmt.Sprintf(`INSERT INTO %s (name, surname, email, password_hash)
						VALUES ($1, $2, $3, $4)
						RETURNING id;`, UsersTable)

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

func (r *UserRepository) UpdateProfile(userId int, input domain.ProfileUpdateInput) error {
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
			args = append(args, fmt.Sprintf("%s", value.Interface()))
			argId += 1
		}
	}

	querySet = querySet[0 : len(querySet)-2]
	args = append(args, userId)

	query := fmt.Sprintf("UPDATE %s SET %s WHERE id=$%x", UsersTable, querySet, argId)
	res, err := r.db.Exec(query, args...)
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
