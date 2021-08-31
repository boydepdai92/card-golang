package common

import (
	"errors"
	"net/http"
)

var ErrRecordNotFound = errors.New("record not found")

type AppError struct {
	RootErr    error  `json:"-"`
	StatusCode int    `json:"status_code"`
	Code       int    `json:"code"`
	Message    string `json:"message"`
}

func (e *AppError) RootError() error {
	if err, ok := e.RootErr.(*AppError); ok {
		return err.RootError()
	}

	return e.RootErr
}

func (e *AppError) Error() string {
	return e.RootErr.Error()
}

func NewErrorResponse(rootErr error, statusCode int, code int, message string) *AppError {
	return &AppError{StatusCode: statusCode, Code: code, Message: message, RootErr: rootErr}
}

func NewFailResponse(err error) *AppError {
	return NewErrorResponse(err, http.StatusServiceUnavailable, CodeFail, GetMessageFromCode(CodeFail))
}

func NewBadRequestResponse(err error, code int, message string) *AppError {
	return NewErrorResponse(err, http.StatusBadRequest, code, message)
}

func NewUnauthorizedResponse() *AppError {
	return NewErrorResponse(nil, http.StatusUnauthorized, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
}

func NewForbiddenResponse() *AppError {
	return NewErrorResponse(nil, http.StatusForbidden, http.StatusForbidden, http.StatusText(http.StatusForbidden))
}

func NewNotFoundResponse() *AppError {
	return NewErrorResponse(nil, http.StatusNotFound, http.StatusNotFound, http.StatusText(http.StatusNotFound))
}

func NewMethodNotAllowedResponse() *AppError {
	return NewErrorResponse(nil, http.StatusMethodNotAllowed, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
}

func NewNoContentResponse() *AppError {
	return NewErrorResponse(nil, http.StatusNoContent, http.StatusNoContent, http.StatusText(http.StatusNoContent))
}
