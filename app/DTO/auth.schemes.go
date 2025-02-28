package DTO

type PhoneUserRegValid struct {
	Phone    string `json:"phone"  binding:"required,phone=8,e164"`
	Password string `json:"password"  binding:"required,min=8,max=64,containsany=!@#$%^&*"`
	Captcha  string `json:"captcha"  binding:"required"`
	Nickname string `json:"nickname"  binding:"omitempty,min=3,max=64"`
}

type LoginUserValid struct {
	Phone    string `json:"phone"  binding:"required,phone=8,e164"`
	Password string `json:"password"  binding:"required,min=8,max=64,containsany=!@#$%^&*"`
	DeviceId string `json:"device_id"  binding:"required,min=12,max=64"`
	Captcha  string `json:"captcha"  binding:"required"`
	Nickname string `json:"nickname"  binding:"omitempty,min=3,max=64"`
}

type AnswerUserReg struct {
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
	ProjectID string `json:"project_id"  binding:"required,min=23,max=25"`
}

type LogoutValid struct {
	All bool `json:"all"`
}

type AnswerToken struct {
	RefreshToken string `json:"refresh_token"`
	AccessToken  string `json:"access_token"`
}

type IdDevice struct {
	Name     string `form:"name"`
	IdDevice string `form:"id_device"   binding:"required"`
}

type AnswerProjectList struct {
	ListProjects map[string]string `json:"list_projects"`
}

type MembersProjectIdValid struct {
	ID string `uri:"id"  binding:"required,min=23,max=25"`
}

type MembersAddIdValid struct {
	ID     string `uri:"id"  binding:"required,min=23,max=25"`
	UserID string `json:"user_id"  binding:"required,min=23,max=25"`
	Role   string `json:"role"  binding:"omitempty,max=25"`
}

type AnswerMembers struct {
	UserID   string `json:"user_id"`
	Nickname string `json:"nickname"`
	Role     string `json:"role"`
	JoinedAt string `json:"joined_at"`
}
type AnswerProjectID struct {
	ID      string          `json:"id"`
	Name    string          `json:"name"`
	Domain  string          `json:"domain"`
	Members []AnswerMembers `json:"string"`
}
