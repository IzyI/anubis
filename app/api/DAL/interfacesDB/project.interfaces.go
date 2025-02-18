package interfacesDB

import (
	"anubis/app/api/DAL/entitiesDB"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProjectRepository interface {
	CreateProject(service string, project *entitiesDB.MdProject) (primitive.ObjectID, error)
	AddMemberToProject(service string, projectID primitive.ObjectID, userID primitive.ObjectID, role string) error
	GetProjectsListByUser(service string, domain string, userID primitive.ObjectID) (map[string]string, error)
	GetProjectsByUser(service string, project *entitiesDB.MdProject, userID primitive.ObjectID) error
}
