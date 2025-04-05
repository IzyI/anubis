package interfacesDB

import (
	"anubis/app/DAL/entitiesDB"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserRepository interface {
	CreateUser(service string, user *entitiesDB.MdUser) error
	GetUserByID(service string, user *entitiesDB.MdUser) error
	GetUsersByIDs(service string, userIDs []primitive.ObjectID) (map[primitive.ObjectID]entitiesDB.MdUser, error)
	CheckOldAndUpdateSession(service string, userSession *entitiesDB.MdUsersSession) error
	DeactivateOldAndCreateSession(service string, userSession *entitiesDB.MdUsersSession) error
	DeactivateUserSessionsByDomain(
		service string,
		userID primitive.ObjectID,
		domain string,
	) error
	DeactivateUserSessionsByTokenFamily(service string, idTokenFamily primitive.ObjectID) error
	GetUsersSessionByID(service string, id primitive.ObjectID, userSession *entitiesDB.MdUsersSession) error
	UpdateSessionsByID(service string, userSession *entitiesDB.MdUsersSession) error
	InsertUsersSession(service string, userSession *entitiesDB.MdUsersSession) error
	DeleteSessionsByID(service string, userSession *entitiesDB.MdUsersSession) error
	UserSessionsSetActive(service string, userSessionID primitive.ObjectID, active bool) error
}
