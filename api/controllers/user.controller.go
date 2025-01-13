package controllers

import (
	entytes "anubis/api/elements"
	"anubis/core"
	"github.com/gin-gonic/gin"
)

type ControllerUser struct {
	user entytes.InfUserReg
}

func NewControllerUser(user entytes.InfUserReg) *ControllerUser {
	return &ControllerUser{user: user}
}

func (s *ControllerUser) HandlerRegPOST(ctx *gin.Context) {
	core.PostHandler(ctx, s.user.RegUser)
}

func (s *ControllerUser) HandlerValidSmsPOST(ctx *gin.Context) {
	core.PostHandler(ctx, s.user.ValidSmsUser)
}

func (s *ControllerUser) HandlerLoginPOST(ctx *gin.Context) {
	core.PostHandler(ctx, s.user.LoginUser)
}

func (s *ControllerUser) HandlerRefreshTokenPOST(ctx *gin.Context) {
	core.PostHandler(ctx, s.user.RefreshTokenUser)
}
