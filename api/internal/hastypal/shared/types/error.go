package types

import "fmt"

type ApiErrorInterface interface {
	Error() string
	IsDomain() bool
	PresentError() string
}

type ApiError struct {
	Msg      string
	Function string
	File     string
	Values   []string
	Domain   bool
	Err      error
}

func NewApiError(msg string, function string, file string, values []string, domain bool, err error) error {
	return &ApiError{
		Msg:      msg,
		Function: function,
		File:     file,
		Values:   values,
		Domain:   domain,
		Err:      err,
	}
}

func WrapError(function string, file string, err error) error {
	return &ApiError{
		Function: function,
		File:     file,
		Err:      err,
	}
}

func (e ApiError) Error() string {
	if len(e.Values) > 0 {
		return fmt.Sprintf("Error %s in file %s calling function %s with values %s", e.Msg, e.File, e.Function, e.Values)
	}

	return fmt.Sprintf("Error %s in file %s calling function %s", e.Msg, e.File, e.Function)
}

func (e ApiError) IsDomain() bool {
	return e.Domain
}

func (e ApiError) PresentError() string {
	return fmt.Sprintf("%s", e.Msg)
}
