package helpers

import (
	"anubis/app/api/DAL/entitiesDB"
	schemesAuth "anubis/app/api/DTO"
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
