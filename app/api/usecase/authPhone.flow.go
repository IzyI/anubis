package usecase

import (
	"anubis/app/api/DAL/entitiesDB"
	"anubis/app/api/DAL/interfacesDB"
	"anubis/app/api/DAL/storage"
	"anubis/app/api/helpers"
	schemesAuth "anubis/app/api/schemes"
	"anubis/app/core"
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

func (s *ServiceAuth) RegUserFlow(ctx *gin.Context, input schemesAuth.PhoneValidUserReg) (schemesAuth.AnswerUserReg, error) {
	var answer schemesAuth.AnswerUserReg
	service, err := helpers.CheckDomain(s.config, input.Domain)
	if err != nil {
		return answer, err
	}
	//TODO.MD:  нет проверки КАПЧИ

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

	//TODO.MD:  установить триггер на удаление старого юзера phone.UserUuid вот в чем вопрос
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

func (s *ServiceAuth) ValidSmsUserFlow(ctx *gin.Context, input schemesAuth.ValidSms) (schemesAuth.RegAnswerToken, error) {
	var answer schemesAuth.RegAnswerToken
	service, err := helpers.CheckDomain(s.config, input.Domain)
	if err != nil {
		return answer, err
	}
	objectID, err := primitive.ObjectIDFromHex(input.SmsId)
	if err != nil {
		return answer, &schemes.ErrorResponse{Code: 105, Err: "Bad SmsId", ErrBase: err}
	}
	objectIDTimestamp := objectID.Timestamp()
	tenMinutesAgo := time.Now().Add(-10 * time.Minute)
	if !objectIDTimestamp.After(tenMinutesAgo) {
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
	var tokens schemesAuth.AnswerToken
	err = helpers.FillJWTTokens(&tokens, phone.UserID.Hex(), []string{}, s.config, 18, 1)
	if err != nil {
		return answer, &schemes.ErrorResponse{Code: 97, Err: "Couldn't create a token", ErrBase: err}
	}
	groupMap, err := s.usr.GetGroupUser(service, input.Domain, phone.UserID)
	if err != nil {
		return answer, &schemes.ErrorResponse{Code: 97, Err: "Couldn't create a token", ErrBase: err}
	}
	answer.GroupId = groupMap
	//answer.RefreshToken = tokens.RefreshToken
	answer.AccessToken = tokens.AccessToken
	return answer, nil
}

func (s *ServiceAuth) PhoneValidUserReg(ctx *gin.Context, input schemesAuth.PhoneValidUserReg) (schemesAuth.AnswerToken, error) {
	var answer schemesAuth.AnswerToken
	//TODO.MD:  нет проверки КАПЧИ
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

	err = helpers.FillJWTTokens(&answer, UserID.Hex(), []string{input.Domain}, s.config, 0, 0)
	if err != nil {
		return answer, &schemes.ErrorResponse{Code: 97, Err: "Couldn't create a token", ErrBase: err}
	}

	return answer, nil
}

//
//func (s *ServiceAuth) RefreshTokenUserFlow(ctx *gin.Context, input schemesAuth.ValidRefresh) (schemesAuth.AnswerToken, error) {
//	var token schemesAuth.AnswerToken
//	//
//	//authorized, _ := utils.IsAuthorized(input.RefreshToken, s.config.RefreshTokenSecret)
//	//if !authorized {
//	//	return token, &schemes.ErrorResponse{Code: 98, Err: "Not authorized"}
//	//
//	//}
//	//
//	//var userID, err = utils.ExtractToken(input.RefreshToken, s.config.RefreshTokenSecret)
//	//if err != nil {
//	//	return token, &schemes.ErrorResponse{Code: 98, Err: "Not find User"}
//	//}
//	//
//	//err = s.ath.GetUuidUser(userID)
//	//if err != nil {
//	//	return token, &schemes.ErrorResponse{Code: 105, Err: "User not found"}
//	//}
//	//
//	//accessToken, err := utils.CreateAccessToken(userID, s.config.AccessTokenSecret, s.config.AccessTokenHour)
//	//if err != nil {
//	//	return token, &schemes.ErrorResponse{Code: 96, Err: "Couldn't create a token"}
//	//}
//	//
//	//refreshToken, err := utils.CreateRefreshToken(userID, s.config.RefreshTokenSecret, s.config.AccessTokenHour)
//	//if err != nil {
//	//	return token, &schemes.ErrorResponse{Code: 97, Err: "Couldn't create a token"}
//	//}
//	//token.AccessToken = accessToken
//	//token.RefreshToken = refreshToken
//	return token, nil
//}
