package entytes

import "github.com/google/uuid"

type MdPhoneAuth struct {
	Phone        int       `json:"phone"`
	CountryCode  int       `json:"country_code"`
	PasswordHash string    `json:"password_hash"`
	UpdatedAt    int       `json:"updated_at"`
	CreatedAt    int       `json:"created_at"`
	Verification bool      `json:"verification"`
	UserUuid     uuid.UUID `json:"user_uuid"`
}

type SmsAuth struct {
	UserUuid   uuid.UUID `json:"user_uuid"`
	Phone      int       `json:"phone"`
	SmsCode    string    `json:"sms_code"`
	SmsService string    `json:"sms_service"`
	IdSend     string    `json:"id_send"`
	CreatedAt  int       `json:"created_at"`
}
