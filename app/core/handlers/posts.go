package handlers

import (
	"anubis/app/core/common"
	"github.com/gin-gonic/gin"
	"net/http"
)

func PostHandler[CTX *gin.Context, T any, R any](ctx *gin.Context, processFunc func(CTX, T) (R, error)) {
	var body T
	if err := ctx.ShouldBindJSON(&body); err != nil {
		common.ValidateErrorResponse(ctx, body, err)
		return
	}
	result, err := processFunc(ctx, body)
	if err != nil {
		common.HandlerError(ctx, err)
		return
	}
	common.APIResponse(ctx, http.StatusOK, result)
}
