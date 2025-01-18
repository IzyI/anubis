package schemes

type ValidUserReg struct {
	Phone    string `json:"phone"  binding:"required,phone=8,e164"`
	Password string `json:"password"  binding:"required"`
	Captcha  string `json:"captcha"  binding:"required"`
}

type AnswerUserReg struct {
	Uuid string `json:"uuid"`
}

type ValidRefresh struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}
type AnswerToken struct {
	RefreshToken string `json:"refresh_token"`
	AccessToken  string `json:"access_token"`
}

type ValidSms struct {
	Uuid string `json:"uuid" binding:"omitempty,uuid"`
	Sms  string `json:"sms"`
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
