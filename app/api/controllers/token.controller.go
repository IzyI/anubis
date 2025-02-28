package controllers

import (
	"anubis/app/api/api/interfaces"
	"anubis/app/core/handlers"
	"github.com/gin-gonic/gin"
)

type ControllerToken struct {
	tokenUC interfaces.TokenUseCase
}

func NewControllerToken(tokenUC interfaces.TokenUseCase) *ControllerToken {
	return &ControllerToken{tokenUC: tokenUC}
}

func (c *ControllerToken) HandlerRefreshTokenDomainFlowPOST(ctx *gin.Context) {
	handlers.PostHandler(ctx, c.tokenUC.RefreshTokenDomainFlow)
}

func (c *ControllerToken) HandlerLogoutDomainFlowPOST(ctx *gin.Context) {
	handlers.PostHandler(ctx, c.tokenUC.LogoutDomainFlow)
}
