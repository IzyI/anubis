package services

import (
	"anubis/app/api/entytes"
	schemesAuth "anubis/app/api/schemes"
	"anubis/app/api/storage"
	"anubis/app/core"
	"anubis/app/core/schemes"
	"anubis/tools/utils"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"strconv"
	"strings"
)

type ServiceAuth struct {
	ath    entytes.InfAuthPhoneDB
	usr    entytes.InfUserDB
	config core.ServiceConfig
}

func NewServiceAuth(usr *storage.RepositoryPsqlUser, ath *storage.RepositoryPsqlAuth, config core.ServiceConfig) *ServiceAuth {
	return &ServiceAuth{ath, usr, config}
}

func sendSms(s string) (string, string) {
	//TODO.MD:  написать отправку sms
	//TODO.MD:  написать  обработку что потдерживаем покачто только +7 (россию)
	fmt.Printf("Send sms %s \n", s)
	return s + "IdSend", "my_sms_SmsService"

}
func (s *ServiceAuth) RegUserFlow(input schemesAuth.ValidUserReg) (schemesAuth.AnswerUserReg, error) {
	var user *entytes.MdUser
	var phone *entytes.MdPhoneAuth
	var answer schemesAuth.AnswerUserReg
	//TODO.MD:  понять как можно сделать защиту от большого количества отправки смс
	number, err := strconv.Atoi(input.Phone)
	if err != nil {
		return answer, &schemes.ErrorResponse{Code: 104, Err: "Bad phone"}
	}

	phone, err = s.ath.GetUserPhone(number)
	if err != nil {
		return answer, &schemes.ErrorResponse{Code: 104, Err: "Ифв2"}
	}
	phone.Phone = number
	phone.CountryCode = 2122
	phone.PasswordHash, _ = utils.GeneratePasswordHash(input.Password)
	phone.Verification = false

	user, err = s.usr.CreateUser()
	phone.UserUuid = user.Uuid
	err = s.ath.SavePhone(phone)
	if err != nil {
		return answer, err
	}
	var my_sms entytes.SmsAuth
	my_sms.SmsCode = utils.RandStringBytes(6)
	my_sms.IdSend, my_sms.SmsService = sendSms(my_sms.SmsCode)
	my_sms.Phone = number
	my_sms.UserUuid = user.Uuid
	err = s.ath.SmsSaveUser(my_sms)
	if err != nil {
		return answer, err
	}

	return answer, nil
}

func (s *ServiceAuth) ValidSmsUserFlow(input schemesAuth.ValidSms) (schemesAuth.AnswerUserReg, error) {
	var user schemesAuth.AnswerUserReg
	err := s.ath.SmsValidUser(input.Uuid, input.Sms)
	if err != nil {
		if strings.Contains(err.Error(), "no rows") {
			return user, &schemes.ErrorResponse{Code: 104, Err: "User with sms-code not found"}
		}
		return user, err
	}
	user.Uuid = input.Uuid
	return user, nil
}

func (s *ServiceAuth) LoginUserFlow(input schemesAuth.ValidUserReg) (schemesAuth.AnswerToken, error) {
	var token schemesAuth.AnswerToken
	uuid, hashedPassword, err := s.ath.LoginUser(input.Phone)
	if err != nil {
		if strings.Contains(err.Error(), "no rows") {
			return token, &schemes.ErrorResponse{Code: 104, Err: "User not found"}
		}
		return token, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(input.Password))
	if err != nil {
		return token, &schemes.ErrorResponse{Code: 104, Err: "Invalid username or password"}
	}

	accessToken, err := utils.CreateAccessToken(uuid, s.config.AccessTokenSecret, s.config.AccessTokenHour)
	if err != nil {
		return token, &schemes.ErrorResponse{Code: 96, Err: "Couldn't create a token"}
	}

	refreshToken, err := utils.CreateRefreshToken(uuid, s.config.RefreshTokenSecret, s.config.AccessTokenHour)
	if err != nil {
		return token, &schemes.ErrorResponse{Code: 97, Err: "Couldn't create a token"}
	}
	token.AccessToken = accessToken
	token.RefreshToken = refreshToken
	return token, nil
}

func (s *ServiceAuth) RefreshTokenUserFlow(input schemesAuth.ValidRefresh) (schemesAuth.AnswerToken, error) {
	var token schemesAuth.AnswerToken

	authorized, _ := utils.IsAuthorized(input.RefreshToken, s.config.RefreshTokenSecret)
	if !authorized {
		return token, &schemes.ErrorResponse{Code: 98, Err: "Not authorized"}

	}

	var userID, err = utils.ExtractToken(input.RefreshToken, s.config.RefreshTokenSecret)
	if err != nil {
		return token, &schemes.ErrorResponse{Code: 98, Err: "Not find User"}
	}

	err = s.ath.GetUuidUser(userID)
	if err != nil {
		return token, &schemes.ErrorResponse{Code: 104, Err: "User not found"}
	}

	accessToken, err := utils.CreateAccessToken(userID, s.config.AccessTokenSecret, s.config.AccessTokenHour)
	if err != nil {
		return token, &schemes.ErrorResponse{Code: 96, Err: "Couldn't create a token"}
	}

	refreshToken, err := utils.CreateRefreshToken(userID, s.config.RefreshTokenSecret, s.config.AccessTokenHour)
	if err != nil {
		return token, &schemes.ErrorResponse{Code: 97, Err: "Couldn't create a token"}
	}
	token.AccessToken = accessToken
	token.RefreshToken = refreshToken
	return token, nil
}
