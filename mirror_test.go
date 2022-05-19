// Copyright 2022 Ainsley Clark. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package errors

import (
	"fmt"
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	want := fmt.Errorf("error")
	got := New("error")
	if !reflect.DeepEqual(want, got) {
		t.Fatalf("expecting %s, got %s", want, got)
	}
}

func TestUnwrap(t *testing.T) {
	e := NewE(fmt.Errorf("error"), "message", "op")
	got := Unwrap(e)
	want := "error"
	if !reflect.DeepEqual(want, got.Error()) {
		t.Fatalf("expecting %s, got %s", want, got.Error())
	}
}

func TestIs(t *testing.T) {
	err := fmt.Errorf("error")
	e := NewE(err, "message", "op")
	got := Is(e, err)
	if !got {
		t.Fatalf("expecting true, got %t", got)
	}
}

func TestAs(t *testing.T) {
	err := fmt.Errorf("error")
	e := NewE(err, "message", "op")
	got := As(e, &err)
	if !got {
		t.Fatalf("expecting true, got %t", got)
	}
}
