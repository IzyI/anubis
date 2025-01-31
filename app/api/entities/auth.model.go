package entities

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type MdPhoneAuth struct {
	Phone        int64
	CountryCode  int32
	PasswordHash string
	CreatedAt    pgtype.Timestamptz
	Verification bool
	UserUuid     uuid.UUID
}

type MdSmsAuth struct {
	SmsUuid    uuid.UUID
	UserUuid   uuid.UUID
	Phone      int64
	SmsCode    string
	SmsService string
	IDSend     string
	CreatedAt  pgtype.Timestamptz
}
