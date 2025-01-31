package handlers

import (
	"anubis/app/core/helpers"
	"github.com/gin-gonic/gin"
	"net/http"
)

func PostHandler[CTX *gin.Context, T any, R any](ctx *gin.Context, processFunc func(CTX, T) (R, error)) {
	var body T
	if err := ctx.ShouldBindJSON(&body); err != nil {
		helpers.ValidateErrorResponse(ctx, body, err)
		return
	}
	result, err := processFunc(ctx, body)
	if err != nil {
		helpers.HandlerError(ctx, err)
		return
	}
	helpers.APIResponse(ctx, "ok", http.StatusOK, result)
}
