package helpers

import (
	"anubis/core/schemes"
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgconn"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func APIResponse(ctx *gin.Context, message string, StatusCode int, data interface{}) {
	jsonResponse := schemes.ShmResponses{
		StatusCode: StatusCode,
		Message:    message,
		Data:       data,
	}
	ctx.JSON(StatusCode, jsonResponse)
}

func ValidatePhone(fl validator.FieldLevel) bool {
	phone := fl.Field().String()
	minLengthStr := fl.Param()
	minLength, err := strconv.Atoi(minLengthStr)
	if err != nil {
		log.Fatal(err)
	}
	if len(phone) < minLength {
		return false
	}
	return true
}

func msgForTag(tag string, param string) string {
	switch tag {
	case "required":
		return "This field is required"
	case "email":
		return "Invalid email"
	case "numeric":
		return "This field must be a number"
	case "phone":
		return fmt.Sprintf("Phone number should be at least %s digits", param)
	case "e164":
		return "Invalid phone number, it should be in E.164 format"
	case "len":
		return "Invalid len"
	case "startswith":
		return "Invalid startswith"
	}
	return ""
}

func ValidateErrorResponse(ctx *gin.Context, Error error) {

	var ve validator.ValidationErrors
	var out []schemes.ValidateError
	if errors.As(Error, &ve) {
		out = make([]schemes.ValidateError, len(ve))
		for i, fe := range ve {
			out[i] = schemes.ValidateError{Field: fe.Field(), Msg: msgForTag(fe.Tag(), fe.Param())}
		}
	} else {
		err := schemes.ShmErrorResponse{
			Code: 101,
			Err:  "Bad json",
		}
		ctx.AbortWithStatusJSON(http.StatusBadRequest, err)
		return
	}

	err := schemes.ShmValidateErrorResponse{
		Code:  100,
		Error: out,
	}
	ctx.AbortWithStatusJSON(http.StatusBadRequest, err)
}

func ErrorResponse(ctx *gin.Context, code int, error string) {
	err := schemes.ShmErrorResponse{
		Code: code,
		Err:  error,
	}
	ctx.AbortWithStatusJSON(http.StatusBadRequest, err)
}

func HandlerError(ctx *gin.Context, err error) {
	var _ = ctx.Error(err)
	var pgErr *pgconn.PgError
	var errResp *schemes.ShmErrorResponse
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
		fmt.Printf("strange error: %s", err)
		ErrorResponse(ctx, 99, "Very strange error. please write to the administrator.")
	}
}
