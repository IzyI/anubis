package helpers

import (
	"anubis/app/api/entities"
	schemesAuth "anubis/app/api/schemes"
	"anubis/tools/utils"
	"github.com/nyaruka/phonenumbers"
	"strconv"
)

func FillPhoneReg(phone *entities.MdPhoneAuth, input schemesAuth.ValidUserReg) error {
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
