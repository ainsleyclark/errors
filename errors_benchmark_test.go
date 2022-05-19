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
		_ = New(errors.New("error"), "message", INTERNAL)
	}
}

func BenchmarkNewInternal(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = New(errors.New("error"), "message", INTERNAL)
	}
}
