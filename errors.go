package mailinabox

import (
	"errors"
	"fmt"
)

var ErrServerError = errors.New("server error")

// APIError represents a Parse API Error response.
type APIError struct {
	Code         int    `json:"code"`
	ErrorMessage string `json:"error"`
}

// Message displays the Code and Error as string message.
func (e *APIError) Error() string {

	if len(e.ErrorMessage) > 0 {
		return fmt.Sprintf("%v - %s", e.Code, e.ErrorMessage)
	}

	return ""

}
