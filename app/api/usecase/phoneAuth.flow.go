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
	usr    interfacesDB.UserRepository
	pr     interfacesDB.ProjectRepository
	ath    interfacesDB.AuthPhoneRepository
	config core.ServiceConfig
}

func NewServiceAuth(
	usr *storage.RepositoryMongoUser,
	pr *storage.RepositoryMongoProjects,
	ath *storage.RepositoryMongoAuthPhone,
	config core.ServiceConfig) *ServiceAuth {
	return &ServiceAuth{
		usr,
		pr,
		ath,
		config,
	}
}

func (s *ServiceAuth) RegUserFlow(ctx *gin.Context, input DTO.PhoneValidUserReg) (DTO.AnswerUserReg, error) {
	//TODO.MD:  нет проверки КАПЧИ
	var answer DTO.AnswerUserReg

	service, err := helpers.CheckDomain(s.config, input.Domain)
	if err != nil {
		return answer, err
	}

	//
	var phone = &entitiesDB.MdPhoneAuth{}
	err = helpers.FillPhoneReg(phone, input)
	if err != nil {
		return answer, &schemes.ErrorResponse{Code: 105, Err: "Bad phone", ErrBase: err}
	}

	err = s.ath.GetPhone(service, phone)
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		return answer, &schemes.ErrorResponse{Code: 105, Err: "Bad phone for user", ErrBase: err}
	} else if phone.Verification != false {
		return answer, &schemes.ErrorResponse{Code: 109, Err: "The user already exists", ErrBase: err}
	}

	var user entitiesDB.MdUser
	if input.Nickname != "" {
		user.Nickname = input.Nickname
	} else {
		user.Nickname = "name_" + utils.RandStringBytes(10)
	}

	if phone.UserID.IsZero() {
		err = s.usr.CreateUser(service, &user)
		if err != nil {
			return answer, &schemes.ErrorResponse{Code: 105, Err: "Bad Req", ErrBase: err}
		}
		phone.UserID = user.ID
	}

	err = s.ath.SavePhone(service, phone)
	if err != nil {
		return answer, &schemes.ErrorResponse{Code: 109, Err: "Not save phone", ErrBase: err}
	}
	if phone.Verification == true {
		return answer, &schemes.ErrorResponse{Code: 109, Err: "The phone already exists", ErrBase: err}
	}

	var mySms entitiesDB.MdSmsAuth
	mySms.SmsCode = utils.RandStringBytes(6)
	mySms.IDSend, mySms.SmsService, err = sms.Sender(mySms.SmsCode)
	mySms.Phone = phone.Phone
	mySms.UserID = phone.UserID
	err = s.ath.SaveSmsAuth(service, &mySms)
	if err != nil {
		return answer, &schemes.ErrorResponse{Code: 105, Err: "Bad sms", ErrBase: err}
	}
	answer.SmsId = mySms.ID.Hex()
	return answer, nil
}

func (s *ServiceAuth) ValidSmsUserFlow(ctx *gin.Context, input DTO.ValidSms) (DTO.AnswerRegToken, error) {
	var answer DTO.AnswerRegToken

	service, err := helpers.CheckDomain(s.config, input.Domain)
	if err != nil {
		return answer, err
	}

	objectID, err := primitive.ObjectIDFromHex(input.SmsId)
	if err != nil {
		return answer, &schemes.ErrorResponse{Code: 105, Err: "Bad SmsId", ErrBase: err}
	}

	tenMinutesAgo := time.Now().Add(-10 * time.Minute)
	if !objectID.Timestamp().After(tenMinutesAgo) {
		return answer, &schemes.ErrorResponse{Code: 105, Err: "The time for the SMS code has expired.", ErrBase: nil}
	}

	phoneNum, err := s.ath.SmsValidUser(service, objectID, input.SmsCode)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return answer, &schemes.ErrorResponse{Code: 105, Err: "User with sms-code not found", ErrBase: err}
		}
		return answer, err
	}

	var phone = &entitiesDB.MdPhoneAuth{}
	phone.Phone = phoneNum
	err = s.ath.GetPhone(service, phone)

	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		return answer, &schemes.ErrorResponse{Code: 105, Err: "Bad phone for user", ErrBase: err}
	} else if phone.Verification != false {
		return answer, &schemes.ErrorResponse{Code: 109, Err: "The user already exists", ErrBase: err}
	}

	err = s.ath.SaveVerifyPhone(service, true, phone.Phone)
	if err != nil {
		return answer, &schemes.ErrorResponse{Code: 109, Err: "Not save verification", ErrBase: err}
	}

	projectMap, err := s.pr.GetProjectsListByUser(service, input.Domain, phone.UserID)
	if err != nil {
		return answer, &schemes.ErrorResponse{Code: 105, Err: "Couldn't create a project", ErrBase: err}
	}

	expiresAt := time.Now().Add(time.Minute * time.Duration(20))
	usersSession := entitiesDB.MdUsersSession{
		DeviceId:   "_",
		DeviceType: common.GetDeviceType(ctx),
		UserID:     phone.UserID,
		Domain:     "*******",
		CreatedAt:  time.Now(),
		ExpiresAt:  expiresAt,
		IP:         common.GetClientIP(ctx),
		IsActive:   true,
	}
	err = s.usr.DeactivateOldAndCreateSession(service, &usersSession)
	if err != nil {
		return answer, &schemes.ErrorResponse{Code: 105, Err: "Couldn't create a session", ErrBase: err}
	}
	answer.ListProjects = projectMap
	//create TOKEN----------------------------------------------------
	refreshToken, err := utils.CreateRefreshToken(usersSession.ID.Hex(), s.config.RefreshTokenSecret, 19)
	if err != nil {
		return answer, &schemes.ErrorResponse{Code: 97, Err: "Couldn't create a token", ErrBase: err}
	}
	if s.config.ShortJwt {
		answer.RefreshToken = utils.RemoveFirstPart(refreshToken)
	} else {
		answer.RefreshToken = refreshToken
	}
	//-----------------------------------------------------------------
	if err != nil {
		return answer, &schemes.ErrorResponse{Code: 97, Err: "Couldn't create a token", ErrBase: err}
	}

	return answer, nil
}

func (s *ServiceAuth) PhoneLoginFlow(ctx *gin.Context, input DTO.PhoneValidUserReg) (DTO.AnswerRegToken, error) {
	//TODO.MD:  нет проверки КАПЧИ
	var answer DTO.AnswerRegToken

	service, err := helpers.CheckDomain(s.config, input.Domain)
	if err != nil {
		return answer, err
	}

	number, err := strconv.ParseInt(input.Phone, 10, 64)
	if err != nil {
		return answer, &schemes.ErrorResponse{Code: 105, Err: "Bad phone", ErrBase: err}
	}

	UserID, hashedPassword, err := s.ath.GetPhoneUserID(service, number)
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

	projectMap, err := s.pr.GetProjectsListByUser(service, input.Domain, UserID)
	if err != nil {
		return answer, &schemes.ErrorResponse{Code: 105, Err: "Couldn't create a project", ErrBase: err}
	}

	expiresAt := time.Now().Add(time.Minute * time.Duration(20))
	usersSession := entitiesDB.MdUsersSession{
		DeviceId:   "_",
		DeviceType: common.GetDeviceType(ctx),
		UserID:     UserID,
		Domain:     "*******",
		CreatedAt:  time.Now(),
		ExpiresAt:  expiresAt,
		IP:         common.GetClientIP(ctx),
		IsActive:   true,
	}
	err = s.usr.DeactivateOldAndCreateSession(service, &usersSession)
	if err != nil {
		return answer, &schemes.ErrorResponse{Code: 105, Err: "Couldn't create a session", ErrBase: err}
	}
	answer.ListProjects = projectMap
	//create TOKEN----------------------------------------------------
	refreshToken, err := utils.CreateRefreshToken(usersSession.ID.Hex(), s.config.RefreshTokenSecret, 19)
	if err != nil {
		return answer, &schemes.ErrorResponse{Code: 97, Err: "Couldn't create a token", ErrBase: err}
	}
	if s.config.ShortJwt {
		answer.RefreshToken = utils.RemoveFirstPart(refreshToken)
	} else {
		answer.RefreshToken = refreshToken
	}
	//-----------------------------------------------------------------
	if err != nil {
		return answer, &schemes.ErrorResponse{Code: 97, Err: "Couldn't create a token", ErrBase: err}
	}

	return answer, nil
}
