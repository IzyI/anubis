package controllers

import (
	"anubis/app/api/interfaces"
	"anubis/app/core/handlers"
	"github.com/gin-gonic/gin"
)

type ControllerAuth struct {
	authUC interfaces.AuthUseCase
}

func NewControllerAuth(authUC interfaces.AuthUseCase) *ControllerAuth {
	return &ControllerAuth{authUC: authUC}
}

func (s *ControllerAuth) HandlerRegPOST(ctx *gin.Context) {
	handlers.PostHandler(ctx, s.authUC.RegUserFlow)
}

func (s *ControllerAuth) HandlerValidSmsPOST(ctx *gin.Context) {
	handlers.PostHandler(ctx, s.authUC.ValidSmsUserFlow)
}

func (s *ControllerAuth) HandlerLoginPOST(ctx *gin.Context) {
	handlers.PostHandler(ctx, s.authUC.PhoneValidUserReg)
}

//
//func (s *ControllerAuth) HandlerRefreshTokenPOST(ctx *gin.Context) {
//	handlers.PostHandler(ctx, s.authUC.RefreshTokenUserFlow)
//}
