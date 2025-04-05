package controllers

import (
	"anubis/app/api/interfaces"
	"anubis/app/core/handlers"
	"github.com/gin-gonic/gin"
)

type ControllerAuth struct {
	authUC interfaces.AuthPhoneUseCase
}

func NewControllerAuth(authUC interfaces.AuthPhoneUseCase) *ControllerAuth {
	return &ControllerAuth{authUC: authUC}
}

func (c *ControllerAuth) HandlerPOSTReg(ctx *gin.Context) {
	handlers.JsonHandler(ctx, c.authUC.RegUserFlow)
}

func (c *ControllerAuth) HandlerPOSTValidSms(ctx *gin.Context) {
	handlers.JsonHandler(ctx, c.authUC.ValidSmsUserFlow)
}

func (c *ControllerAuth) HandlerPOSTLogin(ctx *gin.Context) {
	handlers.JsonHandler(ctx, c.authUC.PhoneLoginFlow)
}
