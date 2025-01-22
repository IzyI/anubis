package entytes

import "anubis/app/api/schemes"

//	type InfUserDB interface {
//		CreateUser(input MdUser) (*MdUser, error)
//		SmsSaveUser(userUuid string, sms string) error
//		SmsValidUser(userUuid string, sms string) error
//		LoginUser(phone string) (string, string, error)
//		GetUuidUser(uuid string) error
//	}
type InfUserFlow interface {
	RegUserFlow(input schemes.ValidUserReg) (schemes.AnswerUserReg, error)
	ValidSmsUserFlow(input schemes.ValidSms) (schemes.AnswerUserReg, error)
	LoginUserFlow(input schemes.ValidUserReg) (schemes.AnswerToken, error)
	RefreshTokenUserFlow(input schemes.ValidRefresh) (schemes.AnswerToken, error)
}
