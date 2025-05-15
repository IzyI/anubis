package interfaces

import (
	"anubis/app/DTO"
	"github.com/gin-gonic/gin"
)

type AuthPhoneUseCase interface {
	RegUserPhoneFlow(ctx *gin.Context, input *DTO.PhoneUserRegValid) (DTO.AnswerUserRegSms, error)
	ValidSmsUserFlow(ctx *gin.Context, input *DTO.SmsValid) (DTO.AnswerRegToken, error)
	PhoneLoginFlow(ctx *gin.Context, input *DTO.LoginPhoneUserValid) (DTO.AnswerRegToken, error)
}
