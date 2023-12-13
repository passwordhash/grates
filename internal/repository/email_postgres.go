package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
)

type EmailRepository struct {
	db *sqlx.DB
}

func (e *EmailRepository) ReplaceEmail(userId int, hash string) error {
	tx, err := e.db.Beginx()
	if err != nil {
		return err
	}

	query := fmt.Sprintf("DELETE FROM %s WHERE users_id=$1", AuthEmailsTable)
	_, err = tx.Exec(query, userId)
	if err != nil {
		tx.Rollback()
		return err
	}

	query = fmt.Sprintf("INSERT INTO %s (users_id, hash) VALUES ($1, $2)", AuthEmailsTable)
	_, err = tx.Exec(query, userId, hash)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func NewEmailRepository(db *sqlx.DB) *EmailRepository {
	return &EmailRepository{db: db}
}
