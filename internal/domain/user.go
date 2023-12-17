package domain

import (
	"database/sql"
	"grates/pkg/utils"
)

type Gnd = string

const (
	Male   Gnd = "M"
	Female     = "F"
	NotSet     = "N"
)

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

type UserSignUpInput struct {
	Name     string `json:"name" binding:"required"`
	Surname  string `json:"surname"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type ProfileUpdateInput struct {
	Name      string     `json:"name"`
	Surname   string     `json:"surname"`
	Gender    Gnd        `json:"gender" db:"gender" default:"M" enum:"M,F,N"`
	BirthDate utils.Date `json:"birth_date" db:"birth_date" time_format:"2006-01-02" swaggertype:"string" example:"2006-01-02"`
	Status    string     `json:"status" db:"status" default:""`
}

func (p *ProfileUpdateInput) DBifyFields() map[string]string {
	return map[string]string{
		"Name":      "name",
		"Surname":   "surname",
		"Gender":    "gender",
		"BirthDate": "birth_date",
		"Status":    "status",
	}
}

type AuthEmail struct {
	Id      int    `db:"id"`
	UsersId int    `db:"users_id"`
	Hash    string `db:"hash"`
}

func (u *User) IsEmtpty() bool {
	return u.Id == 0 || u.Email == ""
}
