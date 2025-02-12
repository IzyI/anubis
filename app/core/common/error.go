package common

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"log"
	"reflect"
	"strconv"
)

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

func msgForTag(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "This field is required"
	case "email":
		return "Invalid email"
	case "numeric":
		return "This field must be a number"
	case "phone":
		return fmt.Sprintf("Phone number should be at least %s digits", fe.Param())
	case "e164":
		return "Invalid phone number, it should be in E.164 format"
	case "len":
		return "Invalid len"
	case "startswith":
		return "Invalid startswith"
	case "min":
		return fmt.Sprintf("Minimum length is %s", fe.Param())
	case "max":
		return fmt.Sprintf("Maximum length is %s", fe.Param())
	case "containsany":
		return fmt.Sprintf("Must contain at least one of the following characters: %s", fe.Param())
	default:
		return fmt.Sprintf("Error tag: %s", fe.Tag())
	}
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
