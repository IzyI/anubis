package usecase

import (
	"anubis/app/DAL/entitiesDB"
	interfacesDB2 "anubis/app/DAL/interfacesDB"
	storage2 "anubis/app/DAL/storage"
	"anubis/app/DTO"
	"anubis/app/core"
	"anubis/app/core/common"
	"anubis/app/core/middlewares"
	"anubis/app/core/schemes"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type ServiceProject struct {
	usr    interfacesDB2.UserRepository
	pr     interfacesDB2.ProjectRepository
	config core.ServiceConfig
}

func NewServiceProject(
	usr *storage2.RepositoryMongoUser,
	pr *storage2.RepositoryMongoProjects,
	config core.ServiceConfig) *ServiceProject {
	return &ServiceProject{
		usr,
		pr,
		config,
	}
}

func (s *ServiceProject) GetProjectsFlow(ctx *gin.Context, input *common.Empty) (DTO.AnswerProjectList, error) {
	var answer DTO.AnswerProjectList
	objectUserID, err := primitive.ObjectIDFromHex(ctx.GetString(middlewares.UserIDKey))
	if err != nil {
		return answer, &schemes.ErrorResponse{Code: 105, Err: "Bad user id", ErrBase: nil}
	}
	projectMap, err := s.pr.GetProjectsListByUser(
		ctx.GetString(middlewares.Service),
		ctx.GetString(middlewares.Domain), objectUserID)
	if err != nil {
		return answer, &schemes.ErrorResponse{Code: 105, Err: "Couldn't create a project", ErrBase: err}
	}
	answer.ListProjects = projectMap
	return answer, nil
}

func (s *ServiceProject) PostProjectsFlow(ctx *gin.Context, input *DTO.CreateProjectValid) (DTO.AnswerProject, error) {
	var answer DTO.AnswerProject

	objectUserID, err := primitive.ObjectIDFromHex(ctx.GetString(middlewares.UserIDKey))
	if err != nil {
		return answer, &schemes.ErrorResponse{Code: 105, Err: "Bad user id", ErrBase: nil}
	}
	var projectMember entitiesDB.MdProjectMember
	var project entitiesDB.MdProject
	projectMember.UserID = objectUserID
	projectMember.Role = "O"
	projectMember.JoinedAt = time.Now()
	project.Domain = ctx.GetString(middlewares.Domain)
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

func (s *ServiceProject) UpdateProjectsFlow(ctx *gin.Context, uri *DTO.UriIDValid, body *DTO.UpdateProjectValid) (common.Empty, error) {
	var answer common.Empty
	objectUserID, err := primitive.ObjectIDFromHex(ctx.GetString(middlewares.UserIDKey))
	if err != nil {
		return answer, &schemes.ErrorResponse{Code: 105, Err: "Bad user id", ErrBase: nil}
	}

	projectID, err := primitive.ObjectIDFromHex(uri.ID)
	if err != nil {
		return answer, &schemes.ErrorResponse{Code: 105, Err: "Bad project id", ErrBase: nil}
	}
	var projectMember entitiesDB.MdProjectMember
	var project entitiesDB.MdProject
	project.ID = projectID
	projectMember.UserID = objectUserID
	project.Domain = ctx.GetString(middlewares.Domain)
	project.Name = body.ProjectName
	project.Members = append(project.Members, projectMember)
	err = s.pr.UpdateProjectName(
		ctx.GetString(middlewares.Service), ctx.GetString(middlewares.Domain), &project)
	if err != nil {
		return answer, &schemes.ErrorResponse{Code: 105, Err: "Couldn't update a project", ErrBase: err}
	}
	return answer, nil
}

func (s *ServiceProject) DelProjectsFlow(ctx *gin.Context, uri *DTO.UriIDValid) (common.Empty, error) {
	var answer common.Empty
	objectUserID, err := primitive.ObjectIDFromHex(ctx.GetString(middlewares.UserIDKey))
	if err != nil {
		return answer, &schemes.ErrorResponse{Code: 105, Err: "Bad user id", ErrBase: nil}
	}

	projectID, err := primitive.ObjectIDFromHex(uri.ID)
	if err != nil {
		return answer, &schemes.ErrorResponse{Code: 105, Err: "Bad project id", ErrBase: nil}
	}
	var projectMember entitiesDB.MdProjectMember
	var project entitiesDB.MdProject
	project.ID = projectID
	projectMember.UserID = objectUserID
	project.Domain = ctx.GetString(middlewares.Domain)
	project.Members = append(project.Members, projectMember)
	err = s.pr.DelProjectID(
		ctx.GetString(middlewares.Service),
		ctx.GetString(middlewares.Domain),
		&project)

	if err != nil {
		return answer, &schemes.ErrorResponse{Code: 105, Err: "Couldn't delete  project", ErrBase: err}
	}
	return answer, nil
}
