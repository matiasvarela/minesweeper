package apperrors

import "github.com/matiasvarela/errors"

var (
	NotFound     = errors.Define("not_found")
	Validation   = errors.Define("validation")
	InvalidInput = errors.Define("invalid_input")
	Internal     = errors.Define("internal")
)

type ApiError struct {
	Status  int         `json:"status"`
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func NewApiError(status int, code string, message string, data interface{}) ApiError {
	return ApiError{Status: status, Code: code, Message: message, Data: data}
}

func ToApiError(err error) ApiError {
	switch errors.Code(err) {
	case "not_found":
		return NewApiError(404, errors.Code(err), err.Error(), errors.Data(err))
	case "validation":
		return NewApiError(400, errors.Code(err), err.Error(), errors.Data(err))
	case "invalid_input":
		return NewApiError(400, errors.Code(err), err.Error(), errors.Data(err))
	default:
		return NewApiError(500, errors.Code(err), err.Error(), errors.Data(err))
	}
}
