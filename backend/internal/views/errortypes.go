package views

import "fmt"

type ErrorTitle string

const (
	NotFound   ErrorTitle = "NotFound"
	Unexpected ErrorTitle = "Unexpected"
	BadRequest ErrorTitle = "BadRequest"
)

type InvalidParam struct {
	Field  string `json:"field"`
	Reason string `json:"reason"`
}

type ApiError struct {
	StatusCode    int32          `json:"-"`
	Title         ErrorTitle     `json:"title"`
	Detail        string         `json:"detail"`
	InvalidParams []InvalidParam `json:"invalidParams,omitempty"`
}

func (err ApiError) Error() string {
	return fmt.Sprintf("code: %d, title: %s, detail: %s, invalidParams: %v",
		err.StatusCode, err.Title, err.Detail, err.InvalidParams)
}
