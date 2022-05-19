// Copyright 2022 Ainsley Clark. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// @see https://blog.carlmjohnson.net/post/2020/working-with-errors-as/

package errors

import (
	"bytes"
	"fmt"
	"net/http"
	"runtime"
	"strconv"
	"strings"
)

// Application error codes.
const (
	// CONFLICT - An action cannot be performed.
	CONFLICT = "conflict"
	// INTERNAL - Error within the application.
	INTERNAL = "internal" // Internal error
	// INVALID - Validation failed
	INVALID = "invalid" // Validation failed
	// NOTFOUND - Entity does not exist
	NOTFOUND = "not_found"
	// UNKNOWN - Application unknown error.
	UNKNOWN = "unknown"
	// MAXIMUMATTEMPTS - More than allowed action.
	MAXIMUMATTEMPTS = "maximum_attempts"
	// EXPIRED - Subscription expired.
	EXPIRED = "expired"
)

var (
	// DefaultCode is the default code returned when
	// none is specified.
	DefaultCode = INTERNAL
	// GlobalError is a general message when no error message
	// has been found.
	GlobalError = "An error has occurred."
)

// Error defines a standard application error.
type Error struct {
	Code      string `json:"code"`
	Message   string `json:"message"`
	Operation string `json:"operation"`
	Err       error  `json:"error"`
	fileLine  string
	pcs       []uintptr
}

// Error returns the string representation of the error
// message by implementing the error interface.
func (e *Error) Error() string {
	var buf bytes.Buffer

	// Print the error code if there is one.
	if e.Code != "" {
		buf.WriteString("<" + e.Code + "> ")
	}

	// Print the file-line, if any.
	if e.fileLine != "" {
		buf.WriteString(e.fileLine + " - ")
	}

	// Print the current operation in our stack, if any.
	if e.Operation != "" {
		buf.WriteString(e.Operation + ": ")
	}

	// Print the original error message, if any.
	if e.Err != nil {
		buf.WriteString(e.Err.Error() + ", ")
	}

	// Print the message, if any.
	if e.Message != "" {
		buf.WriteString(e.Message)
	}

	return strings.TrimSuffix(strings.TrimSpace(buf.String()), ",")
}

// New is a wrapper for the stdlib new function.
func New(err error, message, op string) *Error {
	return newError(err, message, DefaultCode, op)
}

// Newf - TODO
func Newf(err error, format string, args ...any) *Error {
	message := fmt.Sprintf(format, args...)
	return New(err, message, "")
}

// Errorf - TODO
func Errorf(err error, format string, args ...any) *Error {
	return Newf(err, format, args...)
}

// FileLine returns the file and line in which the error
// occurred.
func (e *Error) FileLine() string {
	return e.fileLine
}

// Wrap returns an error annotating err with a stack trace
// at the point Wrap is called, and the supplied message.
// If err is nil, Wrap returns nil.
func Wrap(err error, message string) error {
	if err == nil {
		return nil
	}
	return &Error{
		Err:     err,
		Message: message,
	}
}

// HTTPStatusCode is a convenience method used to get the appropriate
// HTTP response status code for the respective error type.
func (e *Error) HTTPStatusCode() int {
	status := http.StatusInternalServerError
	switch e.Code {
	case CONFLICT:
		return http.StatusConflict
	case INVALID:
		return http.StatusBadRequest
	case NOTFOUND:
		return http.StatusNotFound
	case EXPIRED:
		return http.StatusPaymentRequired
	case MAXIMUMATTEMPTS:
		return http.StatusTooManyRequests
	}
	return status
}

// RuntimeFrames returns function/file/line information.
func (e *Error) RuntimeFrames() *runtime.Frames {
	return runtime.CallersFrames(e.pcs)
}

// ProgramCounters returns the slice of PC values associated
// with the error.
func (e *Error) ProgramCounters() []uintptr {
	return e.pcs
}

// StackTrace returns a string representation of the errors
// stacktrace, where each trace is separated by a newline
// and tab '\t'.
func (e *Error) StackTrace() string {
	trace := make([]string, 0, 100)
	rFrames := e.RuntimeFrames()
	frame, ok := rFrames.Next()
	line := strconv.Itoa(frame.Line)
	trace = append(trace, frame.Function+"(): "+e.Message)

	for ok {
		trace = append(trace, "\t"+frame.File+":"+line)
		frame, ok = rFrames.Next()
	}

	return strings.Join(trace, "\n")
}

// StackTraceSlice returns a string slice of the errors
// stacktrace.
func (e *Error) StackTraceSlice() []string {
	trace := make([]string, 0, 100)
	rFrames := e.RuntimeFrames()
	frame, ok := rFrames.Next()
	line := strconv.Itoa(frame.Line)
	trace = append(trace, frame.Function+"(): "+e.Message)

	for ok {
		trace = append(trace, frame.File+":"+line)
		frame, ok = rFrames.Next()
	}

	return trace
}
