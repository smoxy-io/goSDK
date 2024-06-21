package errors

import (
	"bytes"
	"errors"
	"fmt"
	"runtime"
)

var (
	MaxStackDepth = 50
)

type Error struct {
	msg    string
	vars   []any
	pcs    []uintptr
	frames []StackFrame
	stack  []byte
}

func (e *Error) Error() string {
	if e.vars == nil {
		return e.msg
	}

	return fmt.Sprintf(e.msg, e.vars...)
}

func (e *Error) WithVars(vars ...any) *Error {
	return &Error{msg: e.msg, vars: vars, pcs: e.pcs, frames: e.frames, stack: e.stack}
}

func (e *Error) WithStack() *Error {
	return &Error{msg: e.msg, vars: e.vars, pcs: Callers(3), frames: e.frames, stack: e.stack}
}

func (e *Error) Stack() []byte {
	if e.stack != nil {
		return e.stack
	}

	buf := bytes.Buffer{}

	for _, frame := range e.StackFrames() {
		buf.WriteString(frame.String())
	}

	e.stack = buf.Bytes()

	return e.stack
}

func (e *Error) StackFrames() []StackFrame {
	if e.frames != nil {
		return e.frames
	}

	e.frames = make([]StackFrame, len(e.pcs))

	for i, pc := range e.pcs {
		e.frames[i] = NewStackFrame(pc)
	}

	return e.frames
}

// Callers implements the bugsnag ErrorWithCallerS() interface
func (e *Error) Callers() []uintptr {
	return e.pcs
}

// ErrorStack returns a string that contains both the error message and the callstack.
func (e *Error) ErrorStack() string {
	return e.Error() + "\n" + string(e.Stack())
}

func (e *Error) SameAs(err error) bool {
	var errString string
	var et *Error

	switch {
	case errors.As(err, &et):
		errString = et.WithVars().Error()
	default:
		errString = err.Error()
	}

	return e.WithVars().Error() == errString
}

func New(msg string, vars ...any) *Error {
	return &Error{msg: msg, vars: vars, stack: nil}
}

func Callers(skip int) []uintptr {
	stack := make([]uintptr, MaxStackDepth)
	length := runtime.Callers(skip, stack[:])

	return stack[:length]
}
