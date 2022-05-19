// Copyright 2022 Ainsley Clark. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package errors

import (
	"errors"
	"testing"
)

func BenchmarkNew(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = NewE(errors.New("error"), "message", INTERNAL)
	}
}

func BenchmarkNewInternal(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = NewE(errors.New("error"), "message", INTERNAL)
	}
}

func BenchmarkError_Error(b *testing.B) {
	e := NewE(errors.New("error"), "message", INTERNAL)
	for i := 0; i < b.N; i++ {
		_ = e.Error()
	}
}

func BenchmarkError_Code(b *testing.B) {
	e := NewE(errors.New("error"), "message", INTERNAL)
	for i := 0; i < b.N; i++ {
		_ = Code(e)
	}
}

func BenchmarkError_Message(b *testing.B) {
	e := NewE(errors.New("error"), "message", INTERNAL)
	for i := 0; i < b.N; i++ {
		_ = Message(e)
	}
}

func BenchmarkError_ToError(b *testing.B) {
	e := NewE(errors.New("error"), "message", INTERNAL)
	for i := 0; i < b.N; i++ {
		_ = ToError(e)
	}
}

func BenchmarkError_HTTPStatusCode(b *testing.B) {
	e := NewE(errors.New("error"), "message", INTERNAL)
	for i := 0; i < b.N; i++ {
		_ = e.HTTPStatusCode()
	}
}
