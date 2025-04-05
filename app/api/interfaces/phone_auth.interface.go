package interfaces

import (
	"anubis/app/DTO"
	"github.com/gin-gonic/gin"
)

type AuthPhoneUseCase interface {
	RegUserFlow(ctx *gin.Context, input *DTO.PhoneUserRegValid) (DTO.AnswerUserReg, error)
	ValidSmsUserFlow(ctx *gin.Context, input *DTO.SmsValid) (DTO.AnswerRegToken, error)
	PhoneLoginFlow(ctx *gin.Context, input *DTO.LoginUserValid) (DTO.AnswerRegToken, error)
}
