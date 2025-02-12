package entitiesDB

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MdUser struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Nickname string             `bson:"nickname,omitempty"`
}

//type MdUsersSession struct {
//	Ip           netip.Addr
//	Uuid         pgtype.UUID
//	IDDevice     string
//	Type         string
//	Revoked      bool
//	RefreshToken string
//	UserUuid     pgtype.UUID
//	IDService    int32
//}

type MdUsersGroup struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	UserID    primitive.ObjectID `bson:"user_id"`
	GroupName string             `bson:"group_name"`
	Owner     bool               `bson:"owner"`
	Role      string             `bson:"role"`
	Domain    string             `bson:"domain"`
}
