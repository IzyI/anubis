package interfaces

import (
	"anubis/app/api/schemes"
	"github.com/gin-gonic/gin"
)

type AuthUseCase interface {
	RegUserFlow(ctx *gin.Context, input schemes.PhoneValidUserReg) (schemes.AnswerUserReg, error)
	ValidSmsUserFlow(ctx *gin.Context, input schemes.ValidSms) (schemes.RegAnswerToken, error)
	PhoneValidUserReg(ctx *gin.Context, input schemes.PhoneValidUserReg) (schemes.AnswerToken, error)
	//RefreshTokenUserFlow(ctx *gin.Context, input schemes.ValidRefresh) (schemes.AnswerToken, error)
}
