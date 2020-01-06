package errors

import (
	"fmt"

	"github.com/pkg/errors"
)

// UnknownCode represents unknown cause of error.
const UnknownCode = "UNKNOWN_CODE"

// ErrorFunc represents a factory function to an error.
type ErrorFunc func() error

type simpleError struct {
	err  error
	code string
}

type stackTracer interface {
	StackTrace() errors.StackTrace
}

type causer interface {
	Cause() error
}

type errorCoder interface {
	Code() string
}

// Func creates and return a factory function.
func Func(msg, code string) ErrorFunc {
	return func() error {
		return &simpleError{
			code: code,
			err:  errors.New(msg),
		}
	}
}

// New creates an error from message and code.
func New(msg, code string) error {
	return &simpleError{
		code: code,
		err:  errors.New(msg),
	}
}

// Wrap annotates err with a stack trace at the point Wrap is called.
// If err is nil, Wrap return nil.
func Wrap(err error, code string) error {
	if err == nil {
		return nil
	}

	return &simpleError{
		code: code,
		err:  errors.WithStack(err),
	}
}

// Cause returns the underlying cause of error, if any.
// An error value has a cause if it implements the following
// interface:
//
//	type causer interface {
//		Cause() error
//	}
//
// If the error does not implement Cause, the original error will
// be returned. If the error is nil, nil will be returned.
func Cause(err error) error {
	for err != nil {
		cause, ok := err.(causer)
		if !ok {
			break
		}
		err = cause.Cause()
	}
	return err
}

// Code returns the underlying error code, if any.
// An error value has a code if it implements the following
// interface:
//
//	type errorCoder interface {
//		Code() string
//	}
//
// If the error does not implement Code or the error is nil, UnknownCode will be returned
func Code(err error) string {
	code := UnknownCode

	for err != nil {
		coder, ok := err.(errorCoder)
		if !ok {
			break
		}
		code = coder.Code()
	}

	return code
}

func (err *simpleError) Error() string {
	return err.err.Error()
}

func (err *simpleError) Cause() error {
	return errors.Cause(err.err)
}

func (err *simpleError) Code() string {
	return err.code
}

func (err *simpleError) StackTrace() errors.StackTrace {
	if e, ok := err.err.(stackTracer); ok {
		return e.StackTrace()[1:]
	}
	return nil
}

func (err *simpleError) Format(s fmt.State, c rune) {
	if s.Flag('+') {
		if c == 'v' {
			fmt.Fprintf(s, "title=%s, code=%s", err.Error(), err.code)
			fmt.Fprintf(s, "%+v", err.StackTrace())
			return
		}
	}

	fmt.Fprintf(s, "%s", err.err)
}
