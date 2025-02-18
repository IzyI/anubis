package usecase

import (
	"anubis/app/api/DAL/entitiesDB"
	"anubis/app/api/DAL/interfacesDB"
	"anubis/app/api/DAL/storage"
	"anubis/app/api/DTO"
	"anubis/app/api/helpers"
	"anubis/app/core"
	"anubis/app/core/common"
	"anubis/app/core/schemes"
	"anubis/tools/utils"
	"errors"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type ServiceToken struct {
	pr     interfacesDB.ProjectRepository
	usr    interfacesDB.UserRepository
	config core.ServiceConfig
}

func NewServiceToken(
	pr *storage.RepositoryMongoProjects,
	usr *storage.RepositoryMongoUser,
	config core.ServiceConfig) *ServiceToken {
	return &ServiceToken{
		pr,
		usr,
		config,
	}
}

func (s *ServiceToken) RefreshTokenDomainFlow(ctx *gin.Context, input DTO.RefreshTokenProjectI) (DTO.AnswerToken, error) {
	var answer DTO.AnswerToken

	service, err := helpers.CheckDomain(s.config, input.Domain)
	if err != nil {
		return answer, err
	}

	token := input.RefreshToken
	if s.config.ShortJwt {
		token = s.config.ShortJwtValue + "." + input.RefreshToken
	}
	authorized, err := utils.IsAuthorized(token, s.config.RefreshTokenSecret)
	if !authorized {
		return answer, &schemes.ErrorResponse{Code: 98, Err: "Not authorized", ErrBase: err}

	}

	var refreshClaims utils.RefreshClaims
	err = utils.ExtractToken(token, s.config.RefreshTokenSecret, &refreshClaims)
	if err != nil {
		return answer, &schemes.ErrorResponse{Code: 98, Err: "Not extract token", ErrBase: err}
	}
	objectID, err := primitive.ObjectIDFromHex(refreshClaims.ID)
	if err != nil {
		return answer, &schemes.ErrorResponse{Code: 105, Err: "Bad SessionId", ErrBase: err}
	}

	var userSession entitiesDB.MdUsersSession
	err = s.usr.GetUsersSessionByID(service, objectID, &userSession)
	if err != nil {
		return answer, &schemes.ErrorResponse{Code: 107, Err: "Couldn't find session", ErrBase: err}
	}
	if time.Now().After(userSession.ExpiresAt) {
		return answer, &schemes.ErrorResponse{Code: 107, Err: "Token has expired", ErrBase: nil}
	}
	oldDeviceId := userSession.DeviceId
	if oldDeviceId != "_" {
		if userSession.HashToken != utils.HashToken(token) {
			if userSession.IsActive == true {
				_ = s.usr.UserSessionsSetActive(service, &userSession, false)
			}
			return answer, &schemes.ErrorResponse{Code: 105, Err: "Bad Token hash", ErrBase: err}
		}
		if userSession.DeviceId != input.DeviceId {
			return answer, &schemes.ErrorResponse{Code: 105, Err: "Bad Device", ErrBase: err}
		}
	}

	if !userSession.IsActive {
		if err != nil {
			return answer, &schemes.ErrorResponse{Code: 107, Err: "Session deactivate", ErrBase: err}
		}
		return answer, &schemes.ErrorResponse{Code: 107, Err: "Couldn't find active session", ErrBase: err}
	}

	var project entitiesDB.MdProject

	project.ID, err = primitive.ObjectIDFromHex(input.ProjectID)
	if err != nil {
		return answer, &schemes.ErrorResponse{Code: 105, Err: "Bad projectID", ErrBase: err}
	}
	err = s.pr.GetProjectsByUser(service, &project, userSession.UserID)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return answer, &schemes.ErrorResponse{Code: 105, Err: "Not found project for user ", ErrBase: err}
		}
		return answer, &schemes.ErrorResponse{Code: 105, Err: "Bad projectID", ErrBase: err}
	}

	expiresAt := time.Now().Add(time.Minute * time.Duration(s.config.RefreshTokenMinute))

	userSession.Domain = input.Domain
	userSession.DeviceType = common.GetDeviceType(ctx)
	userSession.IP = common.GetClientIP(ctx)
	userSession.ExpiresAt = expiresAt
	userSession.DeviceId = input.DeviceId
	refreshToken, err := utils.CreateRefreshToken(userSession.ID.Hex(), s.config.RefreshTokenSecret, s.config.RefreshTokenMinute)
	if err != nil {
		return answer, &schemes.ErrorResponse{Code: 97, Err: "Couldn't create a token", ErrBase: err}
	}

	userSession.HashToken = utils.HashToken(refreshToken)
	if oldDeviceId == "_" {
		err = s.usr.CheckOldAndUpdateSession(service, &userSession)
		if err != nil {
			return answer, &schemes.ErrorResponse{Code: 107, Err: "ManySession not update", ErrBase: err}
		}
	} else {
		err = s.usr.UpdateSessionsByID(service, &userSession)
		if err != nil {
			return answer, &schemes.ErrorResponse{Code: 107, Err: "Session not update", ErrBase: err}
		}
	}

	roleUser, _ := project.GetUserRole(userSession.UserID)
	accessToken, err := utils.CreateAccessToken(
		userSession.UserID.Hex()+"|"+roleUser,
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

func (s *ServiceToken) LogoutDomainFlow(ctx *gin.Context, input DTO.Logout) (schemes.EmptyResponses, error) {
	var answer schemes.EmptyResponses

	service, err := helpers.CheckDomain(s.config, input.Domain)
	if err != nil {
		return answer, err
	}

	token := input.RefreshToken
	if s.config.ShortJwt {
		token = s.config.ShortJwtValue + "." + input.RefreshToken
	}
	authorized, err := utils.IsAuthorized(token, s.config.RefreshTokenSecret)
	if !authorized {
		return answer, &schemes.ErrorResponse{Code: 98, Err: "Not authorized", ErrBase: err}

	}

	var refreshClaims utils.RefreshClaims
	err = utils.ExtractToken(token, s.config.RefreshTokenSecret, &refreshClaims)
	if err != nil {
		return answer, &schemes.ErrorResponse{Code: 98, Err: "Not extract token", ErrBase: err}
	}
	objectID, err := primitive.ObjectIDFromHex(refreshClaims.ID)
	if err != nil {
		return answer, &schemes.ErrorResponse{Code: 105, Err: "Bad SessionId", ErrBase: err}
	}

	var userSession entitiesDB.MdUsersSession
	err = s.usr.GetUsersSessionByID(service, objectID, &userSession)
	if err != nil {
		return answer, &schemes.ErrorResponse{Code: 107, Err: "Couldn't find session", ErrBase: err}
	}
	if time.Now().After(userSession.ExpiresAt) {
		return answer, &schemes.ErrorResponse{Code: 107, Err: "Token has expired", ErrBase: nil}
	}

	if userSession.HashToken != utils.HashToken(token) {
		return answer, &schemes.ErrorResponse{Code: 105, Err: "Bad Token hash", ErrBase: err}
	}
	if !userSession.IsActive {
		if err != nil {
			return answer, &schemes.ErrorResponse{Code: 107, Err: "Session deactivate", ErrBase: err}
		}
		return answer, &schemes.ErrorResponse{Code: 107, Err: "Couldn't find active session", ErrBase: err}
	}

	if input.All {
		err = s.usr.DeactivateUserSessionsByDomain(service, &userSession)
		if err != nil {
			return answer, &schemes.ErrorResponse{Code: 107, Err: "Session not all deactivate", ErrBase: err}
		}
	} else {
		err = s.usr.UserSessionsSetActive(service, &userSession, false)
		if err != nil {
			return answer, &schemes.ErrorResponse{Code: 107, Err: "Session not  deactivate", ErrBase: err}
		}
	}

	return answer, nil
}
