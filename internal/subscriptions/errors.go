package subscriptions

import (
	"errors"
	"fmt"
)

type ValidationError struct {
	Errors []error
}

func (err *ValidationError) Error() string {
	message := "invalid request data: "
	for _, e := range err.Errors {
		message = message + e.Error() + " "
	}
	return message
}

func (err *ValidationError) Unwrap() error {
	return errors.Join(err.Errors...)
}

type ResourceNotFoundError struct {
	Err error
}

func (err *ResourceNotFoundError) Error() string {
	return fmt.Sprintf("recource not found %v", err.Err)
}
func (err *ResourceNotFoundError) Unwrap() error {
	return err.Err
}

type DatabaseError struct {
	Query string
	Err   error
}

func (err *DatabaseError) Error() string {
	return fmt.Sprintf("error during the query execution: %s %v", err.Query, err.Err)
}
func (err *DatabaseError) Unwrap() error {
	return err.Err
}

type JsonError struct {
	Json string
	Err  error
}

func (err *JsonError) Error() string {
	return fmt.Sprintf("error during json coding/encoding: %s %v", err.Json, err.Err)
}
func (err *JsonError) Unwrap() error {
	return err.Err
}

type MethodNotAllowedError struct {
	RequiredMethod string
}

func (err *MethodNotAllowedError) Error() string {
	return fmt.Sprintf("Only %s method is allowed!", err.RequiredMethod)
}
func (err *MethodNotAllowedError) Unwrap() error {
	return nil
}

type InvalidParameterError struct {
	ParamName string
}

func (err *InvalidParameterError) Error() string {
	return fmt.Sprintf("invalid parameter: %s", err.ParamName)
}
func (err *InvalidParameterError) Unwrap() error {
	return nil
}
