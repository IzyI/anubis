package entitiesDB

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type MdUser struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Nickname string             `bson:"nickname,omitempty"`
}

type MdUsersSession struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	UserID     primitive.ObjectID `bson:"user_id"`
	Domain     string             `bson:"domain"`
	DeviceId   string             `bson:"device_id"`
	HashToken  string             `bson:"hash_token"`
	DeviceType string             `bson:"device_type"` // Идентификатор сервиса
	CreatedAt  time.Time          `bson:"created_at"`
	ExpiresAt  time.Time          `bson:"expires_at"`
	IP         string             `bson:"ip"`
	IsActive   bool               `bson:"is_active"`
}
