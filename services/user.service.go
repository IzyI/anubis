package services

import (
	entytes "anubis/api/elements"
	"anubis/core"
	"anubis/core/helpers"
	"anubis/core/schemes"
	"anubis/tools/utils"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"strings"
)

type ServiceUser struct {
	user entytes.InfUser
	env  core.Env
}

func NewServiceUser(e entytes.InfUser, env core.Env) *ServiceUser {
	return &ServiceUser{user: e, env: env}
}

func sendSms(s string) error {
	//TODO.MD:  написать отправку sms
	fmt.Printf("Send sms %s", s)
	return nil

}
func (s *ServiceUser) RegUser(input entytes.ShmValidUserReg) (entytes.ShmAnswerUserReg, error) {
	var t entytes.MdUser
	var user entytes.ShmAnswerUserReg
	//TODO.MD:  понять как можно сделать защиту от большого количества отправки смс
	t.Phone = input.Phone
	t.PasswordHash, _ = helpers.GeneratePasswordHash(input.Password)
	//TODO.MD:  понять есть ли проверка что такой пользователь уже уществует
	u, err := s.user.EntityCreate(t)
	if err != nil {
		return user, err
	}
	sms := utils.RandStringBytes(6)
	err = s.user.SmsSave(u.Uuid, sms)
	if err != nil {
		return user, err
	}
	err = sendSms(sms)
	if err != nil {
		return user, err
	}
	user.Uuid = u.Uuid
	return user, nil
}

func (s *ServiceUser) ValidSmsUser(input entytes.ShmValidSms) (entytes.ShmAnswerUserReg, error) {
	var user entytes.ShmAnswerUserReg
	err := s.user.SmsValid(input.Uuid, input.Sms)
	if err != nil {
		if strings.Contains(err.Error(), "no rows") {
			return user, &schemes.ShmErrorResponse{Code: 104, Err: "User with sms-code not found"}
		}
		return user, err
	}
	user.Uuid = input.Uuid
	return user, nil
}

func (s *ServiceUser) LoginUser(input entytes.ShmValidUserReg) (entytes.ShmAnswerToken, error) {
	var token entytes.ShmAnswerToken
	uuid, hashedPassword, err := s.user.LoginUser(input.Phone)
	if err != nil {
		if strings.Contains(err.Error(), "no rows") {
			return token, &schemes.ShmErrorResponse{Code: 104, Err: "User not found"}
		}
		return token, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(input.Password))
	if err != nil {
		return token, &schemes.ShmErrorResponse{Code: 104, Err: "Invalid username or password"}
	}

	accessToken, err := utils.CreateAccessToken(uuid, s.env.AccessTokenSecret, s.env.AccessTokenHour)
	if err != nil {
		return token, &schemes.ShmErrorResponse{Code: 96, Err: "Couldn't create a token"}
	}

	refreshToken, err := utils.CreateRefreshToken(uuid, s.env.RefreshTokenSecret, s.env.AccessTokenHour)
	if err != nil {
		return token, &schemes.ShmErrorResponse{Code: 97, Err: "Couldn't create a token"}
	}
	token.AccessToken = accessToken
	token.RefreshToken = refreshToken
	return token, nil
}

func (s *ServiceUser) RefreshTokenUser(input entytes.ShmValidRefresh) (entytes.ShmAnswerToken, error) {
	var token entytes.ShmAnswerToken

	authorized, _ := utils.IsAuthorized(input.RefreshToken, s.env.RefreshTokenSecret)
	if !authorized {
		return token, &schemes.ShmErrorResponse{Code: 98, Err: "Not authorized"}

	}

	var userID, err = utils.ExtractToken(input.RefreshToken, s.env.RefreshTokenSecret)
	if err != nil {
		return token, &schemes.ShmErrorResponse{Code: 98, Err: "Not find User"}
	}

	err = s.user.GetUuidUser(userID)
	if err != nil {
		return token, &schemes.ShmErrorResponse{Code: 104, Err: "User not found"}
	}

	accessToken, err := utils.CreateAccessToken(userID, s.env.AccessTokenSecret, s.env.AccessTokenHour)
	if err != nil {
		return token, &schemes.ShmErrorResponse{Code: 96, Err: "Couldn't create a token"}
	}

	refreshToken, err := utils.CreateRefreshToken(userID, s.env.RefreshTokenSecret, s.env.AccessTokenHour)
	if err != nil {
		return token, &schemes.ShmErrorResponse{Code: 97, Err: "Couldn't create a token"}
	}
	token.AccessToken = accessToken
	token.RefreshToken = refreshToken
	return token, nil
}
