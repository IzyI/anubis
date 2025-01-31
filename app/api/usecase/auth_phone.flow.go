package usecase

import (
	"anubis/app/api/entities"
	"anubis/app/api/helpers"
	schemesAuth "anubis/app/api/schemes"
	"anubis/app/api/storage"
	"anubis/app/core"
	"anubis/app/core/schemes"
	"anubis/tools/providers/sms"
	"anubis/tools/utils"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"strings"
)

type ServiceAuth struct {
	ath    entities.AuthPhoneRepository
	usr    entities.UserRepository
	config core.ServiceConfig
}

func NewServiceAuth(ath *storage.RepositoryPsqlAuthPhone, usr *storage.RepositoryPsqlUser, config core.ServiceConfig) *ServiceAuth {
	return &ServiceAuth{ath, usr, config}
}

func (s *ServiceAuth) RegUserFlow(ctx *gin.Context, input schemesAuth.ValidUserReg) (schemesAuth.AnswerUserReg, error) {
	var user *entities.MdUser
	var answer schemesAuth.AnswerUserReg
	//TODO.MD:  нет проверки КАПЧИ
	var phone = &entities.MdPhoneAuth{}
	err := helpers.FillPhoneReg(phone, input)
	if err != nil {
		return answer, &schemes.ErrorResponse{Code: 104, Err: "Bad phone", ErrBase: err}
	}

	phone, err = s.ath.GetUserPhone(phone.Phone)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return answer, &schemes.ErrorResponse{Code: 104, Err: "Bad phone for user", ErrBase: err}
	} else if phone.Verification != false {
		return answer, &schemes.ErrorResponse{Code: 109, Err: "The user already exists", ErrBase: err}
	}

	user, err = s.usr.CreateUser()
	phone.UserUuid = user.Uuid

	//TODO.MD:  установить триггер на удаление старого юзера phone.UserUuid
	verifySms, err := s.ath.SavePhone(phone)
	if err != nil {
		return answer, &schemes.ErrorResponse{Code: 109, Err: "Not save phone", ErrBase: err}
	}
	if verifySms == true {
		return answer, &schemes.ErrorResponse{Code: 109, Err: "The phone already exists", ErrBase: err}
	}

	var mySms entities.MdSmsAuth
	mySms.SmsCode = utils.RandStringBytes(6)
	mySms.IDSend, mySms.SmsService, err = sms.Sender(mySms.SmsCode)
	mySms.Phone = phone.Phone
	mySms.UserUuid = user.Uuid
	smsUuid, err := s.ath.SaveSmsAuth(&mySms)
	if err != nil {
		return answer, &schemes.ErrorResponse{Code: 104, Err: "Bad sms", ErrBase: err}
	}
	answer.SmsId = smsUuid.String()

	return answer, nil
}

func (s *ServiceAuth) ValidSmsUserFlow(ctx *gin.Context, input schemesAuth.ValidSms) (schemesAuth.AnswerToken, error) {
	var answer schemesAuth.AnswerToken
	var smsId uuid.UUID
	err := smsId.Scan(input.SmsId)
	if err != nil {
		return answer, &schemes.ErrorResponse{Code: 104, Err: "Bad SmsId", ErrBase: err}

	}
	phoneNum, err := s.ath.SmsValidUser(smsId, input.SmsCode)
	if err != nil {
		if strings.Contains(err.Error(), "no rows") {
			return answer, &schemes.ErrorResponse{Code: 104, Err: "User with sms-code not found", ErrBase: err}
		}
		return answer, err
	}
	var phone = &entities.MdPhoneAuth{}
	phone, err = s.ath.GetUserPhone(phoneNum)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return answer, &schemes.ErrorResponse{Code: 104, Err: "Bad phone for user", ErrBase: err}
	} else if phone.Verification != false {
		return answer, &schemes.ErrorResponse{Code: 109, Err: "The user already exists", ErrBase: err}
	}
	println(phone)
	err = s.ath.SaveVerifyPhone(true, phone.Phone)
	if err != nil {

	}
	accessToken, err := utils.CreateAccessToken(phone.UserUuid, s.config.AccessTokenSecret, s.config.AccessTokenHour)
	if err != nil {
		return answer, &schemes.ErrorResponse{Code: 96, Err: "Couldn't create a token"}
	}

	refreshToken, err := utils.CreateRefreshToken(phone.UserUuid, s.config.RefreshTokenSecret, s.config.AccessTokenHour)
	if err != nil {
		return answer, &schemes.ErrorResponse{Code: 97, Err: "Couldn't create a token"}
	}
	answer.AccessToken = accessToken
	answer.RefreshToken = refreshToken
	return answer, nil
}

func (s *ServiceAuth) LoginUserFlow(ctx *gin.Context, input schemesAuth.ValidUserReg) (schemesAuth.AnswerToken, error) {
	var token schemesAuth.AnswerToken
	//uuid, hashedPassword, err := s.ath.LoginUser(input.Phone)
	//if err != nil {
	//	if strings.Contains(err.Error(), "no rows") {
	//		return token, &schemes.ErrorResponse{Code: 104, Err: "User not found"}
	//	}
	//	return token, err
	//}
	//
	//err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(input.Password))
	//if err != nil {
	//	return token, &schemes.ErrorResponse{Code: 104, Err: "Invalid username or password"}
	//}
	//
	//accessToken, err := utils.CreateAccessToken(uuid, s.settings.AccessTokenSecret, s.settings.AccessTokenHour)
	//if err != nil {
	//	return token, &schemes.ErrorResponse{Code: 96, Err: "Couldn't create a token"}
	//}
	//
	//refreshToken, err := utils.CreateRefreshToken(uuid, s.settings.RefreshTokenSecret, s.settings.AccessTokenHour)
	//if err != nil {
	//	return token, &schemes.ErrorResponse{Code: 97, Err: "Couldn't create a token"}
	//}
	//token.AccessToken = accessToken
	//token.RefreshToken = refreshToken
	return token, nil
}

func (s *ServiceAuth) RefreshTokenUserFlow(ctx *gin.Context, input schemesAuth.ValidRefresh) (schemesAuth.AnswerToken, error) {
	var token schemesAuth.AnswerToken
	//
	//authorized, _ := utils.IsAuthorized(input.RefreshToken, s.settings.RefreshTokenSecret)
	//if !authorized {
	//	return token, &schemes.ErrorResponse{Code: 98, Err: "Not authorized"}
	//
	//}
	//
	//var userID, err = utils.ExtractToken(input.RefreshToken, s.settings.RefreshTokenSecret)
	//if err != nil {
	//	return token, &schemes.ErrorResponse{Code: 98, Err: "Not find User"}
	//}
	//
	//err = s.ath.GetUuidUser(userID)
	//if err != nil {
	//	return token, &schemes.ErrorResponse{Code: 104, Err: "User not found"}
	//}
	//
	//accessToken, err := utils.CreateAccessToken(userID, s.settings.AccessTokenSecret, s.settings.AccessTokenHour)
	//if err != nil {
	//	return token, &schemes.ErrorResponse{Code: 96, Err: "Couldn't create a token"}
	//}
	//
	//refreshToken, err := utils.CreateRefreshToken(userID, s.settings.RefreshTokenSecret, s.settings.AccessTokenHour)
	//if err != nil {
	//	return token, &schemes.ErrorResponse{Code: 97, Err: "Couldn't create a token"}
	//}
	//token.AccessToken = accessToken
	//token.RefreshToken = refreshToken
	return token, nil
}
