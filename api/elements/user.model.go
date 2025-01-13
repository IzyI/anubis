package elements

type MdUser struct {
	Uuid         string `json:"uuid"`
	Phone        string `json:"phone"`
	PasswordHash string `json:"password_hash"`
	CreatedAt    int    `json:"created_at"`
	UpdatedAt    int    `json:"updated_at"`
	Verification bool   `json:"verification"`
	Sms          string `json:"sms"`
}

type MdUserDevice struct {
	Id       int    `json:"id"`
	IdDevice int    `json:"id_device"`
	IP       string `json:"ip"`
	IdUser   int    `json:"id_user"`
}
