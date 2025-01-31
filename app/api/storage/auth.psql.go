package storage

import (
	"anubis/app/api/entities"
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RepositoryPsqlAuthPhone struct {
	db *pgxpool.Pool
}

func NewRepositoryPsqlAuthPhone(db *pgxpool.Pool) *RepositoryPsqlAuthPhone {
	return &RepositoryPsqlAuthPhone{db: db}
}

// GetUserPhone---------------------------------------------------------------------------------------------------
const getQueryUserPhone = `-- name: getQueryUserPhone :one
SELECT  phone,
        country_code,
        password_hash,
        created_at,
        verification,
        user_uuid  FROM phone_auth   WHERE phone=($1)
`

func (r RepositoryPsqlAuthPhone) GetUserPhone(phone int64) (*entities.MdPhoneAuth, error) {
	row := r.db.QueryRow(context.Background(), getQueryUserPhone, phone)
	var i entities.MdPhoneAuth
	err := row.Scan(
		&i.Phone,
		&i.CountryCode,
		&i.PasswordHash,
		&i.CreatedAt,
		&i.Verification,
		&i.UserUuid,
	)
	return &i, err
}

// ---------------------------------------------------------------------------------------------------

// SavePhone---------------------------------------------------------------------------------------------------
const saveQueryPhone = `-- name: saveQueryPhone :one
INSERT INTO phone_auth (phone, country_code, password_hash, verification, user_uuid)
VALUES ($1, $2, $3, $4, $5)
ON CONFLICT (phone)
    DO UPDATE SET phone = EXCLUDED.phone ,user_uuid =  $5 
RETURNING verification
`

func (r *RepositoryPsqlAuthPhone) SavePhone(phone *entities.MdPhoneAuth) (bool, error) {
	row := r.db.QueryRow(context.Background(), saveQueryPhone,
		phone.Phone,
		phone.CountryCode,
		phone.PasswordHash,
		phone.Verification,
		phone.UserUuid,
	)
	var verification bool
	err := row.Scan(&verification)
	return verification, err
}

// ---------------------------------------------------------------------------------------------------

// SavePhone---------------------------------------------------------------------------------------------------
const SaveQueryVerifyPhone = `-- name: SaveQueryVerifyPhone :one
UPDATE phone_auth  SET verification = $1  WHERE  phone = $2
`

func (r *RepositoryPsqlAuthPhone) SaveVerifyPhone(verify bool, phone int64) error {
	row := r.db.QueryRow(context.Background(), SaveQueryVerifyPhone,
		verify, phone,
	)
	err := row.Scan()
	return err
}

// ---------------------------------------------------------------------------------------------------

// SaveSmsAuth---------------------------------------------------------------------------------------------------

const saveQuerySmsAuth = `-- name: saveQuerySmsAuth :one
INSERT INTO sms_auth (user_uuid, phone, sms_code, sms_service, id_send)
VALUES ($1, $2, $3, $4, $5)
RETURNING sms_uuid
`

func (r *RepositoryPsqlAuthPhone) SaveSmsAuth(sms *entities.MdSmsAuth) (uuid.UUID, error) {
	row := r.db.QueryRow(context.Background(), saveQuerySmsAuth,
		sms.UserUuid,
		sms.Phone,
		sms.SmsCode,
		sms.SmsService,
		sms.IDSend,
	)
	var smsUuid uuid.UUID
	err := row.Scan(&smsUuid)
	return smsUuid, err
}

// ---------------------------------------------------------------------------------------------------

// SmsValidUser---------------------------------------------------------------------------------------------------

const getQuerySmsValidUser = `-- name: getQuerySmsValidUser :one
SELECT phone FROM sms_auth WHERE sms_code = $1 and sms_uuid = $2 AND created_at >= NOW() - INTERVAL '10 minutes'
`

func (r *RepositoryPsqlAuthPhone) SmsValidUser(userUuid uuid.UUID, smsCode string) (int64, error) {

	row := r.db.QueryRow(context.Background(), getQuerySmsValidUser, smsCode, userUuid)
	var phone int64
	err := row.Scan(&phone)
	return phone, err
}

// ---------------------------------------------------------------------------------------------------

func (r *RepositoryPsqlAuthPhone) LoginUser(phone string) (uuid.UUID, uuid.UUID, error) {
	var uid uuid.UUID
	//var passwordHash string
	//sql := `SELECT uuid ,password_hash FROM users   WHERE phone=($1) and verification=true`
	//rows := r.db.QueryRow(context.Background(), sql, phone)
	//
	//err := rows.Scan(&uid, &passwordHash)
	//if err != nil {
	//	return "", "", err
	//}

	return uid, uid, nil
}

func (r *RepositoryPsqlAuthPhone) GetUuidUser(uuid uuid.UUID) error {
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
