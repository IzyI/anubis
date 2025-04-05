package entitiesDB

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MdPhoneAuth struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	Phone        int64              `bson:"phone"`
	CountryCode  int32              `bson:"country_code"`
	PasswordHash string             `bson:"password_hash"`
	Verification bool               `bson:"verification"`
	UserID       primitive.ObjectID `bson:"user_id"`
}

type MdSmsAuth struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	UserID     primitive.ObjectID `bson:"user_id"`
	Phone      int64              `bson:"phone"`
	SmsCode    string             `bson:"sms_code"`
	SmsService string             `bson:"sms_service"`
	IDSend     string             `bson:"id_send"`
}
