package entities

import (
	"anubis/app/api/schemes"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AuthPhoneRepository interface {
	SavePhone(phone *MdPhoneAuth) (bool, error)
	SaveVerifyPhone(verify bool, phone int64) error
	SaveSmsAuth(sms *MdSmsAuth) (uuid.UUID, error)
	GetUserPhone(phone int64) (*MdPhoneAuth, error)
	SmsValidUser(userUuid uuid.UUID, sms string) (int64, error)
	LoginUser(phone string) (uuid.UUID, uuid.UUID, error)
	GetUuidUser(uuid uuid.UUID) error
}
type UserRepository interface {
	CreateUser() (*MdUser, error)
}
type AuthUseCase interface {
	RegUserFlow(ctx *gin.Context, input schemes.ValidUserReg) (schemes.AnswerUserReg, error)
	ValidSmsUserFlow(ctx *gin.Context, input schemes.ValidSms) (schemes.AnswerToken, error)
	LoginUserFlow(ctx *gin.Context, input schemes.ValidUserReg) (schemes.AnswerToken, error)
	RefreshTokenUserFlow(ctx *gin.Context, input schemes.ValidRefresh) (schemes.AnswerToken, error)
}
