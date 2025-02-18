package interfaces

import (
	"anubis/app/api/DTO"
	"github.com/gin-gonic/gin"
)

type AuthPhoneUseCase interface {
	RegUserFlow(ctx *gin.Context, input DTO.PhoneValidUserReg) (DTO.AnswerUserReg, error)
	ValidSmsUserFlow(ctx *gin.Context, input DTO.ValidSms) (DTO.AnswerRegToken, error)
	PhoneLoginFlow(ctx *gin.Context, input DTO.PhoneValidUserReg) (DTO.AnswerRegToken, error)
}
