package domain

import (
	"database/sql"
	"grates/pkg/repository"
	"grates/pkg/utils"
)

type Gnd = string

const (
	Male   Gnd = "M"
	Female     = "F"
	NotSet     = "N"
)

// User представляет собой пользователя.
type User struct {
	Id          int          `db:"id"`
	Name        string       `db:"name" binding:"required"`
	Surname     string       `db:"surname"`
	Email       string       `db:"email" binding:"required"`
	Password    string       `db:"password_hash"`
	IsConfirmed bool         `db:"is_confirmed" default:"false"`
	Gender      Gnd          `db:"gender" default:"N"`
	BirthDate   sql.NullTime `db:"birth_date"`
	Status      string       `db:"status" default:""`
	IsDeleted   bool         `db:"is_deleted" default:"false"`
}

// UserResponse представляет собой пользователя в пригодном для ответа виде.
type UserResponse struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Surname     string `json:"surname"`
	Email       string `json:"email"`
	IsConfirmed bool   `json:"is_confirmed"`
	BirthDate   string `json:"birth_date" example:"2006-01-02"`
	Gender      Gnd    `json:"gender" default:"N" enum:"M,F,N"`
	Status      string `json:"status"`
	IsDeleted   bool   `json:"is_deleted"`
}

// ToResponse возвращает пользователя в пригодном для ответа виде.
func (u *User) ToResponse() UserResponse {
	dateS := utils.Date{Time: u.BirthDate.Time}.String()
	if !u.BirthDate.Valid {
		dateS = ""
	}
	return UserResponse{
		Id:          u.Id,
		Name:        u.Name,
		Surname:     u.Surname,
		Email:       u.Email,
		IsConfirmed: u.IsConfirmed,
		Gender:      u.Gender,
		BirthDate:   dateS,
		Status:      u.Status,
		IsDeleted:   u.IsDeleted,
	}
}

// UserSignUpInput представляет собой данные, необходимые для регистрации пользователя.
type UserSignUpInput struct {
	Name     string `json:"name" binding:"required"`
	Surname  string `json:"surname"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// ProfileUpdateInput представляет собой данные, необходимые для обновления профиля пользователя.
type ProfileUpdateInput struct {
	Name                 string     `json:"name"`
	Surname              string     `json:"surname"`
	Gender               Gnd        `json:"gender" db:"gender" default:"M" enum:"M,F,N"`
	BirthDate            utils.Date `json:"birth_date" db:"birth_date" time_format:"2006-01-02" swaggertype:"string" example:"2006-01-02"`
	Status               string     `json:"status" db:"status" default:""`
	repository.DBifyable `json:"-"`
}

// DBifyFields возращает соответствие полей структуры и полей в БД в виде отображения
func (p *ProfileUpdateInput) DBifyFields() map[string]string {
	return map[string]string{
		"Name":      "name",
		"Surname":   "surname",
		"Gender":    "gender",
		"BirthDate": "birth_date",
		"Status":    "status",
	}
}

// IsEmtpty	пользователь пустой, если у него нет id или email
func (u *User) IsEmtpty() bool {
	return u.Id == 0 || u.Email == ""
}

func UserListToResponse(users []User) []UserResponse {
	var res []UserResponse
	for _, u := range users {
		res = append(res, u.ToResponse())
	}

	return res
}

type AuthEmail struct {
	Id      int    `db:"id"`
	UsersId int    `db:"users_id"`
	Hash    string `db:"hash"`
}
