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

func (e *EmailRepository) ConfirmEmail(hash string) error {
	// TODO: может стоит удалять запись из таблицы auth_emails после подтверждения?
	query := fmt.Sprintf(`
	UPDATE %s
	    SET is_confirmed = TRUE
	FROM %s
	WHERE users.id = auth_emails.users_id
	AND auth_emails.hash = $1
	AND users.is_confirmed = FALSE
`, UsersTable, AuthEmailsTable)
	res, err := e.db.Exec(query, hash)
	if err != nil {
		return err
	}
	changes, _ := res.RowsAffected()

	if changes == 0 {
		return NoChangesErr
	}

	return err
}

func (e *EmailRepository) IsConfirmed(hash string) (bool, error) {
	var isConfirmed bool
	query := fmt.Sprintf(`
	SELECT is_confirmed FROM %s
	INNER JOIN %s ON %s.id = %s.users_id
	WHERE %s.hash=$1
	`, UsersTable, AuthEmailsTable, UsersTable, AuthEmailsTable, AuthEmailsTable)

	err := e.db.Get(&isConfirmed, query, hash)

	return isConfirmed, err
}

func NewEmailRepository(db *sqlx.DB) *EmailRepository {
	return &EmailRepository{db: db}
}
