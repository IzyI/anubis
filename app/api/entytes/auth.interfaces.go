package entytes

import "anubis/app/api/schemes"

type InfAuthPhoneDB interface {
	SavePhone(phone *MdPhoneAuth) error
	SmsSaveUser(sms SmsAuth) error
	GetUserPhone(phone int) (*MdPhoneAuth, error)
	SmsValidUser(userUuid string, sms string) error
	LoginUser(phone string) (string, string, error)
	GetUuidUser(uuid string) error
}
type InfUserDB interface {
	CreateUser() (*MdUser, error)
	//SmsSaveUser(userUuid string, sms string) error
	//GetUserPhone(phone int) (*MdPhoneAuth, error)
	//SmsValidUser(userUuid string, sms string) error
	//LoginUser(phone string) (string, string, error)
	//GetUuidUser(uuid string) error
}
type InfAuthFlow interface {
	RegUserFlow(input schemes.ValidUserReg) (schemes.AnswerUserReg, error)
	ValidSmsUserFlow(input schemes.ValidSms) (schemes.AnswerUserReg, error)
	LoginUserFlow(input schemes.ValidUserReg) (schemes.AnswerToken, error)
	RefreshTokenUserFlow(input schemes.ValidRefresh) (schemes.AnswerToken, error)
}
