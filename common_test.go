package errs_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/rtkym/errs-go"
)

func TestCommonHaveStackTrace(t *testing.T) {
	t.Run("errorインスタンスがError", func(t *testing.T) {
		err := errs.New("test error message")
		assert.True(t, errs.ExpHaveStackTrace(err))
	})

	t.Run("ネストされたError", func(t *testing.T) {
		var err error = errs.New("test error message")
		err = fmt.Errorf("fmt wrap: %w", err)
		assert.True(t, errs.ExpHaveStackTrace(err))
	})

	t.Run("WrapされたError", func(t *testing.T) {
		var err error = errs.Wrap(errs.Wrap(errors.New("root cause"), "test error message"), "2")
		err = fmt.Errorf("fmt wrap: %w", err)
		assert.True(t, errs.ExpHaveStackTrace(err))
	})

	t.Run("WrapされたError", func(t *testing.T) {
		err := errors.New("root cause")
		err = fmt.Errorf("fmt wrap: %w", err)
		assert.False(t, errs.ExpHaveStackTrace(err))
	})
}

func TestCommonExpGetStackRecursive(t *testing.T) {
	t.Run("errorインスタンスがError", func(t *testing.T) {
		err := errs.New("test error message")
		assert.NotEmpty(t, errs.ExpGetStackRecursive(err))
	})

	t.Run("ネストされたError", func(t *testing.T) {
		var err error = errs.New("test error message")
		err = fmt.Errorf("fmt wrap: %w", err)
		assert.NotEmpty(t, errs.ExpGetStackRecursive(err))
	})

	t.Run("WrapされたError", func(t *testing.T) {
		var err error = errs.Wrap(errs.Wrap(errors.New("root cause"), "test error message"), "2")
		err = fmt.Errorf("fmt wrap: %w", err)
		assert.NotEmpty(t, errs.ExpGetStackRecursive(err))
	})

	t.Run("WrapされたError", func(t *testing.T) {
		err := errors.New("root cause")
		err = fmt.Errorf("fmt wrap: %w", err)
		assert.Empty(t, errs.ExpGetStackRecursive(err))
	})
}
