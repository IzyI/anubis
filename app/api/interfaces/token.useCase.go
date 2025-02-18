package interfaces

import (
	"anubis/app/api/DTO"
	"anubis/app/core/schemes"
	"github.com/gin-gonic/gin"
)

type TokenUseCase interface {
	RefreshTokenDomainFlow(ctx *gin.Context, input DTO.RefreshTokenProjectI) (DTO.AnswerToken, error)
	LogoutDomainFlow(ctx *gin.Context, input DTO.Logout) (schemes.EmptyResponses, error)
}
