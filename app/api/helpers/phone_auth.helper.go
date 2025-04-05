package helpers

import (
	"anubis/app/DAL/entitiesDB"
	schemesAuth "anubis/app/DTO"
	"anubis/tools/utils"
	"github.com/nyaruka/phonenumbers"
	"strconv"
)

func FillPhoneReg(phone *entitiesDB.MdPhoneAuth, input *schemesAuth.PhoneUserRegValid) error {
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
