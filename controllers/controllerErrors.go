package controllers

import (
	"fmt"
	"net/http"
)

type ControllerErrors interface {
	Error() string
	HumanReadableError() string
	HTTPStatusCode() int
}

type WrappedControllerError struct {
	Operation        string
	HumanReadableMsg string
	StatusCode       int
	Err              error
}

func (e *WrappedControllerError) Error() string {
	if e.Operation == "" {
		return e.Err.Error()
	}

	return fmt.Sprintf("%s: %v", e.Operation, e.Err)
}

func (e *WrappedControllerError) HumanReadableError() string {
	if e.HumanReadableMsg != "" {
		return e.HumanReadableMsg
	}

	if e.Operation != "" {
		return fmt.Sprintf("An unexpected error occurred while %s", e.Operation)
	}

	return "An unexpected error occurred"
}

func (e *WrappedControllerError) HTTPStatusCode() int {
	if e.StatusCode != 0 {
		return e.StatusCode
	}

	return http.StatusInternalServerError
}

func (e *WrappedControllerError) Unwrap() error {
	return e.Err
}

func WrapControllerError(err error, operation, humanReadableMsg string) error {
	if err == nil {
		return nil
	}

	if _, ok := err.(ControllerErrors); ok {
		return err
	}

	return &WrappedControllerError{
		Operation:        operation,
		HumanReadableMsg: humanReadableMsg,
		StatusCode:       http.StatusInternalServerError,
		Err:              err,
	}
}
