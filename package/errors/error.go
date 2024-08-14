package errors

import (
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"

	logger "oauth-server/package/log"
	_validator "oauth-server/package/validator"
)

type CustomError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

const (
	LangVN = "vi"
)

const (
	ErrCodeValidatorRequired     = 1
	ErrCodeValidatorFormat       = 2
	ErrCodeValidatorVerifiedData = 3
	ErrCodeMethodNotSupported    = 4

	ErrCodeUserNotFound    = 10
	ErrCodeUserExisted     = 11
	ErrCodeUserDeactivated = 12

	ErrCodeRefreshTokenExpired = 20
	ErrCodeIncorrectPassword   = 21
	ErrCodeAccessTokenExpired  = 22

	ErrCodePermissionNotFound       = 30
	ErrCodePermissionActionNotFound = 31

	ErrCodeInternalServerError = 500
	ErrCodeTimeout             = 408
	ErrCodeForbidden           = 403
	ErrCodeUnauthorized        = 401
)

var messages = map[int]map[string]string{
	// Validator Error
	ErrCodeValidatorRequired: {
		LangVN: "không được bỏ trống. Vui lòng kiểm tra lại",
	},
	ErrCodeValidatorFormat: {
		LangVN: "không hợp lệ. Vui lòng kiểm tra lại",
	},
	ErrCodeValidatorVerifiedData: {
		LangVN: "không chính xác. Vui lòng kiểm tra lại",
	},
	ErrCodeMethodNotSupported: {
		LangVN: "Phương thức chưa được hỗ trợ",
	},

	// Basic Error
	ErrCodeInternalServerError: {
		LangVN: "Hệ thống gặp lỗi. Vui lòng thử lại sau",
	},
	ErrCodeTimeout: {
		LangVN: "Hệ thống gặp lỗi. Vui lòng thử lại sau",
	},
	ErrCodeForbidden: {
		LangVN: "Bạn không đủ quyền truy cập tài nguyên này. Vui lòng kiểm tra lại quyền",
	},
	ErrCodeUnauthorized: {
		LangVN: "Bạn không có quyền truy cập tài nguyên này. Vui lòng kiểm tra lại quyền",
	},

	// User Error
	ErrCodeUserNotFound: {
		LangVN: "Không tìm thấy người dùng. Vui lòng kiểm tra lại",
	},
	ErrCodeUserExisted: {
		LangVN: "Người dùng đã đăng ký tài khoản. Vui lòng đăng nhập hoặc sử dụng sđt/email khác",
	},
	ErrCodeUserDeactivated: {
		LangVN: "Tài khoản của bạn đã bị vô hiệu hóa. Vui lòng liên hệ admin để hỗ trợ",
	},

	// OAuth Error
	ErrCodeRefreshTokenExpired: {
		LangVN: "Phiên làm việc đã hết hạn. Vui lòng đăng nhâp lại",
	},
	ErrCodeAccessTokenExpired: {
		LangVN: "Phiên làm việc đã hết hạn.",
	},
	ErrCodeIncorrectPassword: {
		LangVN: "Tên đăng nhập hoặc mật khẩu không đúng. Vui lòng kiểm tra lại",
	},

	// Permission Error
	ErrCodePermissionNotFound: {
		LangVN: "Không tìm thấy quyền. Vui loại kiểm tra lại",
	},
	ErrCodePermissionActionNotFound: {
		LangVN: "Không tìm hành động của quyền. Vui loại kiểm tra lại",
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
	case _validator.EQUAL_FIELD:
		return ErrCodeValidatorVerifiedData
	default:
		return ErrCodeValidatorRequired
	}
}
