package storage

import (
	"anubis/app/api/entytes"
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RepositoryPsqlUser struct {
	db *pgxpool.Pool
}

func NewRepositoryPsqlUser(db *pgxpool.Pool) *RepositoryPsqlUser {
	return &RepositoryPsqlUser{db: db}
}

func (r *RepositoryPsqlUser) CreateUser(u entytes.MdUser) (*entytes.MdUser, error) {
	var uid uuid.UUID
	sql := `INSERT INTO users (phone,password_hash)
          VALUES ($1,$2) RETURNING (uuid)`
	rows := r.db.QueryRow(context.Background(), sql,
		u.Phone, u.PasswordHash)

	err := rows.Scan(&uid)
	if err != nil {
		return &u, err
	}
	u.Uuid = uid.String()

	return &u, nil
}

func (r *RepositoryPsqlUser) SmsSaveUser(userUuid string, s string) error {
	var uid uuid.UUID

	sql := `UPDATE users  SET sms=($1)
          WHERE uuid=($2) RETURNING uuid`
	rows := r.db.QueryRow(context.Background(), sql,
		s, userUuid)

	err := rows.Scan(&uid)
	if err != nil {
		return fmt.Errorf("failed to add sms for user")
	}
	return nil
}

func (r *RepositoryPsqlUser) SmsValidUser(userUuid string, s string) error {
	var uid uuid.UUID

	sql := `UPDATE users  SET verification=true
          WHERE uuid=($2) AND sms=($1) and verification=False RETURNING uuid`
	rows := r.db.QueryRow(context.Background(), sql,
		s, userUuid)

	err := rows.Scan(&uid)
	if err != nil {
		return err
	}
	return nil
}

func (r *RepositoryPsqlUser) LoginUser(phone string) (string, string, error) {
	var uid uuid.UUID
	var passwordHash string
	sql := `SELECT uuid ,password_hash FROM users   WHERE phone=($1) and verification=true`
	rows := r.db.QueryRow(context.Background(), sql, phone)

	err := rows.Scan(&uid, &passwordHash)
	if err != nil {
		return "", "", err
	}

	return uid.String(), passwordHash, nil
}

func (r *RepositoryPsqlUser) GetUuidUser(uuid string) error {
	var p string
	sql := `SELECT phone FROM users   WHERE uuid=($1) and verification=true`
	rows := r.db.QueryRow(context.Background(), sql, uuid)

	err := rows.Scan(&p)
	if err != nil {
		return err
	}

	return nil
}
