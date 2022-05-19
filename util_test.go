// Copyright 2022 Ainsley Clark. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package errors

import (
	"reflect"
	"testing"
)

func TestNewInternal(t *testing.T) {
	got := NewInternal(nil, "message", "op")
	want := INTERNAL
	if !reflect.DeepEqual(want, got.Code) {
		t.Fatalf("expecting %s, got %s", want, got)
	}
}

func TestNewConflict(t *testing.T) {
	got := NewConflict(nil, "message", "op")
	want := CONFLICT
	if !reflect.DeepEqual(want, got.Code) {
		t.Fatalf("expecting %s, got %s", want, got)
	}
}

func TestNewInvalid(t *testing.T) {
	got := NewInvalid(nil, "message", "op")
	want := INVALID
	if !reflect.DeepEqual(want, got.Code) {
		t.Fatalf("expecting %s, got %s", want, got)
	}
}

func TestNewNotFound(t *testing.T) {
	got := NewNotFound(nil, "message", "op")
	want := NOTFOUND
	if !reflect.DeepEqual(want, got.Code) {
		t.Fatalf("expecting %s, got %s", want, got)
	}
}

func TestNewUnknown(t *testing.T) {
	got := NewUnknown(nil, "message", "op")
	want := UNKNOWN
	if !reflect.DeepEqual(want, got.Code) {
		t.Fatalf("expecting %s, got %s", want, got)
	}
}

func TestNewMaximumAttempts(t *testing.T) {
	got := NewMaximumAttempts(nil, "message", "op")
	want := MAXIMUMATTEMPTS
	if !reflect.DeepEqual(want, got.Code) {
		t.Fatalf("expecting %s, got %s", want, got)
	}
}

func TestNewExpired(t *testing.T) {
	got := NewExpired(nil, "message", "op")
	want := EXPIRED
	if !reflect.DeepEqual(want, got.Code) {
		t.Fatalf("expecting %s, got %s", want, got)
	}
}
