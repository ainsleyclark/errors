// Copyright 2022 Ainsley Clark. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package errors

import "fmt"

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
