// Copyright 2022 Ainsley Clark. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// @see https://blog.carlmjohnson.net/post/2020/working-with-errors-as/

package errors

import (
	"bytes"
	"encoding/json"
	"errors"
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
	INTERNAL = "internal"
	// INVALID - Validation failed.
	INVALID = "invalid"
	// NOTFOUND - Entity does not exist.
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
	// The application error code.
	Code string `json:"code" bson:"code"`
	// A human-readable message to send back to the end user.
	Message string `json:"message" bson:"message"`
	// Defines what operation is currently being run.
	Operation string `json:"operation" bson:"op"`
	// The error that was returned from the caller.
	Err      error `json:"error" bson:"err"`
	fileLine string
	pcs      []uintptr
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

// NewE returns an Error with the DefaultCode.
func NewE(err error, message, op string) *Error {
	return newError(err, message, DefaultCode, op)
}

// ErrorF returns an Error with the DefaultCode and
// formatted message arguments.
func ErrorF(err error, op, format string, args ...any) *Error {
	return NewE(err, fmt.Sprintf(format, args...), op)
}

// FileLine returns the file and line in which the error
// occurred.
func (e *Error) FileLine() string {
	return e.fileLine
}

// Unwrap unwraps the original error message.
func (e *Error) Unwrap() error {
	return e.Err
}

// Wrap returns an error annotating err with a stack trace
// at the point Wrap is called, and the supplied message.
// If err is nil, Wrap returns nil.
func Wrap(err error, message string) *Error {
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

// wrappingError is the wrapping error features the error
// and file line in strings suitable for json.Marshal.
type wrappingError struct {
	Code      string `json:"code"`
	Message   string `json:"message"`
	Operation string `json:"operation"`
	Err       string `json:"error"`
	FileLine  string `json:"file_line"`
}

// MarshalJSON implements encoding/Marshaller to wrap the
// error as a string if there is one.
func (e *Error) MarshalJSON() ([]byte, error) {
	err := wrappingError{
		Code:      e.Code,
		Message:   e.Message,
		Operation: e.Operation,
	}
	if e.Err != nil {
		err.Err = e.Err.Error()
		err.FileLine = e.fileLine
	}
	return json.Marshal(err)
}

// UnmarshalJSON implements encoding/Marshaller to unmarshal
// the wrapping error to type Error.
func (e *Error) UnmarshalJSON(data []byte) error {
	var err wrappingError
	mErr := json.Unmarshal(data, &err)
	if mErr != nil {
		return mErr
	}
	e.Code = err.Code
	e.Message = err.Message
	e.Operation = err.Operation
	e.fileLine = err.FileLine
	if err.Err != "" {
		e.Err = errors.New(err.Err)
	}
	return nil
}
