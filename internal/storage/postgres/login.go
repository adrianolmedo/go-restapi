package postgres

import (
	"database/sql"

	"github.com/adrianolmedo/go-restapi-practice/internal/domain"
)

type LoginRepository struct {
	db *sql.DB
}

func NewLoginRepository(db *sql.DB) *LoginRepository {
	return &LoginRepository{
		db: db,
	}
}

func (r LoginRepository) UserByLogin(email, password string) error {
	stmt, err := r.db.Prepare("SELECT email, password FROM users WHERE email = $1 AND password = $2")
	if err != nil {
		return err
	}
	defer stmt.Close()

	result, err := stmt.Exec(email, password)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return domain.ErrUserNotFound
	}

	return nil
}

/*func (r LoginRepository) UserByLogin(email, password string) error {
	stmt, err := r.db.Prepare("SELECT email, password FROM users WHERE email = $1 AND password = $2")
	if err != nil {
		return err
	}
	defer stmt.Close()

	err = stmt.QueryRow(email, password).Scan(&email, &password)
	if err != nil {
		return err
	}

	return nil
}*/