// Copyright 2022 Ainsley Clark. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package errors

import (
	"fmt"
	"net/http"
	"os"
	"reflect"
	"testing"
)

func TestError_Error(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed: %s", err.Error())
	}

	tt := map[string]struct {
		input *Error
		want  string
	}{
		"Normal": {
			NewInternal(fmt.Errorf("error"), "message", "op"),
			"<internal> " + wd + "/errors_test.go:26 - op: error, message",
		},
		"Nil Operation": {
			NewInternal(fmt.Errorf("error"), "message", ""),
			"<internal> " + wd + "/errors_test.go:30 - error, message",
		},
		"Nil Err": {
			NewInternal(nil, "message", ""),
			"<internal> " + wd + "/errors_test.go:34 - message",
		},
		"Nil Message": {
			NewInternal(fmt.Errorf("error"), "", ""),
			"<internal> " + wd + "/errors_test.go:38 - error",
		},
		"Message Error": {
			&Error{Message: "message", Err: fmt.Errorf("err")},
			"err, message",
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			got := test.input.Error()
			if !reflect.DeepEqual(test.want, got) {
				t.Fatalf("expecting %s, got %s", test.want, got)
			}
		})
	}
}

func TestNew(t *testing.T) {
	//want := &Error{
	//	Code:      INVALID,
	//	Message:   "message",
	//	Operation: "op",
	//	Err:       fmt.Errorf("error"),
	//}
	//got := New(fmt.Errorf("error"), "message", INVALID)
	//
	//fmt.Println(got.fileLine)
	//
	//if !reflect.DeepEqual(want, got) {
	//	t.Fatalf("expecting %+v, got %+v", want, got)
	//}
}

func TestWrap(t *testing.T) {
	//want := "message: error"
	//got := Wrap(fmt.Errorf("error"), "message")
	//assert.Equal(t, want, got.Error())
}

func TestError_HTTPStatusCode(t *testing.T) {
	tt := map[string]struct {
		input Error
		want  int
	}{
		"Conflict": {
			Error{Code: CONFLICT},
			http.StatusConflict,
		},
		"Internal": {
			Error{Code: INTERNAL},
			http.StatusInternalServerError,
		},
		"Invalid": {
			Error{Code: INVALID},
			http.StatusBadRequest,
		},
		"Not Found": {
			Error{Code: NOTFOUND},
			http.StatusNotFound,
		},
		"Unknown": {
			Error{Code: UNKNOWN},
			http.StatusInternalServerError,
		},
		"Maximum Attempts": {
			Error{Code: MAXIMUMATTEMPTS},
			http.StatusTooManyRequests,
		},
		"Expired": {
			Error{Code: EXPIRED},
			http.StatusPaymentRequired,
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			got := test.input.HTTPStatusCode()
			if !reflect.DeepEqual(test.want, got) {
				t.Fatalf("expecting %d, got %d", test.want, got)
			}
		})
	}
}

func TestError_ProgramCounters(t *testing.T) {
	e := newError(fmt.Errorf("error"), "message", INTERNAL, "op")
	got := e.ProgramCounters()
	want := 100
	if !reflect.DeepEqual(len(got), want) {
		t.Fatalf("expecting %d, got %d", want, got)
	}
}

func TestError_RuntimeFrames(t *testing.T) {
	e := newError(fmt.Errorf("error"), "message", INTERNAL, "op")
	got := e.RuntimeFrames()
	fmt.Println(got.Next())
	//if !reflect.DeepEqual(len(got), want) {
	//	t.Fatalf("expecting %d, got %d", want, got)
	//}
}

//func TestError_StackTrace(t *testing.T) {
//	e := new(fmt.Errorf("error"), "message", INTERNAL, "op")
//	got := e.StackTrace()
//	fmt.Println(got)
//	//if !reflect.DeepEqual(len(got), want) {
//	//	t.Fatalf("expecting %d, got %d", want, got)
//	//}
//}
