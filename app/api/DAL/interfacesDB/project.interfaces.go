package interfacesDB

import (
	"anubis/app/api/DAL/entitiesDB"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProjectRepository interface {
	CreateProject(service string, project *entitiesDB.MdProject) (primitive.ObjectID, error)
	AddMemberToProject(service string, projectID primitive.ObjectID, userID string, role string) error
	GetProjectsByUser(service string, userID string) ([]entitiesDB.MdProject, error)
}
