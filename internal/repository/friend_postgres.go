package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
)

type FriendRepository struct {
	db *sqlx.DB
}

func NewFriendRepository(db *sqlx.DB) *FriendRepository {
	return &FriendRepository{db: db}
}

func (r *FriendRepository) FriendRequest(fromId, toId int) error {
	var count int

	query := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE (from_id=$1 AND to_id=$2) OR (from_id=$2 AND to_id=$1)", FriendsTable)
	err := r.db.Get(&count, query, fromId, toId)
	if err != nil {
		return err
	}

	if count > 0 {
		return fmt.Errorf("friend request already sent")
	}

	query = fmt.Sprintf("INSERT INTO %s (from_id, to_id) VALUES ($1, $2)", FriendsTable)
	_, err = r.db.Exec(query, fromId, toId)

	return err
}

func (r *FriendRepository) AcceptFriendRequest(id1, id2 int) error {
	query := fmt.Sprintf("UPDATE %s SET is_confirmed=true WHERE (from_id=$1 AND to_id=$2) OR (from_id=$2 AND to_id=$1)", FriendsTable)
	_, err := r.db.Exec(query, id1, id2)

	return err
}
