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
	"anubis/tools/providers"
	"anubis/tools/utils"
	"errors"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type ServiceEmailAuth struct {
	usr    interfacesDB2.UserRepository
	pr     interfacesDB2.ProjectRepository
	ath    interfacesDB2.AuthEmailRepository
	config core.ServiceConfig
}

func NewServiceEmailAuth(
	usr *storage2.RepositoryMongoUser,
	pr *storage2.RepositoryMongoProjects,
	ath *storage2.RepositoryMongoAuthEmail,
	config core.ServiceConfig) *ServiceEmailAuth {
	return &ServiceEmailAuth{
		usr,
		pr,
		ath,
		config,
	}
}

func (s *ServiceEmailAuth) RegUserEmailFlow(ctx *gin.Context, input *DTO.EmailUserRegValid) (DTO.AnswerUserRegCode, error) {
	//TODO.MD:  нет проверки КАПЧИ
	var answer DTO.AnswerUserRegCode

	//
	var email = &entitiesDB2.MdEmailAuth{}

	email.Email = input.Email
	email.PasswordHash, _ = utils.GeneratePasswordHash(input.Password)
	email.Verification = false
	err := s.ath.GetEmail(ctx.GetString(middlewares.Service), email)
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		return answer, &schemes.ErrorResponse{Code: 105, Err: "Bad email for user", ErrBase: err}
	} else if email.Verification != false {
		return answer, &schemes.ErrorResponse{Code: 109, Err: "The user already exists", ErrBase: err}
	}

	var user entitiesDB2.MdUser
	if input.Nickname != "" {
		user.Nickname = input.Nickname
	} else {
		user.Nickname = "name_" + utils.RandStringBytes(10)
	}

	if email.UserID.IsZero() {
		err = s.usr.CreateUser(ctx.GetString(middlewares.Service), &user)
		if err != nil {
			return answer, &schemes.ErrorResponse{Code: 105, Err: "Bad Req", ErrBase: err}
		}
		email.UserID = user.ID
	}

	err = s.ath.SaveEmail(ctx.GetString(middlewares.Service), email)
	if err != nil {
		return answer, &schemes.ErrorResponse{Code: 109, Err: "Not save phone", ErrBase: err}
	}
	if email.Verification == true {
		return answer, &schemes.ErrorResponse{Code: 109, Err: "The email already exists", ErrBase: err}
	}

	var myEmailCode entitiesDB2.MdEmailCodeAuth
	myEmailCode.EmailCode = utils.RandStringBytes(6)
	myEmailCode.IDSend, myEmailCode.EmailService, err = providers.SenderEmail(myEmailCode.EmailCode)
	myEmailCode.Email = email.Email
	myEmailCode.UserID = email.UserID

	err = s.ath.SaveEmailAuth(ctx.GetString(middlewares.Service), &myEmailCode)
	if err != nil {
		return answer, &schemes.ErrorResponse{Code: 105, Err: "Bad sms", ErrBase: err}
	}
	answer.CodeId = myEmailCode.ID.Hex()
	return answer, nil
}

func (s *ServiceEmailAuth) ValidCodeEmailUserFlow(ctx *gin.Context, input *DTO.CodeEmailValid) (DTO.AnswerRegToken, error) {
	var answer DTO.AnswerRegToken

	objectID, err := primitive.ObjectIDFromHex(input.CodeId)
	if err != nil {
		return answer, &schemes.ErrorResponse{Code: 105, Err: "Bad EmailCodeId", ErrBase: err}
	}

	tenMinutesAgo := time.Now().Add(-10 * time.Minute)
	if !objectID.Timestamp().After(tenMinutesAgo) {
		return answer, &schemes.ErrorResponse{Code: 105, Err: "The time for the Email code has expired.", ErrBase: nil}
	}

	mail, err := s.ath.EmailCodeValidUser(ctx.GetString(middlewares.Service), objectID, input.EmailCode)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return answer, &schemes.ErrorResponse{Code: 104, Err: "User with email-code not found", ErrBase: err}
		}
		return answer, err
	}

	var email = &entitiesDB2.MdEmailAuth{}
	email.Email = mail
	err = s.ath.GetEmail(ctx.GetString(middlewares.Service), email)

	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		return answer, &schemes.ErrorResponse{Code: 105, Err: "Bad email for user", ErrBase: err}
	} else if email.Verification != false {
		return answer, &schemes.ErrorResponse{Code: 109, Err: "The user already exists", ErrBase: err}
	}

	err = s.ath.SaveVerifyEmail(ctx.GetString(middlewares.Service), true, email.Email)
	if err != nil {
		return answer, &schemes.ErrorResponse{Code: 109, Err: "Not save verification", ErrBase: err}
	}

	projectMap, err := s.pr.GetProjectsListByUser(ctx.GetString(middlewares.Service), ctx.GetString(middlewares.Domain), email.UserID)
	if err != nil {
		return answer, &schemes.ErrorResponse{Code: 105, Err: "Couldn't create a project", ErrBase: err}
	}

	expiresAt := time.Now().Add(time.Minute * time.Duration(20))
	userSession := entitiesDB2.MdUsersSession{
		DeviceId:   input.DeviceId,
		DeviceType: common.GetDeviceType(ctx),
		UserID:     email.UserID,
		Domain:     ctx.GetString(middlewares.Domain),
		CreatedAt:  time.Now(),
		ExpiresAt:  expiresAt,
		IP:         common.GetClientIP(ctx),
		IsActive:   true,
	}
	err = s.usr.DeactivateOldAndCreateSession(ctx.GetString(middlewares.Service), &userSession)
	if err != nil {
		return answer, &schemes.ErrorResponse{Code: 107, Err: "Couldn't create a session", ErrBase: err}
	}
	answer.ListProjects = projectMap
	//create TOKEN----------------------------------------------------
	refreshToken, err := utils.CreateRefreshToken(
		userSession.ID.Hex(),
		s.config.RefreshTokenSecret,
		s.config.RefreshTokenMinute)
	if err != nil {
		return answer, &schemes.ErrorResponse{Code: 97, Err: "Couldn't create a token", ErrBase: err}
	}

	userSession.HashToken = utils.HashToken(refreshToken)
	err = s.usr.UpdateSessionsByID(ctx.GetString(middlewares.Service), &userSession)
	if err != nil {
		return answer, &schemes.ErrorResponse{Code: 107, Err: "Session not update", ErrBase: err}
	}

	accessToken, err := utils.CreateAccessToken(
		userSession.UserID.Hex(),
		"",
		[]string{userSession.Domain},
		s.config.AccessTokenSecret,
		s.config.AccessTokenMinute)
	if err != nil {
		return answer, &schemes.ErrorResponse{Code: 97, Err: "Couldn't create a token", ErrBase: err}
	}

	answer.RefreshToken = refreshToken
	answer.AccessToken = accessToken

	return answer, nil
}

func (s *ServiceEmailAuth) EmailLoginFlow(ctx *gin.Context, input *DTO.LoginEmailUserValid) (DTO.AnswerRegToken, error) {
	//TODO.MD:  нет проверки КАПЧИ
	var answer DTO.AnswerRegToken

	UserID, hashedPassword, err := s.ath.GetEmailVerificationUserID(
		ctx.GetString(middlewares.Service),
		input.Email,
		true)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return answer, &schemes.ErrorResponse{Code: 104, Err: "User not found :(", ErrBase: err}
		} else {
			return answer, err
		}

	}
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(input.Password))
	if err != nil {
		return answer, &schemes.ErrorResponse{Code: 105, Err: "Invalid username or password", ErrBase: err}
	}

	projectMap, err := s.pr.GetProjectsListByUser(
		ctx.GetString(middlewares.Service),
		ctx.GetString(middlewares.Domain),
		UserID)
	if err != nil {
		return answer, &schemes.ErrorResponse{Code: 105, Err: "Couldn't create a project", ErrBase: err}
	}

	expiresAt := time.Now().Add(time.Minute * time.Duration(20))
	userSession := entitiesDB2.MdUsersSession{
		DeviceId:   input.DeviceId,
		DeviceType: common.GetDeviceType(ctx),
		UserID:     UserID,
		Domain:     ctx.GetString(middlewares.Domain),
		CreatedAt:  time.Now(),
		ExpiresAt:  expiresAt,
		IP:         common.GetClientIP(ctx),
		IsActive:   true,
	}
	err = s.usr.CheckOldAndUpdateSession(ctx.GetString(middlewares.Service), &userSession)
	if err != nil {
		return answer, &schemes.ErrorResponse{Code: 107, Err: "Couldn't create a session", ErrBase: err}
	}
	answer.ListProjects = projectMap
	//create TOKEN----------------------------------------------------
	refreshToken, err := utils.CreateRefreshToken(
		userSession.ID.Hex(),
		s.config.RefreshTokenSecret,
		s.config.RefreshTokenMinute)
	if err != nil {
		return answer, &schemes.ErrorResponse{Code: 97, Err: "Couldn't create a token", ErrBase: err}
	}

	userSession.HashToken = utils.HashToken(refreshToken)
	err = s.usr.UpdateSessionsByID(ctx.GetString(middlewares.Service), &userSession)
	if err != nil {
		return answer, &schemes.ErrorResponse{Code: 107, Err: "Session not update", ErrBase: err}
	}

	accessToken, err := utils.CreateAccessToken(
		userSession.UserID.Hex(),
		"",
		[]string{userSession.Domain},
		s.config.AccessTokenSecret,
		s.config.AccessTokenMinute)
	if err != nil {
		return answer, &schemes.ErrorResponse{Code: 97, Err: "Couldn't create a token", ErrBase: err}
	}

	answer.RefreshToken = refreshToken
	answer.AccessToken = accessToken
	return answer, nil
}
