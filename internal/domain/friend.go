package domain

import "time"

type Friend struct {
	Id          int       `db:"id"`
	FromId      int       `db:"from_id"`
	ToId        int       `db:"to_id"`
	IsConfirmed bool      `db:"is_confirmed"`
	SendAt      time.Time `db:"send_at"`
}

//type FriendInput struct {
//	FromId int `json:"from_id" db:"from_id"`
//	ToId   int `json:"to_id" db:"to_id"`
//}
