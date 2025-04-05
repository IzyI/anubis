package common

import (
	"fmt"
	"github.com/go-playground/validator/v10"
)

func msgForTag(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "This field is required"
	case "email":
		return "Invalid email"
	case "numeric":
		return "This field must be a number"
	case "safe_text":
		return "Contains the forbidden character <>&\\ “ '"
	case "object_id":
		return "Bad 24-character id"
	case "phone":
		return fmt.Sprintf("Phone number should be at least %s digits", fe.Param())
	case "e164":
		return "Invalid phone number, it should be in E.164 format"
	case "len":
		return "Invalid len"
	case "lowercase":
		return "Only lowercase"
	case "alpha":
		return "Only Latin letters are allowed."
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
