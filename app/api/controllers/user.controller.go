package controllers

import (
	entytes "anubis/app/api/entities"
	"anubis/app/core/handlers"
	"github.com/gin-gonic/gin"
)

type ControllerAuth struct {
	authUC entytes.AuthUseCase
}

func NewControllerAuth(authUC entytes.AuthUseCase) *ControllerAuth {
	return &ControllerAuth{authUC: authUC}
}

func (s *ControllerAuth) HandlerRegPOST(ctx *gin.Context) {
	handlers.PostHandler(ctx, s.authUC.RegUserFlow)
}

func (s *ControllerAuth) HandlerValidSmsPOST(ctx *gin.Context) {
	handlers.PostHandler(ctx, s.authUC.ValidSmsUserFlow)
}

func (s *ControllerAuth) HandlerLoginPOST(ctx *gin.Context) {
	handlers.PostHandler(ctx, s.authUC.LoginUserFlow)
}

func (s *ControllerAuth) HandlerRefreshTokenPOST(ctx *gin.Context) {
	handlers.PostHandler(ctx, s.authUC.RefreshTokenUserFlow)
}
