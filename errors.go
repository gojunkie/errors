package errors

import (
	"fmt"

	"github.com/pkg/errors"
)

// CodeUnknown represents unknown cause of error.
const CodeUnknown = "CODE_UNKNOWN"

// ErrorFunc represents a factory function to an error.
type ErrorFunc func() error

type stackTracer interface {
	StackTrace() errors.StackTrace
}

type causer interface {
	Cause() error
}

type errorCoder interface {
	Code() string
}

type contextor interface {
	Context() interface{}
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

	return &wrapError{
		code: code,
		err:  errors.WithStack(err),
	}
}

// WithContext annotates err with a context.
// If err is nill, WithContext return nil.
// An err value has code if it implements the following
// interface:
//
//	type errorCoder interface {
//		Code() string
//	}
//
// If the error does not implement Code, UnknownCode will be added to the error.
func WithContext(err error, ctx interface{}) error {
	if err == nil {
		return nil
	}

	if e, ok := err.(*wrapError); ok {
		e.context = ctx
		return e
	}

	e := &wrapError{
		context: ctx,
		code:    CodeUnknown,
		err:     errors.WithStack(err),
	}
	if c, ok := err.(errorCoder); ok {
		e.code = c.Code()
	}

	return e
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
	code := CodeUnknown

	if err != nil {
		if coder, ok := err.(errorCoder); ok {
			code = coder.Code()
		}
	}

	return code
}

// Context returns the underlying error context, if any.
// An error value has a context if it implements the following
// interface:
//
//	type contextor interface {
//		Code() string
//	}
//
// If the error does not implement Context or the error is nil, nil will be returned.
func Context(err error) interface{} {
	if e, ok := err.(contextor); ok {
		return e.Context()
	}
	return nil
}

type simpleError struct {
	err  error
	code string
}

func (err *simpleError) Error() string {
	return err.err.Error()
}

func (err *simpleError) Code() string {
	return err.code
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

func (err *simpleError) StackTrace() errors.StackTrace {
	if e, ok := err.err.(stackTracer); ok {
		return e.StackTrace()[1:]
	}
	return nil
}

type wrapError struct {
	err     error
	code    string
	context interface{}
}

func (err *wrapError) Error() string {
	return err.err.Error()
}

func (err *wrapError) Cause() error {
	return errors.Cause(err.err)
}

func (err *wrapError) Code() string {
	return err.code
}

func (err *wrapError) Context() interface{} {
	return err.context
}

func (err *wrapError) Format(s fmt.State, c rune) {
	if s.Flag('+') {
		if c == 'v' {
			fmt.Fprintf(s, "title=%s, code=%s", err.Error(), err.code)
			fmt.Fprintf(s, "%+v", err.StackTrace())
			return
		}
	}

	fmt.Fprintf(s, "%s", err.err)
}

func (err *wrapError) StackTrace() errors.StackTrace {
	if e, ok := err.err.(stackTracer); ok {
		return e.StackTrace()[1:]
	}
	return nil
}
