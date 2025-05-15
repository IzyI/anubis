package controllers

import (
	"anubis/app/api/interfaces"
	"anubis/app/api/usecase"
	"anubis/app/core/handlers"
	"github.com/gin-gonic/gin"
)

type ControllerEmailAuth struct {
	authUC interfaces.EmailPhoneUseCase
}

func NewControllerEmailAuth(authUC *usecase.ServiceEmailAuth) *ControllerEmailAuth {
	return &ControllerEmailAuth{authUC: authUC}
}

func (c *ControllerEmailAuth) HandlerPOSTReg(ctx *gin.Context) {
	handlers.JsonHandler(ctx, c.authUC.RegUserEmailFlow)
}

func (c *ControllerEmailAuth) HandlerPOSTValidEmailCode(ctx *gin.Context) {
	handlers.JsonHandler(ctx, c.authUC.ValidCodeEmailUserFlow)
}

func (c *ControllerEmailAuth) HandlerPOSTEmailLogin(ctx *gin.Context) {
	handlers.JsonHandler(ctx, c.authUC.EmailLoginFlow)
}
