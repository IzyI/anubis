package core

import (
	"anubis/app/core/helpers"
	"github.com/gin-gonic/gin"
	"net/http"
)

func PostHandler[T any, R any](ctx *gin.Context, processFunc func(T) (R, error)) {
	var body T
	if err := ctx.ShouldBindJSON(&body); err != nil {
		helpers.ValidateErrorResponse(ctx, body, err)
		return
	}

	result, err := processFunc(body)
	if err != nil {
		helpers.HandlerError(ctx, err)
		return
	}

	helpers.APIResponse(ctx, "ok", http.StatusOK, result)
}
