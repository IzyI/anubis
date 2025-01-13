package elements

type InfUser interface {
	EntityCreate(input MdUser) (*MdUser, error)
	SmsSave(user_uuid string, sms string) error
	SmsValid(user_uuid string, sms string) error
	LoginUser(phone string) (string, string, error)
	GetUuidUser(uuid string) error
}
type InfUserReg interface {
	RegUser(input ShmValidUserReg) (ShmAnswerUserReg, error)
	ValidSmsUser(input ShmValidSms) (ShmAnswerUserReg, error)
	LoginUser(input ShmValidUserReg) (ShmAnswerToken, error)
	RefreshTokenUser(input ShmValidRefresh) (ShmAnswerToken, error)
}
