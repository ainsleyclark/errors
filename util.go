// Copyright 2022 Ainsley Clark. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package errors

import (
	"runtime"
	"strconv"
)

// newError is an alias for New by creating the pcs
// file line and constructing the error message.
func newError(err error, message, code, op string) *Error {
	_, file, line, _ := runtime.Caller(2)
	pcs := make([]uintptr, 100)
	_ = runtime.Callers(3, pcs)
	return &Error{
		Code:      code,
		Message:   message,
		Operation: op,
		Err:       err,
		fileLine:  file + ":" + strconv.Itoa(line),
		pcs:       pcs,
	}
}

// NewInternal returns an Error with a INTERNAL error code.
func NewInternal(err error, message, op string) *Error {
	return newError(err, message, INTERNAL, op)
}

// NewConflict returns an Error with a CONFLICT error code.
func NewConflict(err error, message, op string) *Error {
	return newError(err, message, CONFLICT, op)
}

// NewInvalid returns an Error with a INVALID error code.
func NewInvalid(err error, message, op string) *Error {
	return newError(err, message, INVALID, op)
}

// NewNotFound returns an Error with a NOTFOUND error code.
func NewNotFound(err error, message, op string) *Error {
	return newError(err, message, NOTFOUND, op)
}

// NewUnknown returns an Error with a UNKNOWN error code.
func NewUnknown(err error, message, op string) *Error {
	return newError(err, message, UNKNOWN, op)
}

// NewMaximumAttempts returns an Error with a MAXIMUMATTEMPTS error code.
func NewMaximumAttempts(err error, message, op string) *Error {
	return newError(err, message, MAXIMUMATTEMPTS, op)
}

// NewExpired returns an Error with a EXPIRED error code.
func NewExpired(err error, message, op string) *Error {
	return newError(err, message, EXPIRED, op)
}
