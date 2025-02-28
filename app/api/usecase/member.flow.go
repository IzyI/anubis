package usecase

import (
	"anubis/app/api/DAL/entitiesDB"
	"anubis/app/api/DAL/interfacesDB"
	"anubis/app/api/DAL/storage"
	"anubis/app/api/DTO"
	"anubis/app/core"
	"anubis/app/core/middlewares"
	"anubis/app/core/schemes"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type ServiceProjectMembers struct {
	usr    interfacesDB.UserRepository
	pr     interfacesDB.ProjectRepository
	config core.ServiceConfig
}

func NewServiceServiceProjectMembers(
	usr *storage.RepositoryMongoUser,
	pr *storage.RepositoryMongoProjects,
	config core.ServiceConfig) *ServiceProjectMembers {
	return &ServiceProjectMembers{
		usr,
		pr,
		config,
	}
}

func (s *ServiceProjectMembers) GetProjectMembersFlow(ctx *gin.Context, input *DTO.MembersProjectIdValid) (DTO.AnswerProjectID, error) {
	var answer DTO.AnswerProjectID
	var memberAnswer []DTO.AnswerMembers

	userId := ctx.GetString(middlewares.UserIDKey)
	objectID, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return answer, &schemes.ErrorResponse{Code: 107, Err: "Bad user id", ErrBase: nil}
	}
	projectID, err := primitive.ObjectIDFromHex(input.ID)
	if err != nil {
		return answer, &schemes.ErrorResponse{Code: 107, Err: "Bad project id", ErrBase: nil}
	}
	var project entitiesDB.MdProject
	project.ID = projectID
	err = s.pr.GetProjectIDByUserID(
		ctx.GetString(middlewares.Service), &project, objectID)
	if err != nil {
		return answer, &schemes.ErrorResponse{Code: 105, Err: "Couldn't find  project for owner", ErrBase: err}
	}

	for _, member := range project.Members {
		// Создаём новый объект AnswerMembers, заполняя его данными из MdProjectMember
		answerMember := DTO.AnswerMembers{
			UserID:   member.UserID.Hex(), // Преобразуем ObjectID в строку
			Role:     member.Role,
			JoinedAt: member.JoinedAt.Format(time.RFC3339), // Форматируем дату как строку
		}

		// Добавляем объект в срез memberAnswer
		memberAnswer = append(memberAnswer, answerMember)
	}

	answer.Domain = project.Domain
	answer.Name = project.Name
	answer.ID = project.ID.Hex()
	answer.Members = memberAnswer
	return answer, nil
}

func (s *ServiceProjectMembers) PostProjectMembersFlow(ctx *gin.Context, input *DTO.MembersAddIdValid) (DTO.AnswerMembers, error) {
	var answer DTO.AnswerMembers

	ownerID, err := primitive.ObjectIDFromHex(ctx.GetString(middlewares.UserIDKey))
	if err != nil {
		return answer, &schemes.ErrorResponse{Code: 107, Err: "Bad owner id", ErrBase: nil}
	}
	userID, err := primitive.ObjectIDFromHex(input.UserID)
	if err != nil {
		return answer, &schemes.ErrorResponse{Code: 107, Err: "Bad user id", ErrBase: nil}
	}
	projectID, err := primitive.ObjectIDFromHex(input.ID)
	if err != nil {
		return answer, &schemes.ErrorResponse{Code: 107, Err: "Bad project id", ErrBase: nil}
	}
	if input.Role == "O" {
		return answer, &schemes.ErrorResponse{Code: 107, Err: "Bad role", ErrBase: nil}
	}
	if input.Role == "" {
		input.Role = "E"
	}
	err = s.pr.AddMemberToProject(
		ctx.GetString(middlewares.Service), projectID, ownerID, userID, input.Role)
	if err != nil {
		return answer, &schemes.ErrorResponse{Code: 105, Err: err.Error(), ErrBase: err}
	}

	return answer, nil
}

func (s *ServiceProjectMembers) PutProjectMembersFlow(ctx *gin.Context, input *DTO.MembersAddIdValid) (DTO.AnswerMembers, error) {
	var answer DTO.AnswerMembers

	ownerID, err := primitive.ObjectIDFromHex(ctx.GetString(middlewares.UserIDKey))
	if err != nil {
		return answer, &schemes.ErrorResponse{Code: 107, Err: "Bad owner id", ErrBase: nil}
	}
	userID, err := primitive.ObjectIDFromHex(input.UserID)
	if err != nil {
		return answer, &schemes.ErrorResponse{Code: 107, Err: "Bad user id", ErrBase: nil}
	}
	projectID, err := primitive.ObjectIDFromHex(input.ID)
	if err != nil {
		return answer, &schemes.ErrorResponse{Code: 107, Err: "Bad project id", ErrBase: nil}
	}
	if ownerID == userID {
		return answer, &schemes.ErrorResponse{Code: 107, Err: "owner == user", ErrBase: nil}
	}
	if input.Role == "O" {
		return answer, &schemes.ErrorResponse{Code: 107, Err: "Bad role", ErrBase: nil}
	}
	if input.Role == "" {
		input.Role = "E"
	}
	err = s.pr.UpdateMemberRole(
		ctx.GetString(middlewares.Service), projectID, ownerID, userID, input.Role)
	if err != nil {
		return answer, &schemes.ErrorResponse{Code: 105, Err: err.Error(), ErrBase: err}
	}

	return answer, nil
}

func (s *ServiceProjectMembers) DelProjectMembersFlow(ctx *gin.Context, input *DTO.MembersAddIdValid) (DTO.AnswerMembers, error) {
	var answer DTO.AnswerMembers

	ownerID, err := primitive.ObjectIDFromHex(ctx.GetString(middlewares.UserIDKey))
	if err != nil {
		return answer, &schemes.ErrorResponse{Code: 107, Err: "Bad owner id", ErrBase: nil}
	}
	userID, err := primitive.ObjectIDFromHex(input.UserID)
	if err != nil {
		return answer, &schemes.ErrorResponse{Code: 107, Err: "Bad user id", ErrBase: nil}
	}
	projectID, err := primitive.ObjectIDFromHex(input.ID)
	if err != nil {
		return answer, &schemes.ErrorResponse{Code: 107, Err: "Bad project id", ErrBase: nil}
	}
	if ownerID == userID {
		return answer, &schemes.ErrorResponse{Code: 107, Err: "owner == user", ErrBase: nil}
	}

	err = s.pr.RemoveMemberFromProject(
		ctx.GetString(middlewares.Service), projectID, ownerID, userID)
	if err != nil {
		return answer, &schemes.ErrorResponse{Code: 105, Err: err.Error(), ErrBase: err}
	}

	return answer, nil
}
