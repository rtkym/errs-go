package errs

import (
	"errors"
)

func haveStackTrace(cause error) bool {
	return getErrorRecursive(cause) != nil
}

func getErrorRecursive(cause error) *Error {
	var errs *Error
	if errors.As(cause, &errs) {
		return errs
	}

	return nil
}

func getStackRecursive(cause error) StackTrace {
	for {
		err := getErrorRecursive(cause)
		if err == nil {
			return nil
		}

		if st := err.stack; st != nil {
			return st
		}

		cause = err.Unwrap()
	}
}
