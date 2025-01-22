package entytes

import "github.com/google/uuid"

type MdUser struct {
	Uuid      uuid.UUID `json:"uuid"`
	CreatedAt int       `json:"created_at"`
}

type MdUsersSession struct {
	Id          int       `json:"id"`
	Ip          int       `json:"ip"`
	IdSession   string    `json:"id_session"`
	IdDevice    string    `json:"id_device"`
	Type        string    `json:"type"`
	StatusActiv bool      `json:"status_activ"`
	Token       string    `json:"token"`
	UserUuid    uuid.UUID `json:"user_uuid"`
	IdService   int       `json:"id_service"`
}

type MdUsersInfo struct {
	UserUuid uuid.UUID `json:"user_uuid"`
	Email    string    `json:"email"`
	Nickname string    `json:"nickname"`
}

type MdUserSession struct {
	Id        int    `json:"id"`
	IdSession string `json:"id_session"`
	Device    string `json:"device"`
	IP        string `json:"ip"`
	IdUser    int    `json:"id_user"`
}
