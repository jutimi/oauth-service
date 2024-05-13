package errors

import (
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"

	logger "gin-boilerplate/package/log"
	_validator "gin-boilerplate/package/validator"
)

type CustomError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

const (
	LangVN = "vi"
)

const (
	ErrCodeRequired        = 1
	ErrCodeValidatorFormat = 2

	ErrCodeInternalServerError = 500
	ErrCodeTimeout             = 408
)

var messages = map[int]map[string]string{
	ErrCodeRequired: {
		LangVN: "không được bỏ trống. Vui lòng kiểm tra lại",
	},
	ErrCodeValidatorFormat: {
		LangVN: "không hợp lệ. Vui lòng kiểm tra lại",
	},
	ErrCodeInternalServerError: {
		LangVN: "Hệ thống gặp lỗi. Vui lòng thử lại sau",
	},
	ErrCodeTimeout: {
		LangVN: "Hệ thống gặp lỗi. Vui lòng thử lại sau",
	},
}

func New(code int) *CustomError {
	return &CustomError{
		Code:    code,
		Message: messages[code][LangVN],
	}
}

func NewCustomError(code int, message string) *CustomError {
	return &CustomError{
		Code:    code,
		Message: message,
	}
}

func NewValidatorError(err error) *CustomError {
	var validatorErr validator.ValidationErrors
	if errors.As(err, &validatorErr) {
		errDetail := validatorErr[0]

		field := errDetail.Field()
		tag := errDetail.Tag()

		code := convertValidatorTag(tag)
		return &CustomError{
			Code:    code,
			Message: fmt.Sprintf("%s %s", field, messages[code][LangVN]),
		}

	}

	return New(ErrCodeInternalServerError)
}

func (err *CustomError) Error() string {
	return err.Message
}

func (err *CustomError) GetCode() int {
	return err.Code
}

// --------------------------------------
func convertValidatorTag(tag string) int {
	logger.GetLogger().Info("validation_tag: ", tag)
	switch tag {
	case _validator.EMAIL, _validator.PHONE_NUMBER:
		return ErrCodeValidatorFormat
	default:
		return ErrCodeRequired
	}
}
