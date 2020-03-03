package walhallapi

import (
	"errors"
	"fmt"
	"net/http"
)

var ErrNotFound = errors.New("not found")

type HTTPError struct {
	StatusCode int
	Message    string
}

func NewHTTPError(code int, entity string) HTTPError {
	var message string
	switch code {
	case http.StatusForbidden:
		message = fmt.Sprintf("access to resource %s forbidden", entity)
	case http.StatusNotFound:
		message = fmt.Sprintf("resource %s not found", entity)
	default:
		message = fmt.Sprintf("error access resource %s", entity)
	}
	return HTTPError{
		StatusCode: code,
		Message:    message,
	}
}
func (e HTTPError) Error() string {
	return e.Message
}
