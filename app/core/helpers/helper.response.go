package helpers

import (
	schemes2 "anubis/app/core/schemes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgconn"
	"log"
	"net/http"
	"reflect"
	"strconv"
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
func getJSONFieldName(v interface{}, fieldName string) string {
	t := reflect.TypeOf(v)
	field, found := t.FieldByName(fieldName)
	if !found {
		return fieldName // Возвращаем имя поля как есть, если не нашли
	}
	jsonTag := field.Tag.Get("json")
	if jsonTag == "" {
		return fieldName
	}
	return jsonTag
}

func ValidateErrorResponse(ctx *gin.Context, body interface{}, Error error) {

	var ve validator.ValidationErrors
	var jsonErr *json.UnmarshalTypeError
	if errors.As(Error, &ve) {
		var out []schemes2.ValidateError
		out = make([]schemes2.ValidateError, len(ve))

		for i, fe := range ve {
			jsonFieldName := getJSONFieldName(body, fe.StructField())
			out[i] = schemes2.ValidateError{Field: jsonFieldName, Msg: msgForTag(fe.Tag(), fe.Param())}
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
		err := schemes2.ErrorResponse{
			Code: 101,
			Err:  "Bad json",
		}
		ctx.AbortWithStatusJSON(http.StatusBadRequest, err)
		return
	}

}

func ErrorResponse(ctx *gin.Context, code int, error string) {
	err := schemes2.ErrorResponse{
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
		fmt.Printf("strange error: %s", err)
		ErrorResponse(ctx, 99, "Very strange error. please write to the administrator.")
	}
}
