package domain

import (
	"grates/pkg/utils"
	"time"
)

type Gnd = string

const (
	Male   Gnd = "M"
	Female     = "F"
	NotSet     = "N"
)

type User struct {
	Id          int        `json:"id" db:"id"`
	Name        string     `json:"name" binding:"required"`
	Surname     string     `json:"surname"`
	Email       string     `json:"email" binding:"required"`
	Password    string     `json:"password" db:"password_hash"`
	IsConfirmed bool       `json:"is_confirmed" db:"is_confirmed" default:"false"`
	Gender      Gnd        `json:"gender" db:"gender" default:"N"`
	BirthDate   *time.Time `json:"birth_date" db:"birth_date"`
	Status      string     `json:"status" db:"status" default:""`
	IsDeleted   bool       `json:"is_deleted" db:"is_deleted" default:"false"`
}

func (u *User) GetAge() int {
	return time.Now().Year() - u.BirthDate.Year()
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
