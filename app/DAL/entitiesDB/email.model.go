package entitiesDB

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MdEmailAuth struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	Email        string             `bson:"email"`
	Domain       string             `bson:"domain"`
	PasswordHash string             `bson:"password_hash"`
	Verification bool               `bson:"verification"`
	UserID       primitive.ObjectID `bson:"user_id"`
}

type MdEmailCodeAuth struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	UserID       primitive.ObjectID `bson:"user_id"`
	Email        string             `bson:"email"`
	EmailCode    string             `bson:"email_code"`
	EmailService string             `bson:"email_service"`
	IDSend       string             `bson:"id_send"`
}
