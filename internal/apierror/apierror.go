package apierror

import "net/http"

// ApiError represents an error response from the API
type ApiError struct {
	Message string `json:"message"`
	code    int
}

func (a *ApiError) Error() string {
	return a.Message
}

func (a *ApiError) Code() int {
	return a.code
}

func BadRequest(message string) error {
	return &ApiError{
		Message: message,
		code:    http.StatusBadRequest,
	}
}

func NotFound(message string) error {
	return &ApiError{
		Message: message,
		code:    http.StatusNotFound,
	}
}

func InternalServerError(message string) error {
	return &ApiError{
		Message: message,
		code:    http.StatusInternalServerError,
	}
}

//TODO: Implement any other useful function for creating apierror
