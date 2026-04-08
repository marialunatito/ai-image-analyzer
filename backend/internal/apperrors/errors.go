package apperrors

import "errors"

type Code string

const (
	CodeInvalidRequest  Code = "INVALID_REQUEST"
	CodeInvalidImage    Code = "INVALID_IMAGE"
	CodePayloadLarge    Code = "PAYLOAD_TOO_LARGE"
	CodeProviderError   Code = "PROVIDER_ERROR"
	CodeProviderTimeout Code = "PROVIDER_TIMEOUT"
	CodeInternal        Code = "INTERNAL_ERROR"
)

type AppError struct {
	Code    Code
	Message string
	Err     error
}

func (e *AppError) Error() string {
	if e == nil {
		return ""
	}
	return e.Message
}

func (e *AppError) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Err
}

func New(code Code, message string) error {
	return &AppError{Code: code, Message: message}
}

func Wrap(code Code, message string, err error) error {
	return &AppError{Code: code, Message: message, Err: err}
}

func CodeOf(err error) Code {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.Code
	}
	return CodeInternal
}
