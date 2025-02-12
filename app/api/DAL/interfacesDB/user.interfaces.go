package interfacesDB

import (
	"anubis/app/api/DAL/entitiesDB"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserRepository interface {
	CreateUser(domain string, user *entitiesDB.MdUser) error
	GetGroupUser(service string, domain string, userID primitive.ObjectID) (map[string]string, error)
}
