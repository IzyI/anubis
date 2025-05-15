package interfaces

import (
	"anubis/app/DTO"
	"github.com/gin-gonic/gin"
)

type EmailPhoneUseCase interface {
	RegUserEmailFlow(ctx *gin.Context, input *DTO.EmailUserRegValid) (DTO.AnswerUserRegCode, error)
	ValidCodeEmailUserFlow(ctx *gin.Context, input *DTO.CodeEmailValid) (DTO.AnswerRegToken, error)
	EmailLoginFlow(ctx *gin.Context, input *DTO.LoginEmailUserValid) (DTO.AnswerRegToken, error)
}
