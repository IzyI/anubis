package interfacesDB

import (
	"anubis/app/DAL/entitiesDB"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProjectRepository interface {
	CreateProject(service string, project *entitiesDB.MdProject) error
	UpdateProjectName(service string, domain string, project *entitiesDB.MdProject) error
	DelProjectID(service string, domain string, project *entitiesDB.MdProject) error
	AddMemberToProject(
		service string,
		domain string,
		projectID primitive.ObjectID,
		ownerID primitive.ObjectID,
		userID primitive.ObjectID,
		role string,
	) error
	UpdateMemberRole(
		service string,
		domain string,
		projectID primitive.ObjectID,
		ownerID primitive.ObjectID,
		userID primitive.ObjectID,
		newRole string,
	) error

	RemoveMemberFromProject(
		service string,
		domain string,
		projectID primitive.ObjectID,
		ownerID primitive.ObjectID,
		userID primitive.ObjectID,
	) error
	GetProjectsListByUser(service string, domain string, userID primitive.ObjectID) (map[string]string, error)
	GetProjectsByUser(service string, project *entitiesDB.MdProject, userID primitive.ObjectID) error
	// MEMBERS
	GetProjectByUserOwnerID(service string, domain string, project *entitiesDB.MdProject, userId primitive.ObjectID) error
	GetProjectIDByUserID(service string, project *entitiesDB.MdProject, userId primitive.ObjectID) error
}
