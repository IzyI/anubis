package interfacesDB

import (
	"anubis/app/api/DAL/entitiesDB"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserRepository interface {
	CreateUser(domain string, user *entitiesDB.MdUser) error
	CheckOldAndUpdateSession(service string, userSession *entitiesDB.MdUsersSession) error
	DeactivateOldAndCreateSession(service string, userSession *entitiesDB.MdUsersSession) error
	DeactivateUserSessionsByDomain(service string, userSession *entitiesDB.MdUsersSession) error
	DeactivateUserSessionsByTokenFamily(service string, idTokenFamily primitive.ObjectID) error
	GetUsersSessionByID(service string, id primitive.ObjectID, userSession *entitiesDB.MdUsersSession) error
	UpdateSessionsByID(service string, userSession *entitiesDB.MdUsersSession) error
	InsertUsersSession(service string, userSession *entitiesDB.MdUsersSession) error
	DeleteSessionsByID(service string, userSession *entitiesDB.MdUsersSession) error
	UserSessionsSetActive(service string, userSession *entitiesDB.MdUsersSession, active bool) error
}
