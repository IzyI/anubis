package controllers

import (
	"anubis/app/api/interfaces"
	"anubis/app/core/handlers"
	"github.com/gin-gonic/gin"
)

type ControllerPhoneAuth struct {
	authUC interfaces.AuthPhoneUseCase
}

func NewControllerPhoneAuth(authUC interfaces.AuthPhoneUseCase) *ControllerPhoneAuth {
	return &ControllerPhoneAuth{authUC: authUC}
}

func (c *ControllerPhoneAuth) HandlerPOSTReg(ctx *gin.Context) {
	handlers.JsonHandler(ctx, c.authUC.RegUserPhoneFlow)
}

func (c *ControllerPhoneAuth) HandlerPOSTValidSms(ctx *gin.Context) {
	handlers.JsonHandler(ctx, c.authUC.ValidSmsUserFlow)
}

func (c *ControllerPhoneAuth) HandlerPOSTPhoneLogin(ctx *gin.Context) {
	handlers.JsonHandler(ctx, c.authUC.PhoneLoginFlow)
}
