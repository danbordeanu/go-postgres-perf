package utils

import (
	"fmt"
	"github.com/go-errors/errors"
)


type HttpError struct {
	Code    int
	Err     error
	Message string
}

func (err *HttpError) Error() string {
	if err.Message != "" {
		return fmt.Sprintf("%s: %s", err.Message, err.Err.Error())
	}
	return err.Err.Error()
}

func (err *HttpError) StackTrace() string {
	return errors.Wrap(err.Err, 1).ErrorStack()
}

func (err *HttpError) ErrorAndStack() (string, string) {
	return err.Error(), err.StackTrace()
}
