package errs

import (
	"errors"
	"fmt"
)

// StringErrorは指定された文字列をエラーメッセージとして扱います.
type StringError string

// Error returns error message
func (s StringError) Error() string {
	return string(s)
}

// New creates *errs.Error that wraps StringError.
func (s StringError) New() *Error {
	err := newError("", false)
	err.cause = s

	return err
}

// Wrap creates *errs.Error that wraps StringError and cause.
func (s StringError) Wrap(cause error) *Error {
	wrap := &stringError{
		strErr: s,
		cause:  cause,
	}

	err := newError("", false)
	err.cause = wrap

	return err
}

type stringError struct {
	strErr StringError
	cause  error
}

// Error returns error message
func (s stringError) Error() string {
	return fmt.Sprintf("%s: %v", s.strErr, s.cause.Error())
}

// Unwrap returns errors that cause
func (s stringError) Unwrap() error {
	return s.cause
}

// Is follows the specification of errors.Is.
func (s stringError) Is(target error) bool {
	switch x := target.(type) { // nolint:errorlint
	case StringError:
		if string(s.strErr) == string(x) {
			return true
		}
	case *StringError:
		if x != nil && string(s.strErr) == string(*x) {
			return true
		}
	}

	return errors.Is(s.cause, target)
}

// As follows the specification of errors.As.
func (s stringError) As(target interface{}) bool {
	if x, ok := target.(*StringError); ok {
		*x = s.strErr
		return true
	}

	return false
}
