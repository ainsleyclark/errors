// Copyright 2022 Ainsley Clark. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package errors

import (
	"bytes"
	"fmt"
	"net/http"
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
}

// Error returns the string representation of the error
// message.
func (e *Error) Error() string {
	var buf bytes.Buffer

	// Print the current operation in our stack, if any.
	if e.Operation != "" {
		buf.WriteString(fmt.Sprintf("%s: ", e.Operation))
	}

	// If wrapping an error, print its HasError() message.
	// Otherwise, print the error code & message.
	if e.Err != nil {
		buf.WriteString(e.Err.Error())
	} else {
		if e.Code != "" {
			buf.WriteString(fmt.Sprintf("<%s> ", e.Code))
		}
		buf.WriteString(e.Message)
	}

	return buf.String()
}

// New is a wrapper for the stdlib new function.
func New(err error, message, code, op string) error {
	return &Error{
		Code:      code,
		Message:   message,
		Operation: op,
		Err:       err,
	}
}

func Newf(err error, format string, args ...any) error {
	message := fmt.Sprintf(format, args...)
	return New(err, message, DefaultCode, "")
}

func Errorf(format string, args ...any) {

}

// Code returns the code of the root error, if available.
// Otherwise, returns INTERNAL.
func Code(err error) string {
	if err == nil {
		return ""
	} else if e, ok := err.(*Error); ok && e.Code != "" {
		return e.Code
	} else if ok && e.Err != nil {
		return Code(e.Err)
	}
	return INTERNAL
}

// Message returns the human-readable message of the error,
// if available. Otherwise, returns a generic error
// message.
func Message(err error) string {
	if err == nil {
		return ""
	} else if e, ok := err.(*Error); ok && e.Message != "" {
		return e.Message
	} else if ok && e.Err != nil {
		return Message(e.Err)
	}
	return GlobalError
}

// ToError Returns an application error from input. If The type
// is not of type Error, nil will be returned.
func ToError(err any) *Error {
	switch v := err.(type) {
	case *Error:
		return v
	case Error:
		return &v
	case error:
		return &Error{Err: fmt.Errorf(v.Error())}
	case string:
		return &Error{Err: fmt.Errorf(v)}
	default:
		return nil
	}
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
