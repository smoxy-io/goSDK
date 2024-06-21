package errors

import (
	"fmt"
	"strings"
	"testing"
)

var testErrMsgs = []string{
	"error message",
	"error message: %v",
	"error message: %v %v",
	"error message:\n  %v",
}

var testErrVars = [][]any{
	nil,
	{"foo"},
	{"foo", "bar"},
	{"baz"},
}

func TestNew(t *testing.T) {
	for i, msg := range testErrMsgs {
		err := New(msg)

		if err.Error() != msg {
			t.Errorf("msg: %d > expected: '%s', got: '%s'", i, msg, err.Error())
			return
		}
	}
}

func TestCallers(t *testing.T) {
	pcs := Callers(0)

	if len(pcs) != 5 {
		t.Errorf("expected: %d frames, got: %d frames", 3, len(pcs))
	}

	pcs1 := Callers(1)

	if len(pcs1) != 4 {
		t.Errorf("expected: %d frames, got: %d frames", 2, len(pcs1))
	}
}

func TestError_WithStack(t *testing.T) {
	for i, msg := range testErrMsgs {
		err := New(msg).WithStack()

		if err.pcs == nil {
			t.Errorf("msg: %d > no program counters", i)
			return
		}

		if err.frames != nil {
			t.Errorf("msg: %d > expected: %d frames, got: %d frames", i, 0, len(err.frames))
			return
		}

		if err.stack != nil {
			t.Errorf("msg: %d > stack is rendered", i)
			return
		}
	}
}

func TestError_StackFrames(t *testing.T) {
	for i, msg := range testErrMsgs {
		err := New(msg).WithStack()

		_ = err.StackFrames()

		if err.frames == nil {
			t.Errorf("msg: %d > no stack frames", i)
			return
		}

		if len(err.frames) != len(err.pcs) {
			t.Errorf("msg: %d > expected: %d frames, got: %d frames", i, len(err.pcs), len(err.frames))
			return
		}

		if err.stack != nil {
			t.Errorf("msg: %d > stack is rendered", i)
			return
		}
	}
}

func TestError_Stack(t *testing.T) {
	for i, msg := range testErrMsgs {
		err := New(msg).WithStack()

		_ = err.Stack()

		if err.stack == nil {
			t.Errorf("msg: %d > no stack", i)
			return
		}

		if len(err.stack) < 1 {
			t.Errorf("msg: %d > stack is empty", i)
			return
		}
	}
}

func TestError_ErrorStack(t *testing.T) {
	for i, msg := range testErrMsgs {
		err := New(msg).WithStack()

		out := err.ErrorStack()

		if !strings.HasPrefix(out, msg+"\n") {
			t.Errorf("msg: %d > expected error msg prefix: '%s', got: '%s'", i, msg+"\n", out)
			return
		}
	}
}

func TestError_WithVars(t *testing.T) {
	for i, msg := range testErrMsgs {
		err := testGetErrWithVars(i)

		if err.Error() != fmt.Sprintf(msg, testErrVars[i]) {
			t.Errorf("msg: %d > expected: '%s', got: '%s", i, fmt.Sprintf(msg, testErrVars[i]), err.Error())
			return
		}
	}
}

func TestError_SameAs(t *testing.T) {
	for i, msg := range testErrMsgs {
		err := New(msg)

		err2 := testGetErrWithVars(i)

		if err2.SameAs(err) {
			t.Errorf("msg: %d > expected '%s' to be same as '%s'", i, err2.Error(), err.Error())
			return
		}
	}
}

func testGetErrWithVars(i int) *Error {
	err := New(testErrMsgs[i])

	if testErrVars[i] == nil {
		return err
	}

	return err.WithVars(testErrVars[i])
}
