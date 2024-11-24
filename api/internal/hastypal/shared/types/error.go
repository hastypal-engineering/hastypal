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

type ParallelApiErrorInterface interface {
	Error() string
	IsDomain() bool
	PresentError() string
}

type ParallelApiError struct {
	Msg      string
	Function string
	File     string
	Values   []string
	Domain   bool
	Err      error
}

func (e ParallelApiError) Error() string {
	if len(e.Values) > 0 {
		return fmt.Sprintf("Error %s in file %s calling function %s with values %s", e.Msg, e.File, e.Function, e.Values)
	}

	return fmt.Sprintf("Error %s in file %s calling function %s", e.Msg, e.File, e.Function)
}

func (e ParallelApiError) IsDomain() bool {
	return e.Domain
}

func (e ParallelApiError) PresentError() string {
	return fmt.Sprintf("%s", e.Msg)
}

func WrapErr(function string, file string, err error) *ParallelApiError {
	return &ParallelApiError{
		Function: function,
		File:     file,
		Err:      err,
	}
}
