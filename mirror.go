// Copyright 2022 Ainsley Clark. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package errors

import "errors"

// Unwrap calls the Go builtin errors.UnUnwrap
func Unwrap(err error) error {
	return errors.Unwrap(err)
}

// Is calls the Go builtin errors.Is
func Is(err, target error) bool {
	return errors.Is(err, target)
}

// As calls the Go builtin errors.As
func As(err error, target any) bool {
	return errors.As(err, target)
}
