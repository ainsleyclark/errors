// Copyright 2022 Ainsley Clark. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package errors

import (
	"fmt"
	"net/http"
	"os"
	"reflect"
	"strings"
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
			"<internal> " + wd + "/errors_test.go:27 - op: error, message",
		},
		"Nil Operation": {
			NewInternal(fmt.Errorf("error"), "message", ""),
			"<internal> " + wd + "/errors_test.go:31 - error, message",
		},
		"Nil Err": {
			NewInternal(nil, "message", ""),
			"<internal> " + wd + "/errors_test.go:35 - message",
		},
		"Nil Message": {
			NewInternal(fmt.Errorf("error"), "", ""),
			"<internal> " + wd + "/errors_test.go:39 - error",
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

func TestErrorf(t *testing.T) {

}

func TestError_FileLine(t *testing.T) {
	e := &Error{fileLine: "fileline:20"}
	got := e.FileLine()
	want := "fileline:20"
	if !reflect.DeepEqual(want, got) {
		t.Fatalf("expecting %s, got %s", want, got)
	}
}

func TestWrap(t *testing.T) {
	got := Wrap(fmt.Errorf("error"), "message")
	if !reflect.DeepEqual("message", got.Message) {
		t.Fatalf("expecting message, got %s", got)
	}
	if !reflect.DeepEqual(fmt.Errorf("error"), got.Err) {
		t.Fatalf("expecting error, got %s", got)
	}
}

func TestWrap_NilError(t *testing.T) {
	got := Wrap(nil, "")
	var want *Error
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("expecting %s, got %s", want, got)
	}
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
	e := New(fmt.Errorf("error"), "message", "op")
	got := e.ProgramCounters()
	want := 100
	if !reflect.DeepEqual(len(got), want) {
		t.Fatalf("expecting %d, got %d", want, got)
	}
}

func TestError_RuntimeFrames(t *testing.T) {
	e := New(fmt.Errorf("error"), "message", "op")
	got := e.RuntimeFrames()
	frame, _ := got.Next()
	want := "github.com/ainsleyclark/errors.TestError_RuntimeFrames"
	if !reflect.DeepEqual(want, frame.Function) {
		t.Fatalf("expecting %s, got %s", want, frame.Function)
	}
}

func TestError_StackTrace(t *testing.T) {
	e := New(fmt.Errorf("error"), "message", "op")
	got := e.StackTrace()
	want := "github.com/ainsleyclark/errors.TestError_StackTrace(): message"
	if !strings.Contains(got, want) {
		t.Fatalf("expecting %s to contain, got %s", want, got)
	}
}

func TestError_StackTraceSlice(t *testing.T) {
	e := New(fmt.Errorf("error"), "message", "op")
	got := e.StackTraceSlice()[0]
	want := "github.com/ainsleyclark/errors.TestError_StackTrace(): message"
	if reflect.DeepEqual(want, got) {
		t.Fatalf("expecting %s, got %s", want, got)
	}
}
