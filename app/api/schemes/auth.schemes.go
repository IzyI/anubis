package schemes

type PhoneValidUserReg struct {
	Domain   string `json:"domain"  binding:"required,min=3,max=64"`
	Phone    string `json:"phone"  binding:"required,phone=8,e164"`
	Password string `json:"password"  binding:"required,min=8,max=64,containsany=!@#$%^&*"`
	Captcha  string `json:"captcha"  binding:"required"`
	Nickname string `json:"nickname"  binding:"omitempty,min=3,max=64"`
}

type AnswerUserReg struct {
	SmsId string `json:"sms_id"`
}
type ValidSms struct {
	Domain  string `json:"domain"  binding:"required,min=3,max=64"`
	SmsId   string `json:"sms_id"  binding:"required,min=23,max=25"`
	SmsCode string `json:"sms_code" binding:"required"`
}

type ValidRefresh struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}
type RegAnswerToken struct {
	AccessToken string            `json:"access_token"`
	GroupId     map[string]string `json:"group_id"`
}
type AnswerToken struct {
	RefreshToken string `json:"refresh_token"`
	AccessToken  string `json:"access_token"`
}

type UserInfo struct {
	Email string `json:"email"  binding:"required"`
	Phone string `json:"phone"  binding:"required"`
	Name  string `json:"name"`
	FIO   string `json:"fio"`
}

type UserDevice struct {
	IdDevice int    `json:"id_device"   binding:"required"`
	IP       string `json:"ip"   binding:"required"`
	IdUser   int    `json:"id_user"   binding:"required"`
}
