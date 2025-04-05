package interfacesDB

import (
	"anubis/app/DAL/entitiesDB"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AuthPhoneRepository interface {
	GetPhone(service string, phone *entitiesDB.MdPhoneAuth) error
	SavePhone(service string, phone *entitiesDB.MdPhoneAuth) error
	SaveSmsAuth(service string, sms *entitiesDB.MdSmsAuth) error
	SmsValidUser(service string, ID primitive.ObjectID, smsCode string) (int64, error)
	SaveVerifyPhone(service string, verify bool, phone int64) error
	GetPhoneVerificationUserID(service string, phone int64, verification bool) (primitive.ObjectID, string, error)
}
