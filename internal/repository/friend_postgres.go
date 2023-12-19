package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"grates/internal/domain"
)

type FriendRepository struct {
	db *sqlx.DB
}

func NewFriendRepository(db *sqlx.DB) *FriendRepository {
	return &FriendRepository{db: db}
}

func (r *FriendRepository) Get(id1, id2 int) (domain.Friend, error) {
	var friend domain.Friend

	query := fmt.Sprintf("SELECT * FROM %s WHERE (from_id=$1 AND to_id=$2) OR (from_id=$2 AND to_id=$1)", FriendsTable)
	err := r.db.Get(&friend, query, id1, id2)

	return friend, err
}

// GetFriendUsers возвращает список друзей пользователя.
func (r *FriendRepository) FriendUsers(userId int) ([]domain.User, error) {
	var friends []domain.User

	query := fmt.Sprintf(`
	SELECT * FROM users
	WHERE %s.id IN
		(SELECT to_id as id FROM %s WHERE from_id=$1 AND is_confirmed=true
		UNION
		SELECT from_id as id FROM %s WHERE to_id=$1 AND is_confirmed=true)
	`, UsersTable, FriendsTable, FriendsTable)
	err := r.db.Select(&friends, query, userId)

	return friends, err
}

func (r *FriendRepository) FriendUsersIds(userId int) ([]int, error) {
	var friends []int

	query := fmt.Sprintf(`
	SELECT id FROM users
	WHERE %s.id IN
		(SELECT to_id as id FROM %s WHERE from_id=$1 AND is_confirmed=true
		UNION
		SELECT from_id as id FROM %s WHERE to_id=$1 AND is_confirmed=true)
	`, UsersTable, FriendsTable, FriendsTable)
	err := r.db.Select(&friends, query, userId)

	return friends, err
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

func (r *FriendRepository) AcceptFriendRequest(fromId, toId int) error {
	query := fmt.Sprintf("UPDATE %s SET is_confirmed=true WHERE from_id=$1 AND to_id=$2 ", FriendsTable)
	_, err := r.db.Exec(query, fromId, toId)

	return err
}

// Unfriend применяется в случае, если пользователи были друзьями (т. е. is_confirmed=true).
// В этом случае удаляется запись из таблицы и добавляется обратно, но с другими id, как будто пользователь,
// которого удалили, отправил заявку (т. е. is_confirmed=false, и from_id и to_id меняются).
func (r *FriendRepository) Unfriend(userId, friendId int) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}

	// TODO: переписать на один запрос
	query := fmt.Sprintf("DELETE FROM %s WHERE (from_id=$1 AND to_id=$2) OR (from_id=$2 AND to_id=$1)", FriendsTable)
	_, err = tx.Exec(query, userId, friendId)
	if err != nil {
		return fmt.Errorf("error deleting friend: %w", err)
	}

	query = fmt.Sprintf("INSERT INTO %s (from_id, to_id) VALUES ($1, $2)", FriendsTable)
	_, err = tx.Exec(query, friendId, userId)
	if err != nil {
		return fmt.Errorf("error adding friend: %w", err)
	}

	return tx.Commit()
}

// Decline применяется в случае, если пользователи не были друзьями (т. е. is_confirmed=false).
// В этом случае удаляется запись из таблицы.
func (r *FriendRepository) Decline(userId, friendId int) error {
	//query := fmt.Sprintf("DELETE FROM %s WHERE ((from_id=$1 AND to_id=$2) OR (from_id=$2 AND to_id=$1)) AND is_confirmed=false", FriendsTable)
	query := fmt.Sprintf("DELETE FROM %s WHERE (from_id=$1 AND to_id=$2) AND is_confirmed=false", FriendsTable)
	_, err := r.db.Exec(query, userId, friendId)

	return err
}
