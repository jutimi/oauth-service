package validator

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

const (
	REQUIRED     = "required"
	EMAIL        = "email"
	PHONE_NUMBER = "phone_number"
	EQUAL_FIELD  = "eqfield"
)

var validatePhoneNumber validator.Func = func(fl validator.FieldLevel) bool {
	phoneNumber, ok := fl.Field().Interface().(string)
	if !ok {
		return false
	}

	// Check if the phone number matches the Vietnamese format
	vietnamesePhoneNumberPattern := `^(03[2-9]|07[0|6-9]|08[1-5]|09[0-9]|01[2|6|8|9])+([0-9]{7})$`
	match, err := regexp.MatchString(vietnamesePhoneNumberPattern, phoneNumber)
	if err != nil {
		return false
	}

	if !match {
		return false
	}

	return true
}

func RegisterCustomValidators(v *validator.Validate) {
	v.RegisterValidation("phone_number", validatePhoneNumber)
}
