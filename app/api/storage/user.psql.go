package storage

import (
	"anubis/app/api/entytes"
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RepositoryPsqlUser struct {
	db *pgxpool.Pool
}

func NewRepositoryPsqlUser(db *pgxpool.Pool) *RepositoryPsqlUser {
	return &RepositoryPsqlUser{db: db}
}
func (r *RepositoryPsqlUser) CreateUser() (*entytes.MdUser, error) {
	var user entytes.MdUser
	sql := `INSERT INTO users DEFAULT VALUES RETURNING uuid`
	rows := r.db.QueryRow(context.Background(), sql)

	err := rows.Scan(&user.Uuid)
	if err != nil {
		return &user, err
	}

	return &user, nil
}
