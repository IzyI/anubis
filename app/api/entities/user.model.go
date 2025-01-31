package entities

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"net/netip"
)

type MdUser struct {
	Uuid       uuid.UUID
	GlobalRole pgtype.Int4
}

type MdUsersSession struct {
	Ip           netip.Addr
	Uuid         pgtype.UUID
	IDDevice     string
	Type         string
	Revoked      bool
	RefreshToken string
	UserUuid     pgtype.UUID
	IDService    int32
}

type MdUsersInfo struct {
	UserUuid pgtype.UUID
	Email    string
	Nickname pgtype.Text
}

type MdUsersGroup struct {
	UserUuid  pgtype.UUID
	GroupName string
	Service   pgtype.Text
}
