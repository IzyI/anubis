package interfacesDB

import (
	"anubis/app/DAL/entitiesDB"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AuthEmailRepository interface {
	GetEmail(service string, email *entitiesDB.MdEmailAuth) error
	SaveEmail(service string, email *entitiesDB.MdEmailAuth) error
	SaveEmailAuth(service string, emailCode *entitiesDB.MdEmailCodeAuth) error
	EmailCodeValidUser(service string, ID primitive.ObjectID, emailCode string) (string, error)
	SaveVerifyEmail(service string, verify bool, email string) error
	GetEmailVerificationUserID(service string, email string, verification bool) (primitive.ObjectID, string, error)
}
