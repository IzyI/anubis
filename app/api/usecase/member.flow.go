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
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ServiceProjectMembers struct {
	usr    interfacesDB2.UserRepository
	pr     interfacesDB2.ProjectRepository
	config core.ServiceConfig
}

func NewServiceServiceProjectMembers(
	usr *storage2.RepositoryMongoUser,
	pr *storage2.RepositoryMongoProjects,
	config core.ServiceConfig) *ServiceProjectMembers {
	return &ServiceProjectMembers{
		usr,
		pr,
		config,
	}
}

func (s *ServiceProjectMembers) GetProjectMembersFlow(ctx *gin.Context, input *DTO.UriIDValid) (DTO.AnswerProjectID, error) {
	var answer DTO.AnswerProjectID

	userId := ctx.GetString(middlewares.UserIDKey)
	objectID, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return answer, &schemes.ErrorResponse{Code: 105, Err: "Bad user id", ErrBase: nil}
	}
	projectID, err := primitive.ObjectIDFromHex(input.ID)
	if err != nil {
		return answer, &schemes.ErrorResponse{Code: 105, Err: "Bad project id", ErrBase: nil}
	}
	var project entitiesDB.MdProject
	project.ID = projectID

	err = s.pr.GetProjectIDByUserID(
		ctx.GetString(middlewares.Service), &project, objectID)
	if err != nil {
		return answer, &schemes.ErrorResponse{Code: 105, Err: "Couldn't find  project for owner", ErrBase: err}
	}

	userIDs := make([]primitive.ObjectID, 0, len(project.Members))
	for _, member := range project.Members {
		userIDs = append(userIDs, member.UserID)
	}

	usersMap, err := s.usr.GetUsersByIDs(
		ctx.GetString(middlewares.Service),
		userIDs,
	)
	if err != nil {
		return answer, &schemes.ErrorResponse{Code: 105, Err: "Error fetching members", ErrBase: err}
	}

	memberAnswer := make([]DTO.AnswerMembers, 0, len(project.Members))
	for _, member := range project.Members {
		fmt.Printf("member.UserID  %v \n", member.UserID)
		user, exists := usersMap[member.UserID]
		if !exists {
			return answer, &schemes.ErrorResponse{Code: 104, Err: fmt.Sprintf("Member %s not found", member.UserID.Hex())}
		}

		memberAnswer = append(memberAnswer, DTO.AnswerMembers{
			UserID:   member.UserID.Hex(),
			Role:     member.Role,
			Nickname: user.Nickname,
		})
	}

	answer.Domain = project.Domain
	answer.Name = project.Name
	answer.ID = project.ID.Hex()
	answer.Members = memberAnswer
	return answer, nil
}

func (s *ServiceProjectMembers) PostProjectMembersFlow(
	ctx *gin.Context, uri *DTO.UriIDValid, body *DTO.MembersAddIdValid) (DTO.AnswerMembers, error) {
	var answer DTO.AnswerMembers

	ownerID, err := primitive.ObjectIDFromHex(ctx.GetString(middlewares.UserIDKey))
	if err != nil {
		return answer, &schemes.ErrorResponse{Code: 105, Err: "Bad owner id", ErrBase: nil}
	}
	userID, err := primitive.ObjectIDFromHex(body.UserID)
	if err != nil {
		return answer, &schemes.ErrorResponse{Code: 105, Err: "Bad user id", ErrBase: nil}
	}
	projectID, err := primitive.ObjectIDFromHex(uri.ID)
	if err != nil {
		return answer, &schemes.ErrorResponse{Code: 105, Err: "Bad project id", ErrBase: nil}
	}
	if body.Role == s.config.ListServices[ctx.GetString(middlewares.Domain)].Role["owner"] {
		return answer, &schemes.ErrorResponse{Code: 107, Err: "Bad role", ErrBase: nil}
	}
	if body.Role == "" {
		body.Role = s.config.ListServices[ctx.GetString(middlewares.Domain)].Role["user"]
	}

	var user entitiesDB.MdUser
	user.ID = userID
	err = s.usr.GetUserByID(
		ctx.GetString(middlewares.Service), &user)
	if err != nil {
		return answer, &schemes.ErrorResponse{Code: 104, Err: "Not found user", ErrBase: err}
	}
	err = s.pr.AddMemberToProject(
		ctx.GetString(middlewares.Service), ctx.GetString(middlewares.Domain), projectID, ownerID, userID, body.Role)
	if err != nil {
		return answer, &schemes.ErrorResponse{Code: 105, Err: err.Error(), ErrBase: err}
	}
	answer.Role = body.Role
	answer.UserID = body.UserID
	answer.Nickname = user.Nickname
	return answer, nil
}

func (s *ServiceProjectMembers) PutProjectMembersFlow(ctx *gin.Context, uri *DTO.UriIDValid, body *DTO.MembersAddIdValid) (DTO.AnswerMembers, error) {
	var answer DTO.AnswerMembers

	ownerID, err := primitive.ObjectIDFromHex(ctx.GetString(middlewares.UserIDKey))
	if err != nil {
		return answer, &schemes.ErrorResponse{Code: 105, Err: "Bad owner id", ErrBase: nil}
	}
	userID, err := primitive.ObjectIDFromHex(body.UserID)
	if err != nil {
		return answer, &schemes.ErrorResponse{Code: 105, Err: "Bad user id", ErrBase: nil}
	}
	projectID, err := primitive.ObjectIDFromHex(uri.ID)
	if err != nil {
		return answer, &schemes.ErrorResponse{Code: 107, Err: "Bad project id", ErrBase: nil}
	}
	if ownerID == userID {
		return answer, &schemes.ErrorResponse{Code: 105, Err: "owner == user", ErrBase: nil}
	}
	if body.Role == s.config.ListServices[ctx.GetString(middlewares.Domain)].Role["owner"] {
		return answer, &schemes.ErrorResponse{Code: 105, Err: "Bad role", ErrBase: nil}
	}
	if body.Role == "" {
		body.Role = s.config.ListServices[ctx.GetString(middlewares.Domain)].Role["user"]
	}

	var user entitiesDB.MdUser
	user.ID = userID
	err = s.usr.GetUserByID(
		ctx.GetString(middlewares.Service), &user)
	if err != nil {
		return answer, &schemes.ErrorResponse{Code: 104, Err: "Not found user", ErrBase: err}
	}

	err = s.pr.UpdateMemberRole(
		ctx.GetString(middlewares.Service),
		ctx.GetString(middlewares.Domain),
		projectID,
		ownerID,
		userID,
		body.Role)
	if err != nil {
		return answer, &schemes.ErrorResponse{Code: 105, Err: err.Error(), ErrBase: err}
	}

	answer.Role = body.Role
	answer.UserID = body.UserID
	answer.Nickname = user.Nickname
	return answer, nil
}

func (s *ServiceProjectMembers) DelProjectMembersFlow(ctx *gin.Context, uri *DTO.UriIDValid, body *DTO.MembersAddIdValid) (common.Empty, error) {
	var answer common.Empty

	ownerID, err := primitive.ObjectIDFromHex(ctx.GetString(middlewares.UserIDKey))
	if err != nil {
		return answer, &schemes.ErrorResponse{Code: 105, Err: "Bad owner id", ErrBase: nil}
	}
	userID, err := primitive.ObjectIDFromHex(body.UserID)
	if err != nil {
		return answer, &schemes.ErrorResponse{Code: 105, Err: "Bad user id", ErrBase: nil}
	}
	projectID, err := primitive.ObjectIDFromHex(uri.ID)
	if err != nil {
		return answer, &schemes.ErrorResponse{Code: 105, Err: "Bad project id", ErrBase: nil}
	}
	if ownerID == userID {
		return answer, &schemes.ErrorResponse{Code: 105, Err: "owner == user", ErrBase: nil}
	}

	err = s.pr.RemoveMemberFromProject(
		ctx.GetString(middlewares.Service), ctx.GetString(middlewares.Domain), projectID, ownerID, userID)
	if err != nil {
		return answer, &schemes.ErrorResponse{Code: 105, Err: err.Error(), ErrBase: err}
	}

	return answer, nil
}
