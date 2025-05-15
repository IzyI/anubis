package DTO

type EmailUserRegValid struct {
	Email    string `json:"email"  binding:"required,email"`
	Password string `json:"password"  binding:"required,min=8,max=64,containsany=!@#$%^&*"`
	Captcha  string `json:"captcha"  binding:"required"`
	Nickname string `json:"nickname"  binding:"omitempty,safe_text,min=3,max=64"`
}

type AnswerUserRegCode struct {
	CodeId string `json:"code_id"`
}

type CodeEmailValid struct {
	CodeId    string `json:"code_id"  binding:"required,min=23,max=25"`
	DeviceId  string `json:"device_id"  binding:"required,min=12,max=64"`
	EmailCode string `json:"email_code" binding:"required"`
}

type LoginEmailUserValid struct {
	Email    string `json:"email"  binding:"required,email"`
	Password string `json:"password"  binding:"required,min=8,max=64,containsany=!@#$%^&*"`
	DeviceId string `json:"device_id"  binding:"required,min=12,max=64"`
	Captcha  string `json:"captcha"  binding:"required"`
	Nickname string `json:"nickname"  binding:"omitempty,safe_text,min=3,max=64"`
}
