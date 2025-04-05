package usecase

import (
	entitiesDB2 "anubis/app/DAL/entitiesDB"
	interfacesDB2 "anubis/app/DAL/interfacesDB"
	storage2 "anubis/app/DAL/storage"
	"anubis/app/DTO"
	"anubis/app/api/helpers"
	"anubis/app/core"
	"anubis/app/core/common"
	"anubis/app/core/middlewares"
	"anubis/app/core/schemes"
	"anubis/tools/providers/sms"
	"anubis/tools/utils"
	"errors"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"strconv"
	"time"
)

type ServiceAuth struct {
	usr    interfacesDB2.UserRepository
	pr     interfacesDB2.ProjectRepository
	ath    interfacesDB2.AuthPhoneRepository
	config core.ServiceConfig
}

func NewServiceAuth(
	usr *storage2.RepositoryMongoUser,
	pr *storage2.RepositoryMongoProjects,
	ath *storage2.RepositoryMongoAuthPhone,
	config core.ServiceConfig) *ServiceAuth {
	return &ServiceAuth{
		usr,
		pr,
		ath,
		config,
	}
}

func (s *ServiceAuth) RegUserFlow(ctx *gin.Context, input *DTO.PhoneUserRegValid) (DTO.AnswerUserReg, error) {
	//TODO.MD:  нет проверки КАПЧИ
	var answer DTO.AnswerUserReg

	//
	var phone = &entitiesDB2.MdPhoneAuth{}
	err := helpers.FillPhoneReg(phone, input)
	if err != nil {
		return answer, &schemes.ErrorResponse{Code: 105, Err: "Bad phone", ErrBase: err}
	}

	err = s.ath.GetPhone(ctx.GetString(middlewares.Service), phone)
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		return answer, &schemes.ErrorResponse{Code: 105, Err: "Bad phone for user", ErrBase: err}
	} else if phone.Verification != false {
		return answer, &schemes.ErrorResponse{Code: 109, Err: "The user already exists", ErrBase: err}
	}

	var user entitiesDB2.MdUser
	if input.Nickname != "" {
		user.Nickname = input.Nickname
	} else {
		user.Nickname = "name_" + utils.RandStringBytes(10)
	}

	if phone.UserID.IsZero() {
		err = s.usr.CreateUser(ctx.GetString(middlewares.Service), &user)
		if err != nil {
			return answer, &schemes.ErrorResponse{Code: 105, Err: "Bad Req", ErrBase: err}
		}
		phone.UserID = user.ID
	}

	err = s.ath.SavePhone(ctx.GetString(middlewares.Service), phone)
	if err != nil {
		return answer, &schemes.ErrorResponse{Code: 109, Err: "Not save phone", ErrBase: err}
	}
	if phone.Verification == true {
		return answer, &schemes.ErrorResponse{Code: 109, Err: "The phone already exists", ErrBase: err}
	}

	var mySms entitiesDB2.MdSmsAuth
	mySms.SmsCode = utils.RandStringBytes(6)
	mySms.IDSend, mySms.SmsService, err = sms.Sender(mySms.SmsCode)
	mySms.Phone = phone.Phone
	mySms.UserID = phone.UserID
	err = s.ath.SaveSmsAuth(ctx.GetString(middlewares.Service), &mySms)
	if err != nil {
		return answer, &schemes.ErrorResponse{Code: 105, Err: "Bad sms", ErrBase: err}
	}
	answer.SmsId = mySms.ID.Hex()
	return answer, nil
}

func (s *ServiceAuth) ValidSmsUserFlow(ctx *gin.Context, input *DTO.SmsValid) (DTO.AnswerRegToken, error) {
	var answer DTO.AnswerRegToken

	objectID, err := primitive.ObjectIDFromHex(input.SmsId)
	if err != nil {
		return answer, &schemes.ErrorResponse{Code: 105, Err: "Bad SmsId", ErrBase: err}
	}

	tenMinutesAgo := time.Now().Add(-10 * time.Minute)
	if !objectID.Timestamp().After(tenMinutesAgo) {
		return answer, &schemes.ErrorResponse{Code: 105, Err: "The time for the SMS code has expired.", ErrBase: nil}
	}

	phoneNum, err := s.ath.SmsValidUser(ctx.GetString(middlewares.Service), objectID, input.SmsCode)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return answer, &schemes.ErrorResponse{Code: 105, Err: "User with sms-code not found", ErrBase: err}
		}
		return answer, err
	}

	var phone = &entitiesDB2.MdPhoneAuth{}
	phone.Phone = phoneNum
	err = s.ath.GetPhone(ctx.GetString(middlewares.Service), phone)

	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		return answer, &schemes.ErrorResponse{Code: 105, Err: "Bad phone for user", ErrBase: err}
	} else if phone.Verification != false {
		return answer, &schemes.ErrorResponse{Code: 109, Err: "The user already exists", ErrBase: err}
	}

	err = s.ath.SaveVerifyPhone(ctx.GetString(middlewares.Service), true, phone.Phone)
	if err != nil {
		return answer, &schemes.ErrorResponse{Code: 109, Err: "Not save verification", ErrBase: err}
	}

	projectMap, err := s.pr.GetProjectsListByUser(ctx.GetString(middlewares.Service), ctx.GetString(middlewares.Domain), phone.UserID)
	if err != nil {
		return answer, &schemes.ErrorResponse{Code: 105, Err: "Couldn't create a project", ErrBase: err}
	}

	expiresAt := time.Now().Add(time.Minute * time.Duration(20))
	userSession := entitiesDB2.MdUsersSession{
		DeviceId:   input.DeviceId,
		DeviceType: common.GetDeviceType(ctx),
		UserID:     phone.UserID,
		Domain:     ctx.GetString(middlewares.Domain),
		CreatedAt:  time.Now(),
		ExpiresAt:  expiresAt,
		IP:         common.GetClientIP(ctx),
		IsActive:   true,
	}
	err = s.usr.DeactivateOldAndCreateSession(ctx.GetString(middlewares.Service), &userSession)
	if err != nil {
		return answer, &schemes.ErrorResponse{Code: 105, Err: "Couldn't create a session", ErrBase: err}
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

	if s.config.ShortJwt {
		answer.RefreshToken = utils.RemoveFirstPart(refreshToken)
		answer.AccessToken = utils.RemoveFirstPart(accessToken)
	} else {
		answer.RefreshToken = refreshToken
		answer.AccessToken = accessToken
	}

	return answer, nil
}

func (s *ServiceAuth) PhoneLoginFlow(ctx *gin.Context, input *DTO.LoginUserValid) (DTO.AnswerRegToken, error) {
	//TODO.MD:  нет проверки КАПЧИ
	var answer DTO.AnswerRegToken

	number, err := strconv.ParseInt(input.Phone, 10, 64)
	if err != nil {
		return answer, &schemes.ErrorResponse{Code: 105, Err: "Bad phone", ErrBase: err}
	}

	UserID, hashedPassword, err := s.ath.GetPhoneVerificationUserID(
		ctx.GetString(middlewares.Service),
		number,
		true)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return answer, &schemes.ErrorResponse{Code: 105, Err: "User not found :(", ErrBase: err}
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
		return answer, &schemes.ErrorResponse{Code: 105, Err: "Couldn't create a session", ErrBase: err}
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

	if s.config.ShortJwt {
		answer.RefreshToken = utils.RemoveFirstPart(refreshToken)
		answer.AccessToken = utils.RemoveFirstPart(accessToken)
	} else {
		answer.RefreshToken = refreshToken
		answer.AccessToken = accessToken
	}
	return answer, nil
}
