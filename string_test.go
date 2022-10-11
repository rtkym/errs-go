package errs_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/rtkym/errs-go"
)

func TestString(t *testing.T) {
	const ErrSimple errs.StringError = "ErrSimple"

	t.Run("単体でエラーとして機能する", func(t *testing.T) {
		err := ErrSimple

		assert.Error(t, err)
		assert.Equal(t, "ErrSimple", err.Error())

		assert.True(t, errors.Is(err, ErrSimple))

		var asErr errs.StringError
		assert.True(t, errors.As(err, &asErr))
		assert.Equal(t, ErrSimple, asErr)
	})

	t.Run("Newでエラーを生成して機能する", func(t *testing.T) {
		newErr := errNew(ErrSimple)

		assert.Equal(t, "ErrSimple", newErr.Error())
		assert.Equal(t, ErrSimple, errors.Unwrap(newErr))

		assert.True(t, errors.Is(newErr, ErrSimple))

		var asErr errs.StringError
		assert.True(t, errors.As(newErr, &asErr))
		assert.Equal(t, ErrSimple, asErr)
	})

	t.Run("Wrapでエラーを生成して機能する", func(t *testing.T) {
		var (
			Ref1        = ErrSimple
			Ref2  error = ErrSimple
			cause error = errors.New("test")
		)

		wrapErr := errWrap(ErrSimple, cause)

		assert.Equal(t, "ErrSimple: test", wrapErr.Error())
		assert.Equal(t, cause, errors.Unwrap(errors.Unwrap(wrapErr)))

		assert.True(t, errors.Is(wrapErr, ErrSimple))
		assert.True(t, errors.Is(wrapErr, Ref1))
		assert.True(t, errors.Is(wrapErr, &Ref1))
		assert.True(t, errors.Is(wrapErr, Ref2))
		assert.True(t, errors.Is(wrapErr, cause))

		var asErr errs.StringError
		assert.True(t, errors.As(wrapErr, &asErr))
		assert.Equal(t, ErrSimple, asErr)
	})
}

func errNew(e errs.StringError) error {
	return e.New()
}

func errWrap(e errs.StringError, cause error) error {
	return e.Wrap(cause)
}
