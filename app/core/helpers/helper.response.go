package helpers

import (
	schemes2 "anubis/app/core/schemes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgconn"
	"net/http"
	"strings"
)

func APIResponse(ctx *gin.Context, message string, StatusCode int, data interface{}) {
	jsonResponse := schemes2.Responses{
		StatusCode: StatusCode,
		Message:    message,
		Data:       data,
	}
	ctx.JSON(StatusCode, jsonResponse)
}

func ErrorResponse(ctx *gin.Context, code int, error string) {
	err := schemes2.HTTPError{
		Code: code,
		Err:  error,
	}
	ctx.AbortWithStatusJSON(http.StatusBadRequest, err)
}

func HandlerError(ctx *gin.Context, err error) {
	var _ = ctx.Error(err)
	var pgErr *pgconn.PgError
	var errResp *schemes2.ErrorResponse

	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23505":
			ErrorResponse(ctx, 103, strings.ToUpper(pgErr.TableName)+" is already exists")
		default:
			ErrorResponse(ctx, 102, "Error dataBase")
		}
	} else if errors.As(err, &errResp) {
		ErrorResponse(ctx, errResp.Code, errResp.Err)

	} else {
		ErrorResponse(ctx, 99, "Very strange error. please write to the administrator.")
	}
}
func ValidateErrorResponse(ctx *gin.Context, body interface{}, Error error) {

	var ve validator.ValidationErrors
	var jsonErr *json.UnmarshalTypeError
	if errors.As(Error, &ve) {
		var out []schemes2.ValidateError
		out = make([]schemes2.ValidateError, len(ve))

		for i, fe := range ve {
			jsonFieldName := getJSONFieldName(body, fe.StructField())
			out[i] = schemes2.ValidateError{Field: jsonFieldName, Msg: msgForTag(fe)}
		}

		err := schemes2.ValidateErrorResponse{
			Code:  100,
			Error: out,
		}
		ctx.AbortWithStatusJSON(http.StatusBadRequest, err)
		return
	} else if errors.As(Error, &jsonErr) {
		out := []schemes2.ValidateError{
			{
				Field: jsonErr.Field,
				Msg:   fmt.Sprintf("Error in the '%s': field is expected type '%s'", jsonErr.Field, jsonErr.Type),
			},
		}

		err := schemes2.ValidateErrorResponse{
			Code:  100,
			Error: out,
		}
		ctx.AbortWithStatusJSON(http.StatusBadRequest, err)
		return
	} else {
		err := schemes2.HTTPError{
			Code: 101,
			Err:  "Bad json",
		}
		ctx.AbortWithStatusJSON(http.StatusBadRequest, err)
		return
	}

}
