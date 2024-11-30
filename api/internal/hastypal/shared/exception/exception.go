package exception

import (
	"errors"
	"fmt"
	"strings"
)

type HastypalErrorInterface interface {
	Error() string
	IsDomain() bool
	PresentError() string
}

type HastypalError struct {
	Msg       string
	Function  string
	File      string
	Values    []string
	DomainErr bool
	Err       error
}

func New(msg string) HastypalError {
	return HastypalError{
		Msg: msg,
	}
}

func Wrap(function string, file string, err error) HastypalError {
	return HastypalError{
		Function: function,
		File:     file,
		Err:      err,
	}
}

func (e HastypalError) UnWrap() error {
	return e.Err
}

func (e HastypalError) WithValues(data []string) HastypalError {
	e.Values = append(e.Values, data...)

	return e
}

func (e HastypalError) Trace(function string, file string) HastypalError {
	e.Function = function
	e.File = file

	return e
}

func (e HastypalError) Error() string {
	var formattedErr strings.Builder

	stackTrace := e.getStackTrace(e)

	formattedErr.WriteString(fmt.Sprintf("%s:\n", stackTrace[len(stackTrace)-1][0]))

	for i := len(stackTrace); i >= 1; i-- {
		formattedErr.WriteString(fmt.Sprintf("[%d] %s\n", i, stackTrace[i-1][1]))
	}

	return formattedErr.String()
}

func (e HastypalError) IsDomain() bool {
	return e.DomainErr
}

func (e HastypalError) Domain() HastypalError {
	e.DomainErr = true

	return e
}

func (e HastypalError) PresentError() string {
	return fmt.Sprintf("%s", e.Msg)
}

func (e HastypalError) getStackTrace(head error) [][]string {
	var stackTrace [][]string
	var hastypalErr HastypalError

	currentNode := head

	for currentNode != nil {
		if errors.As(currentNode, &hastypalErr) {
			beautifiedStack := fmt.Sprintf("File: %s, Function: %s, Values: %s",
				hastypalErr.File,
				hastypalErr.Function,
				hastypalErr.Values,
			)

			stackTrace = append(stackTrace, []string{hastypalErr.Msg, beautifiedStack})

			currentNode = hastypalErr.Err

			continue
		}

		break
	}

	return stackTrace
}
