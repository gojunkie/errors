package errors

import (
	"errors"
	"fmt"
	"io"
	"os"
	"testing"
)

func TestFunc(t *testing.T) {
	type testStruct struct {
		fn   ErrorFunc
		msg  string
		code string
	}

	suites := []testStruct{
		testStruct{
			fn:  Func("msg1", "code1"),
			msg: "msg1",
		},
		testStruct{
			fn:  Func("msg2", "code2"),
			msg: "msg2",
		},
	}

	for _, suite := range suites {
		err := suite.fn()
		if err.Error() != suite.msg {
			t.Errorf("Func(): got = %q; want %q", err.Error(), suite.msg)
		}
	}
}

func TestNew(t *testing.T) {
	suites := []struct {
		err  string
		want error
	}{
		{"", fmt.Errorf("")},
		{"error1", fmt.Errorf("error1")},
		{"error2", New("error2", "code1")},
	}

	for _, suite := range suites {
		got := New(suite.err, "code1")
		if got.Error() != suite.want.Error() {
			t.Errorf("New.Error(): got %q, want %q", got, suite.want)
		}
	}
}

func TestCause(t *testing.T) {
	type testStruct struct {
		cause error
		err   error
	}

	err1 := errors.New("err1")
	err2 := New("err2", "code2")

	suites := []testStruct{
		testStruct{
			err:   Wrap(err1, "code1"),
			cause: err1,
		},
		testStruct{
			err:   Wrap(err2, "code3"),
			cause: err2,
		},
	}

	for _, suite := range suites {
		cause := Cause(suite.err)
		if cause != suite.cause {
			t.Errorf("Cause(%v): got %v; want %v", suite.err, cause, suite.cause)
		}
	}
}

func TestWrap(t *testing.T) {
	suites := []struct {
		err  error
		want string
	}{
		{io.EOF, "EOF"},
		{os.ErrClosed, "file already closed"},
	}
	for _, suite := range suites {
		got := Wrap(suite.err, "code1")
		if got.Error() != suite.want {
			t.Errorf(`Wrap(%v, "code1"): got %q, want %q`, suite.err, got, suite.want)
		}
	}
}

func TestCode(t *testing.T) {
	suites := []struct {
		err  error
		want string
	}{
		{fmt.Errorf("err1"), CodeUnknown},
		{New("err2", "code2"), "code2"},
		{nil, CodeUnknown},
	}
	for _, suite := range suites {
		got := Code(suite.err)
		if got != suite.want {
			t.Errorf("Code(%v): got %s, want %s", suite.err, got, suite.want)
		}
	}
}
