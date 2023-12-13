package domain

type User struct {
	Id       int    `json:"id" db:"id"`
	Name     string `json:"name" binding:"required"`
	Surname  string `json:"surname"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" db:"password_hash"`
}

type AuthEmail struct {
	Id      int    `db:"id"`
	UsersId int    `db:"users_id"`
	Hash    string `db:"hash"`
}

func (u *User) IsEmtpty() bool {
	var emptyUser User
	return emptyUser == *u
}
