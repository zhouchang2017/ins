package insa

import "fmt"

// Error503 is instagram API error
type Error503 struct {
	Message string
}

func (e Error503) Error() string {
	return e.Message
}

// ErrorN is general instagram error
type ErrorN struct {
	Message   string `json:"message"`
	Status    string `json:"status"`
	ErrorType string `json:"error_type"`
}

func (e ErrorN) Error() string {
	return fmt.Sprintf("%s: %s (%s)", e.Status, e.Message, e.ErrorType)
}

// Error400 is error returned by HTTP 400 status code.
type Error400 struct {
	Action     string `json:"action"`
	StatusCode string `json:"status_code"`
	Payload    struct {
		ClientContext string `json:"client_context"`
		Message       string `json:"message"`
	} `json:"payload"`
	Status string `json:"status"`
}

func (e Error400) Error() string {
	return fmt.Sprintf("%s: %s", e.Status, e.Payload.Message)
}