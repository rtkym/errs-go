package errs

import (
	"fmt"
	"io"
	"strings"

	"github.com/goccy/go-json"
)

// New creates a new *Error.
func New(msg string) *Error {
	return newError(msg, false)
}

// Wrap creates a new *Error. However, if the cause is *Error, cause is returned as is.
func Wrap(cause error, messages ...interface{}) *Error {
	if cause == nil {
		panic("nil error")
	}

	haveStackTrace := haveStackTrace(cause)

	if len(messages) == 0 && haveStackTrace {
		if e, ok := cause.(*Error); ok { //nolint:errorlint
			return e
		}
	}

	msg := "wrap"

	if len(messages) != 0 {
		newMsgs := make([]string, 0, len(messages))
		for _, m := range messages {
			newMsgs = append(newMsgs, fmt.Sprintf("%v", m))
		}

		msg = strings.Join(newMsgs, ": ")
	}

	err := newError(msg, haveStackTrace)
	err.cause = cause

	return err
}

// Error is an error interface with associated attributes.
type Error struct {
	msg    string
	cause  error
	values map[string]interface{}
	stack  StackTrace
}

var (
	_ error            = (*Error)(nil)
	_ fmt.Formatter    = (*Error)(nil)
	_ json.Marshaler   = (*Error)(nil)
	_ json.Unmarshaler = (*Error)(nil)
)

func newError(msg string, haveStackTrace bool) *Error {
	var stack StackTrace
	if !haveStackTrace {
		stack = getStackTrace()
	}

	return &Error{
		msg:    msg,
		cause:  nil,
		values: make(map[string]interface{}),
		stack:  stack,
	}
}

// Error returns error message
func (x *Error) Error() string {
	msg := x.msg

	if cause := x.cause; cause != nil {
		if len(msg) == 0 {
			msg = fmt.Sprintf("%v", cause.Error())
		} else {
			msg = fmt.Sprintf("%s: %v", msg, cause.Error())
		}
	}

	return msg
}

// Unwrap returns cause
func (x *Error) Unwrap() error {
	return x.cause
}

// With adds key and value to error attribute.
func (x *Error) With(key string, value interface{}) *Error {
	x.values[key] = value

	return x
}

// Values ​​returns a map of attributes set by With.
// Error attributes are overridden with the error attributes of the wrapped errs.Error.
func (x *Error) Values() map[string]interface{} {
	var values map[string]interface{}

	for cause := x.Unwrap(); cause != nil; {
		if err, ok := cause.(*Error); ok { //nolint:errorlint
			values = err.Values()
			break
		} else {
			type errorUnwrap interface {
				Unwrap() error
			}
			unwrapable, ok := cause.(errorUnwrap) //nolint:errorlint
			if !ok {
				break
			}
			cause = unwrapable.Unwrap()
		}
	}

	if values == nil {
		values = make(map[string]interface{})
	}

	for key, value := range x.values {
		values[key] = value
	}

	return values
}

// StackTrace returns stack trace
func (x *Error) StackTrace() StackTrace {
	st := getStackRecursive(x)
	if st == nil {
		return make([]StackFrame, 0)
	}

	return st
}

// Format returns:
// - %v, %s, %q: formated message
// - %+v: formated message with values and stack trace
func (x *Error) Format(state fmt.State, verb rune) {
	switch verb {
	case 'v':
		_, _ = io.WriteString(state, x.Error())

		if state.Flag('+') {
			x.formatStack(state, verb)
		}

	case 's':
		_, _ = io.WriteString(state, x.Error())

	case 'q':
		fmt.Fprintf(state, "%q", x.Error())
	}
}

func (x *Error) formatStack(state fmt.State, verb rune) {
	if b, err := json.Marshal(x.Values()); err != nil {
		fmt.Fprintf(state, "\nvalues=%s", err.Error())
	} else {
		fmt.Fprintf(state, "\nvalues=%s", string(b))
	}

	if st := x.StackTrace(); st != nil {
		st.Format(state, verb)
	}
}

// MarshallJSON returns a JSON string of errs.Error.
// However, because the Cause is embedded in the Message, the type information is lost.
func (x *Error) MarshalJSON() ([]byte, error) {
	v, err := json.Marshal(&struct {
		Message    string                 `json:"message"`
		Values     map[string]interface{} `json:"values"`
		StackTrace StackTrace             `json:"stackTrace"`
	}{
		Message:    x.Error(),
		Values:     x.Values(),
		StackTrace: x.StackTrace(),
	})
	if err != nil {
		return nil, Wrap(err)
	}

	return v, nil
}

// UnmarshalJSON recovers Error objects from JSON strings.
// However, nested error structures are not restored.
func (x *Error) UnmarshalJSON(b []byte) error {
	resource := new(struct {
		Message    string                 `json:"message"`
		Values     map[string]interface{} `json:"values"`
		StackTrace StackTrace             `json:"stackTrace"`
	})
	if err := json.Unmarshal(b, resource); err != nil {
		return err
	}

	x.msg = resource.Message
	x.values = resource.Values
	x.stack = resource.StackTrace

	return nil
}
