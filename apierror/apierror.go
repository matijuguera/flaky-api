package apiError

import (
	"fmt"
)

type APIError struct {
	Status  string
	Code    int
	URL     string
	Message string
}

func (err APIError) Error() string {
	return fmt.Sprintf("%s [%s]: %d - %s", err.Message, err.URL, err.Code, err.Status)
}

func NewAPIError(code int, message, URL, status string) APIError {
	return APIError{
		Status:  status,
		Code:    code,
		Message: message,
		URL:     URL,
	}
}
