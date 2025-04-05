package common

import (
	"github.com/go-playground/validator/v10"
	"log"
	"regexp"
	"strconv"
	"unicode"
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

func ValidateObjectId(fl validator.FieldLevel) bool {
	re := regexp.MustCompile(`^[a-f\d]{24}$`)
	return re.MatchString(fl.Field().String())
}

func ValidateSafeText(fl validator.FieldLevel) bool {
	unsafeChars := map[rune]bool{
		'<':  true,
		'>':  true,
		'&':  true,
		'"':  true,
		'\'': true,
	}

	field := fl.Field().String()
	for _, r := range field {
		if unsafeChars[r] || !unicode.IsPrint(r) {
			return false
		}
	}
	return true
}
