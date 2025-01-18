package controllers

import (
	entytes "anubis/app/api/entytes"
	"anubis/app/core"
	"github.com/gin-gonic/gin"
)

type ControllerUser struct {
	user entytes.InfUserFlow
}

func NewControllerAuth(user entytes.InfUserFlow) *ControllerUser {
	return &ControllerUser{user: user}
}

func (s *ControllerUser) HandlerRegPOST(ctx *gin.Context) {
	core.PostHandler(ctx, s.user.RegUserFlow)
}

func (s *ControllerUser) HandlerValidSmsPOST(ctx *gin.Context) {
	core.PostHandler(ctx, s.user.ValidSmsUserFlow)
}

func (s *ControllerUser) HandlerLoginPOST(ctx *gin.Context) {
	core.PostHandler(ctx, s.user.LoginUserFlow)
}

func (s *ControllerUser) HandlerRefreshTokenPOST(ctx *gin.Context) {
	core.PostHandler(ctx, s.user.RefreshTokenUserFlow)
}
