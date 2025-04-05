package common

import (
	schemes2 "anubis/app/core/schemes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"reflect"
	"strings"
)

type Empty struct {
}

func APIResponse(ctx *gin.Context, StatusCode int, data interface{}) {
	jsonResponse := schemes2.Responses{
		StatusCode: StatusCode,
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
	return
}

func getJSONTag(body interface{}, fieldName string) string {
	t := reflect.TypeOf(body)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	field, _ := t.FieldByName(fieldName)
	jsonTag := field.Tag.Get("json")
	formTag := field.Tag.Get("form")
	if jsonTag != "" {
		return strings.Split(jsonTag, ",")[0]
	}
	if formTag != "" {
		return strings.Split(formTag, ",")[0]
	}
	return fieldName
}

//
//func validationMessage(tag string, param string) string {
//	messages := map[string]string{
//		"required":    "This field is required",
//		"email":       "Invalid email format",
//		"min":         fmt.Sprintf("Minimum length: %s", param),
//		"max":         fmt.Sprintf("Maximum length: %s", param),
//		"numeric":     "This field must be a number",
//		"phone":       fmt.Sprintf("Phone number should be at least %s digits", param),
//		"e164":        "Invalid phone number, it should be in E.164 format",
//		"len":         "Invalid len",
//		"lowercase":   "Only lowercase",
//		"alpha":       "Only Latin letters are allowed.",
//		"startswith":  "Invalid startswith",
//		"containsany": fmt.Sprintf("Must contain at least one of the following characters: %s", param),
//	}
//	return messages[tag]
//}

func HandlerError(ctx *gin.Context, err error) {

	var _ = ctx.Error(err)
	var mongoErr mongo.WriteException
	var errResp *schemes2.ErrorResponse

	if errors.As(err, &mongoErr) {
		for _, writeErr := range mongoErr.WriteErrors {
			switch writeErr.Code {
			case 11000: // Duplicate key error code
				ErrorResponse(ctx, 103, "Document already exists")
			default:
				ErrorResponse(ctx, 102, "Database error")
			}
		}
	} else if errors.As(err, &errResp) {
		ErrorResponse(ctx, errResp.Code, errResp.Err)

	} else {
		ErrorResponse(ctx, 99, "Very strange error. please write to the administrator.")
	}
}

//func handleMongoError(ctx *gin.Context, mongoErr mongo.WriteException) {
//	for _, writeErr := range mongoErr.WriteErrors {
//		if writeErr.Code == 11000 {
//			ErrorResponse(ctx, 103, "Duplicate entry")
//			return
//		}
//	}
//	ErrorResponse(ctx, 102, "Database operation failed")
//}

func handleJSONErrors(ctx *gin.Context, jsonErr *json.UnmarshalTypeError) {
	ctx.AbortWithStatusJSON(http.StatusBadRequest, schemes2.ValidateErrorResponse{
		Code: 100,
		Error: []schemes2.ValidateError{{
			Field: jsonErr.Field,
			Msg:   fmt.Sprintf("Expected type %s for field %s", jsonErr.Type, jsonErr.Field),
		}},
	})
	return
}

func handleValidationErrors(ctx *gin.Context, body interface{}, ve validator.ValidationErrors) {
	out := make([]schemes2.ValidateError, len(ve))
	for i, fe := range ve {

		jsonField := getJSONTag(body, fe.StructField())
		out[i] = schemes2.ValidateError{
			Field: jsonField,
			Msg:   msgForTag(fe),
		}
	}
	ctx.AbortWithStatusJSON(http.StatusBadRequest, schemes2.ValidateErrorResponse{
		Code:  100,
		Error: out,
	})
	return
}

func ValidateErrorResponse(ctx *gin.Context, body interface{}, err error) {
	var ve validator.ValidationErrors
	var jsonErr *json.UnmarshalTypeError
	switch {
	case errors.As(err, &ve):
		handleValidationErrors(ctx, body, ve)
	case errors.As(err, &jsonErr):
		handleJSONErrors(ctx, jsonErr)
	default:
		//TODO: узнать что за ошибка
		ErrorResponse(ctx, 101, "Invalid request format")
	}
}

//func ValidateErrorResponse(ctx *gin.Context, body interface{}, Error error) {
//	var ve validator.ValidationErrors
//	var jsonErr *json.UnmarshalTypeError
//	if errors.As(Error, &ve) {
//		var out []schemes2.ValidateError
//		out = make([]schemes2.ValidateError, len(ve))
//
//		for i, fe := range ve {
//			jsonFieldName := getJSONFieldName(body, fe.StructField())
//			out[i] = schemes2.ValidateError{Field: jsonFieldName, Msg: msgForTag(fe)}
//		}
//
//		err := schemes2.ValidateErrorResponse{
//			Code:  100,
//			Error: out,
//		}
//		ctx.AbortWithStatusJSON(http.StatusBadRequest, err)
//		return
//	} else if errors.As(Error, &jsonErr) {
//		out := []schemes2.ValidateError{
//			{
//				Field: jsonErr.Field,
//				Msg:   fmt.Sprintf("Error in the '%s': field is expected type '%s'", jsonErr.Field, jsonErr.Type),
//			},
//		}
//
//		err := schemes2.ValidateErrorResponse{
//			Code:  100,
//			Error: out,
//		}
//		ctx.AbortWithStatusJSON(http.StatusBadRequest, err)
//		return
//	} else {
//		err := schemes2.HTTPError{
//			Code: 101,
//			Err:  "Bad json",
//		}
//		ctx.AbortWithStatusJSON(http.StatusBadRequest, err)
//		return
//	}
//
//}
