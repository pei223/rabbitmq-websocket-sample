package views

func NewDBError(err error) *ApiError {
	return &ApiError{
		StatusCode: 500,
		Title:      Unexpected,
		Detail:     err.Error(),
	}
}

func NewNotFoundErr(detail string) *ApiError {
	return &ApiError{
		StatusCode:    404,
		Title:         NotFound,
		Detail:        detail,
		InvalidParams: []InvalidParam{},
	}
}

func NewUnexpectedErr(detail string) *ApiError {
	return &ApiError{
		StatusCode:    500,
		Title:         Unexpected,
		Detail:        detail,
		InvalidParams: []InvalidParam{},
	}
}

func NewBadRequestErr(detail string, invalidParams []InvalidParam) *ApiError {
	return &ApiError{
		StatusCode:    400,
		Title:         BadRequest,
		Detail:        detail,
		InvalidParams: invalidParams,
	}
}
