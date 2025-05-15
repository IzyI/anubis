package DTO

type PhoneUserRegValid struct {
	Phone    string `json:"phone"  binding:"required,phone=8,e164"`
	Password string `json:"password"  binding:"required,min=8,max=64,containsany=!@#$%^&*"`
	Captcha  string `json:"captcha"  binding:"required"`
	Nickname string `json:"nickname"  binding:"omitempty,safe_text,min=3,max=64"`
}

type LoginPhoneUserValid struct {
	Phone    string `json:"phone"  binding:"required,phone=8,e164"`
	Password string `json:"password"  binding:"required,min=8,max=64,containsany=!@#$%^&*"`
	DeviceId string `json:"device_id"  binding:"required,min=12,max=64"`
	Captcha  string `json:"captcha"  binding:"required"`
	Nickname string `json:"nickname"  binding:"omitempty,safe_text,min=3,max=64"`
}

type AnswerUserRegSms struct {
	SmsId string `json:"sms_id"`
}
type SmsValid struct {
	SmsId    string `json:"sms_id"  binding:"required,min=23,max=25"`
	DeviceId string `json:"device_id"  binding:"required,min=12,max=64"`
	SmsCode  string `json:"sms_code" binding:"required"`
}

type AnswerRegToken struct {
	RefreshToken string            `json:"refresh_token"`
	AccessToken  string            `json:"access_token"`
	ListProjects map[string]string `json:"list_projects"`
}

type RefreshTokenProjectValid struct {
	ProjectID string `json:"project_id"  binding:"required,object_id"`
}

type LogoutValid struct {
	All bool `json:"all"`
}

type AnswerToken struct {
	RefreshToken string `json:"refresh_token"`
	AccessToken  string `json:"access_token"`
}
