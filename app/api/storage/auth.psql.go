package storage

import (
	"anubis/app/api/entytes"
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RepositoryPsqlAuth struct {
	db *pgxpool.Pool
}

func NewRepositoryPsqlAuth(db *pgxpool.Pool) *RepositoryPsqlAuth {
	return &RepositoryPsqlAuth{db: db}
}

func (r *RepositoryPsqlAuth) SavePhone(phone *entytes.MdPhoneAuth) error {
	sql := `INSERT INTO phone_auth (phone,country_code,password_hash,verification,user_uuid)
          VALUES ($1,$2,$3,$4,$5)`
	rows := r.db.QueryRow(context.Background(), sql,
		phone.Phone, phone.CountryCode, phone.PasswordHash, phone.Verification, phone.UserUuid)

	err := rows.Scan()
	if err != nil {
		return err
	}
	return nil
}

func (r *RepositoryPsqlAuth) SmsSaveUser(sms entytes.SmsAuth) error {
	sql := `INSERT INTO sms_auth (user_uuid,phone,sms_code,sms_service,id_send)
          VALUES ($1,$2,$3,$4,$5)`
	rows := r.db.QueryRow(context.Background(), sql,
		sms.UserUuid, sms.Phone, sms.SmsCode, sms.SmsService, sms.IdSend)

	err := rows.Scan()
	if err != nil {
		return fmt.Errorf("failed to add sms for user")
	}
	return nil
}

func (r *RepositoryPsqlAuth) SmsValidUser(userUuid string, s string) error {
	//var uid uuid.UUID
	//
	//sql := `UPDATE users  SET verification=true
	//      WHERE uuid=($2) AND sms=($1) and verification=False RETURNING uuid`
	//rows := r.db.QueryRow(context.Background(), sql,
	//	s, userUuid)
	//
	//err := rows.Scan(&uid)
	//if err != nil {
	//	return err
	//}
	return nil
}

func (r *RepositoryPsqlAuth) GetUserPhone(phone int) (*entytes.MdPhoneAuth, error) {
	var phoneUser entytes.MdPhoneAuth
	sql := `SELECT * FROM phone_auth   WHERE phone=($1)`
	rows := r.db.QueryRow(context.Background(), sql, phone)
	err := rows.Scan(

		&phoneUser.Phone,
		&phoneUser.CountryCode,
		&phoneUser.PasswordHash,
		&phoneUser.UpdatedAt,
		&phoneUser.CreatedAt,
		&phoneUser.Verification,
		&phoneUser.UserUuid,
	)

	if err != nil {
		return &phoneUser, err
	}
	return &phoneUser, nil
}
func (r *RepositoryPsqlAuth) LoginUser(phone string) (string, string, error) {
	var uid uuid.UUID
	//var passwordHash string
	//sql := `SELECT uuid ,password_hash FROM users   WHERE phone=($1) and verification=true`
	//rows := r.db.QueryRow(context.Background(), sql, phone)
	//
	//err := rows.Scan(&uid, &passwordHash)
	//if err != nil {
	//	return "", "", err
	//}

	return uid.String(), uid.String(), nil
}

func (r *RepositoryPsqlAuth) GetUuidUser(uuid string) error {
	//var p string
	//sql := `SELECT phone FROM users   WHERE uuid=($1) and verification=true`
	//rows := r.db.QueryRow(context.Background(), sql, uuid)
	//
	//err := rows.Scan(&p)
	//if err != nil {
	//	return err
	//}

	return nil
}
