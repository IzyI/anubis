package helpers

import (
	"anubis/app/api/DAL/entitiesDB"
	schemesAuth "anubis/app/api/schemes"
	"anubis/app/core"
	"anubis/app/core/schemes"
	"anubis/tools/utils"
	"github.com/nyaruka/phonenumbers"
	"strconv"
)

func FillPhoneReg(phone *entitiesDB.MdPhoneAuth, input schemesAuth.PhoneValidUserReg) error {
	number, err := strconv.ParseInt(input.Phone, 10, 64)
	if err != nil {
		return err
	}
	num, err := phonenumbers.Parse(input.Phone, "")
	if err != nil {
		return err
	}
	phone.Phone = number
	phone.CountryCode = num.GetCountryCode()
	phone.PasswordHash, _ = utils.GeneratePasswordHash(input.Password)
	phone.Verification = false
	return nil
}

func FillJWTTokens(
	answer *schemesAuth.AnswerToken,
	uuid string,
	group []string,
	config core.ServiceConfig,
	aTokenM int,
	rTokenM int,
) error {
	if aTokenM == 0 {
		aTokenM = config.AccessTokenMinute
	}
	if rTokenM == 0 {
		rTokenM = config.RefreshTokenMinute
	}
	println(aTokenM, rTokenM, config.AccessTokenMinute, config.RefreshTokenMinute)
	accessToken, err := utils.CreateAccessToken(uuid, group, config.AccessTokenSecret, aTokenM)
	if err != nil {
		return &schemes.ErrorResponse{Code: 96, Err: "Couldn't create a token", ErrBase: err}
	}

	refreshToken, err := utils.CreateRefreshToken(uuid, []string{}, config.RefreshTokenSecret, rTokenM)
	if err != nil {
		return &schemes.ErrorResponse{Code: 97, Err: "Couldn't create a token", ErrBase: err}
	}
	if config.ShortJwt {
		answer.AccessToken = utils.RemoveFirstPart(accessToken)
		answer.RefreshToken = utils.RemoveFirstPart(refreshToken)
	} else {
		answer.AccessToken = accessToken
		answer.RefreshToken = refreshToken
	}
	return nil
}

func CheckDomain(s core.ServiceConfig, d string) (string, error) {
	domain, ok := s.ListServices[d]
	if !ok {
		return "", &schemes.ErrorResponse{Code: 104, Err: "Domain not found", ErrBase: nil}
	}
	if domain.Auth != nil {
		if !utils.LittleContainsString(domain.Auth, "phone") {
			return "", &schemes.ErrorResponse{Code: 106, Err: "Authorization method denied", ErrBase: nil}
		}
	} else {
		return "", &schemes.ErrorResponse{Code: 106, Err: "Authorization method denied !", ErrBase: nil}
	}
	return domain.Service, nil
}
