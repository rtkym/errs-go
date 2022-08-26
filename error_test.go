package errs_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/rtkym/errs-go"
)

func TestNew(t *testing.T) {
	msg := "test message"
	err := errs.New(msg)
	assert.Equal(t, err.Error(), msg)
	assert.NotEqual(t, len(err.StackTrace()), 0)
}

func TestWrap(t *testing.T) {
	msg1 := "m1"
	msg2 := "m2"
	msg3 := "m3"

	t.Run("LantanaError", func(t *testing.T) {
		err := errs.Wrap(errs.New("test"))
		assert.Equal(t, "test", err.Error())
		assert.NotEqual(t, len(err.StackTrace()), 0)
	})
	t.Run("LantanaError with message", func(t *testing.T) {
		err := errs.Wrap(errs.New("test"), msg1, msg2, msg3)
		assert.Equal(t, fmt.Sprintf("%s: %s: %s: test", msg1, msg2, msg3), err.Error())
		assert.NotEqual(t, len(err.StackTrace()), 0)
	})
	t.Run("OtherError", func(t *testing.T) {
		err := errs.Wrap(errors.New("test"))
		assert.ErrorContains(t, err, "test")
		assert.NotEqual(t, len(err.StackTrace()), 0)
	})
	t.Run("OtherError with message", func(t *testing.T) {
		err := errs.Wrap(errors.New("test"), msg1, msg2, msg3)
		assert.Equal(t, fmt.Sprintf("%s: %s: %s: test", msg1, msg2, msg3), err.Error())
		assert.NotEqual(t, len(err.StackTrace()), 0)
	})
	t.Run("nil", func(t *testing.T) {
		assert.Panics(t, func() {
			errs.Wrap(nil) //nolint
		})
	})
}

type dummy struct{}

func (dummy) Error() string {
	return "dummy"
}

func TestValues(t *testing.T) {
	msg := "test message"

	t.Run("SingleError", func(t *testing.T) {
		err := errs.New(msg).With("1", "a").With("2", "b").With("3", "c")
		expect := map[string]interface{}{"1": "a", "2": "b", "3": "c"}
		assert.Equal(t, len(err.Values()), len(expect))
		for k, v := range err.Values() {
			if v2, ok := expect[k]; ok {
				assert.Equal(t, v, v2)
			} else {
				assert.Fail(t, "Not found key: %s", k)
			}
		}
	})

	t.Run("NestedError", func(t *testing.T) {
		err := errs.New(msg).With("1", "a").With("2", "b").With("3", "c")
		err = errs.Wrap(err).With("2", "B").With("@", "Z")
		expect := map[string]interface{}{"1": "a", "2": "B", "3": "c", "@": "Z"}
		assert.Equal(t, len(err.Values()), len(expect))
		for k, v := range err.Values() {
			if v2, ok := expect[k]; ok {
				assert.Equal(t, v, v2)
			} else {
				assert.Fail(t, "Not found key: %s", k)
			}
		}
	})

	t.Run("NestedError2", func(t *testing.T) {
		var exx error = &dummy{}
		err := errs.Wrap(exx).With("1", "a").With("2", "b").With("3", "c")
		exx = fmt.Errorf("other error: %w", err)
		err = errs.Wrap(exx).With("2", "B").With("@", "Z")
		expect := map[string]interface{}{"1": "a", "2": "B", "3": "c", "@": "Z"}
		assert.Equal(t, len(err.Values()), len(expect))
		for k, v := range err.Values() {
			if v2, ok := expect[k]; ok {
				assert.Equal(t, v, v2)
			} else {
				assert.Fail(t, "Not found key: %s", k)
			}
		}
	})
}
