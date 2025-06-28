package service

import (
	"errors"
	"strings"
)

type PublicMessageError struct {
	Err     error
	Message string
}

func WrapPublicMessage(err error, msg string) error {
	if err == nil {
		return nil
	}
	return &PublicMessageError{Err: err, Message: сapitalizeFirst(msg)}
}

func (e *PublicMessageError) Error() string {
	return e.Err.Error()
}

func (e *PublicMessageError) Unwrap() error {
	return e.Err
}

func isPublicMessageError(err error) bool {
	var publicErr *PublicMessageError
	return errors.As(err, &publicErr)
}

func сapitalizeFirst(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}
