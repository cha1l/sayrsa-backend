package repository

import (
	"fmt"
	"time"

	"github.com/cha1l/sayrsa-2.0/models"
	"github.com/jmoiron/sqlx"
)

type AuthRepo struct {
	db *sqlx.DB
}

func NewAuthRepo(db *sqlx.DB) *AuthRepo {
	return &AuthRepo{
		db: db,
	}
}

func (s *AuthRepo) CreateUser(u models.User, token string, tokenT time.Time) error {
	var id int

	tx, err := s.db.Beginx()
	if err != nil {
		return err
	}

	createUserQuery := fmt.Sprintf(`INSERT INTO %s (username, bio, last_online, password_hash, public_key) 
	VALUES ($1, $2, $3, $4, $5) RETURNING id`, usersTable)
	row := tx.QueryRow(createUserQuery, u.Username, u.Bio, u.LastOnline, u.Password, u.PublicKey)
	if err := row.Scan(&id); err != nil {
		tx.Rollback()
		return err
	}

	createTokenQuery := fmt.Sprintf(`INSERT INTO %s (user_id, token, expires_at) VALUES ($1, $2, $3)`, tokensTable)
	_, err = tx.Exec(createTokenQuery, id, token, tokenT)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
