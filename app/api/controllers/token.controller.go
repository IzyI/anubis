package controllers

import (
	"anubis/app/api/interfaces"
	"anubis/app/core/handlers"
	"github.com/gin-gonic/gin"
)

type ControllerToken struct {
	tokenUC interfaces.TokenUseCase
}

func NewControllerToken(tokenUC interfaces.TokenUseCase) *ControllerToken {
	return &ControllerToken{tokenUC: tokenUC}
}

func (c *ControllerToken) HandlerPOSTRefreshTokenDomain(ctx *gin.Context) {
	handlers.JsonHandler(ctx, c.tokenUC.RefreshTokenDomainFlow)
}

func (c *ControllerToken) HandlerPOSTLogoutDomain(ctx *gin.Context) {
	handlers.JsonHandler(ctx, c.tokenUC.LogoutDomainFlow)
}
