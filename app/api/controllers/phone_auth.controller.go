package controllers

import (
	"anubis/app/api/api/interfaces"
	"anubis/app/core/handlers"
	"github.com/gin-gonic/gin"
)

type ControllerAuth struct {
	authUC interfaces.AuthPhoneUseCase
}

func NewControllerAuth(authUC interfaces.AuthPhoneUseCase) *ControllerAuth {
	return &ControllerAuth{authUC: authUC}
}

func (c *ControllerAuth) HandlerRegPOST(ctx *gin.Context) {
	handlers.PostHandler(ctx, c.authUC.RegUserFlow)
}

func (c *ControllerAuth) HandlerValidSmsPOST(ctx *gin.Context) {
	handlers.PostHandler(ctx, c.authUC.ValidSmsUserFlow)
}

func (c *ControllerAuth) HandlerLoginPOST(ctx *gin.Context) {
	handlers.PostHandler(ctx, c.authUC.PhoneLoginFlow)
}
