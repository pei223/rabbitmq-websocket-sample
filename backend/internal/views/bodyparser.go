package views

import (
	"encoding/json"
	"io"
	"net/http"
)

func ParseBody[T Validatable](r *http.Request) (*T, error) {
	var body T
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&body); err != nil {
		err := NewBadRequestErr("JSON Parse error: "+err.Error(), []InvalidParam{})
		return nil, err
	}
	return &body, nil
}

func ParseBodyAndValidate[T Validatable](req *http.Request) (*T, error) {
	if req.Header.Get("Content-Type") != "application/json" {
		err := NewBadRequestErr("Content-Type must be application/json", []InvalidParam{})
		return nil, err
	}
	b, err := io.ReadAll(req.Body)
	if err != nil {
		err := NewBadRequestErr("Request Body read error: "+err.Error(), []InvalidParam{})
		return nil, err
	}
	var body T
	if err := json.Unmarshal(b, &body); err != nil {
		err := NewBadRequestErr("JSON Parse error: "+err.Error(), []InvalidParam{})
		return nil, err
	}

	invalidParams := body.Validate()
	if invalidParams != nil {
		err := NewBadRequestErr("Invalid Payload", invalidParams)
		return nil, err
	}
	return &body, nil
}
