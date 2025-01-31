package schemes

type ValidUserReg struct {
	Phone    string `json:"phone"  binding:"required,phone=8,e164"`
	Password string `json:"password"  binding:"required,min=8,max=64,containsany=!@#$%^&*"`
	Captcha  string `json:"captcha"  binding:"required"`
}

type AnswerUserReg struct {
	SmsId string `json:"sms_id"`
}
type ValidSms struct {
	SmsId   string `json:"sms_id" binding:"omitempty,uuid"`
	SmsCode string `json:"sms_code"`
}

type ValidRefresh struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
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
