package usecase

import (
	entitiesDB2 "anubis/app/DAL/entitiesDB"
	interfacesDB2 "anubis/app/DAL/interfacesDB"
	storage2 "anubis/app/DAL/storage"
	"anubis/app/DTO"
	"anubis/app/core"
	"anubis/app/core/common"
	"anubis/app/core/middlewares"
	"anubis/app/core/schemes"
	"anubis/tools/utils"
	"errors"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type ServiceToken struct {
	pr     interfacesDB2.ProjectRepository
	usr    interfacesDB2.UserRepository
	config core.ServiceConfig
}

func NewServiceToken(
	pr *storage2.RepositoryMongoProjects,
	usr *storage2.RepositoryMongoUser,
	config core.ServiceConfig) *ServiceToken {
	return &ServiceToken{
		pr,
		usr,
		config,
	}
}

func (s *ServiceToken) RefreshTokenDomainFlow(ctx *gin.Context, input *DTO.RefreshTokenProjectValid) (DTO.AnswerToken, error) {
	var answer DTO.AnswerToken

	var userSession entitiesDB2.MdUsersSession
	sessionID := ctx.GetString(middlewares.SessionIDKey)
	hashToken := ctx.GetString(middlewares.HashToken)

	objectID, err := primitive.ObjectIDFromHex(sessionID)
	if err != nil {
		return answer, &schemes.ErrorResponse{Code: 107, Err: "BAD token", ErrBase: nil}
	}
	err = s.usr.GetUsersSessionByID(ctx.GetString(middlewares.Service), objectID, &userSession)
	if err != nil {
		return answer, &schemes.ErrorResponse{Code: 107, Err: "Couldn't find session", ErrBase: err}
	}
	if time.Now().After(userSession.ExpiresAt) {
		return answer, &schemes.ErrorResponse{Code: 107, Err: "Token has expired", ErrBase: nil}
	}

	if userSession.HashToken != hashToken {
		if userSession.IsActive == true {
			_ = s.usr.UserSessionsSetActive(ctx.GetString(middlewares.Service), userSession.ID, false)
		}
		return answer, &schemes.ErrorResponse{Code: 105, Err: "Bad Token hash", ErrBase: err}
	}

	if !userSession.IsActive {
		if err != nil {
			return answer, &schemes.ErrorResponse{Code: 107, Err: "Session deactivate", ErrBase: err}
		}
		return answer, &schemes.ErrorResponse{Code: 107, Err: "Couldn't find active session", ErrBase: err}
	}

	var project entitiesDB2.MdProject
	project.ID, err = primitive.ObjectIDFromHex(input.ProjectID)
	if err != nil {
		return answer, &schemes.ErrorResponse{Code: 105, Err: "Bad projectID", ErrBase: err}
	}
	err = s.pr.GetProjectsByUser(ctx.GetString(middlewares.Service), &project, userSession.UserID)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return answer, &schemes.ErrorResponse{Code: 105, Err: "Not found project for user ", ErrBase: err}
		}
		return answer, &schemes.ErrorResponse{Code: 105, Err: "Bad projectID", ErrBase: err}
	}

	expiresAt := time.Now().Add(time.Minute * time.Duration(s.config.RefreshTokenMinute))

	userSession.IP = common.GetClientIP(ctx)
	userSession.ExpiresAt = expiresAt
	refreshToken, err := utils.CreateRefreshToken(userSession.ID.Hex(), s.config.RefreshTokenSecret, s.config.RefreshTokenMinute)
	if err != nil {
		return answer, &schemes.ErrorResponse{Code: 97, Err: "Couldn't create a token", ErrBase: err}
	}

	userSession.HashToken = utils.HashToken(refreshToken)
	err = s.usr.UpdateSessionsByID(ctx.GetString(middlewares.Service), &userSession)
	if err != nil {
		return answer, &schemes.ErrorResponse{Code: 107, Err: "Session not update", ErrBase: err}
	}

	roleUser, _ := project.GetUserRole(userSession.UserID)
	accessToken, err := utils.CreateAccessToken(
		userSession.UserID.Hex(),
		project.ID.Hex()+"|"+roleUser,
		[]string{userSession.Domain},
		s.config.AccessTokenSecret,
		s.config.AccessTokenMinute)
	if err != nil {
		return answer, &schemes.ErrorResponse{Code: 97, Err: "Couldn't create a token", ErrBase: err}
	}

	if s.config.ShortJwt {
		answer.RefreshToken = utils.RemoveFirstPart(refreshToken)
		answer.AccessToken = utils.RemoveFirstPart(accessToken)
	} else {
		answer.RefreshToken = refreshToken
		answer.AccessToken = accessToken
	}
	if err != nil {
		return answer, &schemes.ErrorResponse{Code: 97, Err: "Couldn't create a token", ErrBase: err}
	}

	return answer, nil
}

func (s *ServiceToken) LogoutDomainFlow(ctx *gin.Context, input *DTO.LogoutValid) (schemes.EmptyResponses, error) {
	var answer schemes.EmptyResponses

	var userSession entitiesDB2.MdUsersSession
	sessionID := ctx.GetString(middlewares.SessionIDKey)
	hashToken := ctx.GetString(middlewares.HashToken)

	objectID, err := primitive.ObjectIDFromHex(sessionID)
	if err != nil {
		return answer, &schemes.ErrorResponse{Code: 107, Err: "BAD token", ErrBase: nil}
	}

	err = s.usr.GetUsersSessionByID(ctx.GetString(middlewares.Service), objectID, &userSession)
	if err != nil {
		return answer, &schemes.ErrorResponse{Code: 107, Err: "Couldn't find session", ErrBase: err}
	}
	if time.Now().After(userSession.ExpiresAt) {
		return answer, &schemes.ErrorResponse{Code: 107, Err: "Token has expired", ErrBase: nil}
	}

	if userSession.HashToken != hashToken {
		return answer, &schemes.ErrorResponse{Code: 105, Err: "Bad Token hash", ErrBase: err}
	}
	if !userSession.IsActive {
		if err != nil {
			return answer, &schemes.ErrorResponse{Code: 107, Err: "Session deactivate", ErrBase: err}
		}
		return answer, &schemes.ErrorResponse{Code: 107, Err: "Couldn't find active session", ErrBase: err}
	}

	if input.All {
		err = s.usr.DeactivateUserSessionsByDomain(
			ctx.GetString(middlewares.Service),
			userSession.UserID,
			userSession.Domain)
		if err != nil {
			return answer, &schemes.ErrorResponse{Code: 107, Err: "Session not all deactivate", ErrBase: err}
		}
	} else {
		err = s.usr.UserSessionsSetActive(ctx.GetString(middlewares.Service), userSession.ID, false)
		if err != nil {
			return answer, &schemes.ErrorResponse{Code: 107, Err: "Session not  deactivate", ErrBase: err}
		}
	}

	return answer, nil
}
