package usecase

import (
	"anubis/app/api/DAL/entitiesDB"
	"anubis/app/api/DAL/interfacesDB"
	"anubis/app/api/DAL/storage"
	"anubis/app/api/DTO"
	"anubis/app/core"
	"anubis/app/core/common"
	"anubis/app/core/middlewares"
	"anubis/app/core/schemes"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ServiceProject struct {
	usr    interfacesDB.UserRepository
	pr     interfacesDB.ProjectRepository
	config core.ServiceConfig
}

func NewServiceProject(
	usr *storage.RepositoryMongoUser,
	pr *storage.RepositoryMongoProjects,
	config core.ServiceConfig) *ServiceProject {
	return &ServiceProject{
		usr,
		pr,
		config,
	}
}

func (s *ServiceProject) GetProjectsFlow(ctx *gin.Context, input *common.Empty) (DTO.AnswerProjectList, error) {
	var answer DTO.AnswerProjectList
	userId := ctx.GetString(middlewares.UserIDKey)
	objectID, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return answer, &schemes.ErrorResponse{Code: 107, Err: "Bad user id", ErrBase: nil}
	}
	projectMap, err := s.pr.GetProjectsListByUser(
		ctx.GetString(middlewares.Service),
		ctx.GetString(middlewares.Domain), objectID)
	if err != nil {
		return answer, &schemes.ErrorResponse{Code: 105, Err: "Couldn't create a project", ErrBase: err}
	}
	answer.ListProjects = projectMap
	return answer, nil
}

func (s *ServiceProject) PostProjectsFlow(ctx *gin.Context, input *DTO.CreateProjectValid) (DTO.AnswerProject, error) {
	var answer DTO.AnswerProject

	objectID, err := primitive.ObjectIDFromHex(ctx.GetString(middlewares.UserIDKey))
	if err != nil {
		return answer, &schemes.ErrorResponse{Code: 107, Err: "Bad user id", ErrBase: nil}
	}
	var projectMember entitiesDB.MdProjectMember
	var project entitiesDB.MdProject
	projectMember.UserID = objectID
	projectMember.Role = "O"
	project.Domain = middlewares.Domain
	project.Name = input.ProjectName
	project.Members = append(project.Members, projectMember)

	err = s.pr.CreateProject(
		ctx.GetString(middlewares.Service), &project)
	if err != nil {
		return answer, &schemes.ErrorResponse{Code: 105, Err: "Couldn't create a project", ErrBase: err}
	}
	answer.ProjectName = input.ProjectName
	answer.ProjectId = project.ID.Hex()
	return answer, nil
}

func (s *ServiceProject) PutProjectsFlow(ctx *gin.Context, input *DTO.UpdateProjectValid) (common.Empty, error) {
	var answer common.Empty
	objectID, err := primitive.ObjectIDFromHex(ctx.GetString(middlewares.UserIDKey))
	if err != nil {
		return answer, &schemes.ErrorResponse{Code: 107, Err: "Bad user id", ErrBase: nil}
	}

	projectID, err := primitive.ObjectIDFromHex(input.ID)
	if err != nil {
		return answer, &schemes.ErrorResponse{Code: 107, Err: "Bad project id", ErrBase: nil}
	}
	var projectMember entitiesDB.MdProjectMember
	var project entitiesDB.MdProject
	project.ID = projectID
	projectMember.UserID = objectID
	project.Domain = middlewares.Domain
	project.Name = input.ProjectName
	project.Members = append(project.Members, projectMember)
	err = s.pr.UpdateProjectName(
		ctx.GetString(middlewares.Service), &project)
	if err != nil {
		return answer, &schemes.ErrorResponse{Code: 105, Err: "Couldn't update a project", ErrBase: err}
	}
	return answer, nil
}

func (s *ServiceProject) DelProjectsFlow(ctx *gin.Context, input *DTO.DelProjectValid) (common.Empty, error) {
	var answer common.Empty
	objectID, err := primitive.ObjectIDFromHex(ctx.GetString(middlewares.UserIDKey))
	if err != nil {
		return answer, &schemes.ErrorResponse{Code: 107, Err: "Bad user id", ErrBase: nil}
	}

	projectID, err := primitive.ObjectIDFromHex(input.ID)
	if err != nil {
		return answer, &schemes.ErrorResponse{Code: 107, Err: "Bad project id", ErrBase: nil}
	}
	var projectMember entitiesDB.MdProjectMember
	var project entitiesDB.MdProject
	project.ID = projectID
	projectMember.UserID = objectID
	project.Domain = middlewares.Domain
	err = s.pr.DelProjectName(
		ctx.GetString(middlewares.Service), &project)
	if err != nil {
		return answer, &schemes.ErrorResponse{Code: 105, Err: "Couldn't delete  project", ErrBase: err}
	}
	return answer, nil
}
