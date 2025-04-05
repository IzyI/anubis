package interfaces

import (
	"anubis/app/DTO"
	"anubis/app/core/schemes"
	"github.com/gin-gonic/gin"
)

type TokenUseCase interface {
	RefreshTokenDomainFlow(ctx *gin.Context, input *DTO.RefreshTokenProjectValid) (DTO.AnswerToken, error)
	LogoutDomainFlow(ctx *gin.Context, input *DTO.LogoutValid) (schemes.EmptyResponses, error)
}
