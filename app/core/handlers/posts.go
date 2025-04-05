package handlers

import (
	"anubis/app/core/common"
	"github.com/gin-gonic/gin"
	"net/http"
)

//func BaseHandler[CTX *gin.Context, T any, R any](
//	ctx *gin.Context,
//	bindFunc func(ctx *gin.Context, req *T) error,
//	processFunc func(CTX, *T) (R, error)) {
//	var body T
//
//	// Привязка данных
//	if err := bindFunc(ctx, &body); err != nil {
//		common.ValidateErrorResponse(ctx, body, err)
//		return
//	}
//
//	// Обработка логики
//	res, err := processFunc(ctx, &body)
//	if err != nil {
//		common.HandlerError(ctx, err)
//		return
//	}
//
//	// Успешный ответ
//	common.APIResponse(ctx, http.StatusOK, res)
//}

// PostHandler - обработчик POST запросов
//func PostHandler[CTX *gin.Context, T any, R any](
//	ctx *gin.Context,
//	processFunc func(CTX, *T) (R, error),
//) {
//	bindFunc := func(c *gin.Context, req *T) error {
//		// Если структура поддерживает биндинг URI, пробуем
//		if err := c.ShouldBindUri(req); err != nil {
//			// Игнорируем ошибку, если в структуре нет полей для URI
//			var ve gin.Error
//			if errors.As(err, &ve) {
//				// Проверяем, есть ли вообще URI-поля в структуре
//				if reflect.ValueOf(req).Elem().NumField() > 0 {
//					return err
//				}
//			}
//		}
//		if err := c.ShouldBindJSON(req); err != nil {
//			return err
//		}
//		return nil
//	}
//
//	BaseHandler(ctx, bindFunc, processFunc)
//}

//func GetHandler[CTX *gin.Context, T any, R any](
//	ctx *gin.Context,
//	processFunc func(CTX, *T) (R, error),
//) {
//	bindFunc := func(c *gin.Context, req *T) error {
//		if err := c.ShouldBindUri(req); err != nil {
//			return err
//		}
//		if err := c.ShouldBindQuery(req); err != nil {
//			return err
//		}
//		return nil
//	}
//
//	BaseHandler(ctx, bindFunc, processFunc)
//}
//
//func DeleteHandler[CTX *gin.Context, T any, R any](
//	ctx *gin.Context,
//	processFunc func(CTX, *T) (R, error),
//) {
//	bindFunc := func(c *gin.Context, req *T) error {
//
//		if err := c.ShouldBindUri(req); err != nil {
//			return err
//		}
//		if err := c.ShouldBindQuery(req); err != nil {
//			return err
//		}
//		return nil
//	}
//
//	BaseHandler(ctx, bindFunc, processFunc)
//}

func UriJsonHandler[CTX *gin.Context, U any, B any, R any](
	ctx *gin.Context,
	processFunc func(CTX, *U, *B) (R, error),
) {
	var uri U
	var body B

	if err := ctx.ShouldBindUri(&uri); err != nil {
		common.ValidateErrorResponse(ctx, uri, err)
		return
	}
	if err := ctx.ShouldBindJSON(&body); err != nil {
		common.ValidateErrorResponse(ctx, body, err)
		return
	}

	// Обработка логики
	res, err := processFunc(ctx, &uri, &body)
	if err != nil {
		common.HandlerError(ctx, err)
		return
	}

	// Успешный ответ
	common.APIResponse(ctx, http.StatusOK, res)
}

func UriHandler[CTX *gin.Context, U any, R any](
	ctx *gin.Context,
	processFunc func(CTX, *U) (R, error),
) {
	var uri U

	if err := ctx.ShouldBindUri(&uri); err != nil {
		common.ValidateErrorResponse(ctx, uri, err)
		return
	}

	// Обработка логики
	res, err := processFunc(ctx, &uri)
	if err != nil {
		common.HandlerError(ctx, err)
		return
	}

	// Успешный ответ
	common.APIResponse(ctx, http.StatusOK, res)
}

func JsonHandler[CTX *gin.Context, B any, R any](
	ctx *gin.Context,
	processFunc func(CTX, *B) (R, error),
) {
	var body B

	if err := ctx.ShouldBindJSON(&body); err != nil {
		common.ValidateErrorResponse(ctx, body, err)
		return
	}

	// Обработка логики
	res, err := processFunc(ctx, &body)
	if err != nil {
		common.HandlerError(ctx, err)
		return
	}

	// Успешный ответ
	common.APIResponse(ctx, http.StatusOK, res)
}
