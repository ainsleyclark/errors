// Copyright 2022 Ainsley Clark. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package errors

import (
	"errors"
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
			"<internal> " + wd + "/errors_test.go:28 - op: error, message",
		},
		"Nil Operation": {
			NewInternal(fmt.Errorf("error"), "message", ""),
			"<internal> " + wd + "/errors_test.go:32 - error, message",
		},
		"Nil Err": {
			NewInternal(nil, "message", ""),
			"<internal> " + wd + "/errors_test.go:36 - message",
		},
		"Nil Message": {
			NewInternal(fmt.Errorf("error"), "", ""),
			"<internal> " + wd + "/errors_test.go:40 - error",
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

func UtilTestError(t *testing.T, want, got *Error) {
	t.Helper()
	if !reflect.DeepEqual(want.Err, got.Err) {
		t.Fatalf("expecting %s, got %s", want.Err, got.Err)
	}
	if !reflect.DeepEqual(want.Message, got.Message) {
		t.Fatalf("expecting %s, got %s", want.Message, got.Message)
	}
	if !reflect.DeepEqual(want.Code, got.Code) {
		t.Fatalf("expecting %s, got %s", want.Code, got.Code)
	}
	if !reflect.DeepEqual(want.Operation, got.Operation) {
		t.Fatalf("expecting %s, got %s", want.Operation, got.Operation)
	}
}

func TestNewE(t *testing.T) {
	want := &Error{
		Code:      INTERNAL,
		Message:   "message",
		Operation: "op",
		Err:       fmt.Errorf("error"),
	}
	got := NewE(fmt.Errorf("error"), "message", "op")
	UtilTestError(t, want, got)
}

func TestErrorf(t *testing.T) {
	want := &Error{
		Code:      INTERNAL,
		Message:   "message: hello",
		Operation: "op",
		Err:       fmt.Errorf("error"),
	}
	got := ErrorF(fmt.Errorf("error"), "op", "message: %s", "hello")
	UtilTestError(t, want, got)
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
	e := NewE(fmt.Errorf("error"), "message", "op")
	got := e.ProgramCounters()
	want := 100
	if !reflect.DeepEqual(len(got), want) {
		t.Fatalf("expecting %d, got %d", want, got)
	}
}

func TestError_RuntimeFrames(t *testing.T) {
	e := NewE(fmt.Errorf("error"), "message", "op")
	got := e.RuntimeFrames()
	frame, _ := got.Next()
	want := "github.com/ainsleyclark/errors.TestError_RuntimeFrames"
	if !reflect.DeepEqual(want, frame.Function) {
		t.Fatalf("expecting %s, got %s", want, frame.Function)
	}
}

func TestError_StackTrace(t *testing.T) {
	e := NewE(fmt.Errorf("error"), "message", "op")
	got := e.StackTrace()
	want := "github.com/ainsleyclark/errors.TestError_StackTrace(): message"
	if !strings.Contains(got, want) {
		t.Fatalf("expecting %s to contain, got %s", want, got)
	}
}

func TestError_StackTraceSlice(t *testing.T) {
	e := NewE(fmt.Errorf("error"), "message", "op")
	got := e.StackTraceSlice()[0]
	want := "github.com/ainsleyclark/errors.TestError_StackTrace(): message"
	if reflect.DeepEqual(want, got) {
		t.Fatalf("expecting %s, got %s", want, got)
	}
}

func TestError_MarshalJSON(t *testing.T) {
	tt := map[string]struct {
		input *Error
		want  string
	}{
		"With Error": {
			NewInternal(errors.New("error"), "message", "op"),
			`{"code":"internal","message":"message","operation":"op","error":"error"`,
		},
		"No Error": {
			NewInternal(nil, "message", "op"),
			`{"code":"internal","message":"message","operation":"op","error":""`,
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			got, err := test.input.MarshalJSON()
			if err != nil {
				if !strings.Contains(err.Error(), test.want) {
					t.Fatalf("expecting %s to contain, got %s", test.want, got)
				}
				return
			}
			if !strings.Contains(string(got), test.want) {
				t.Fatalf("expecting %s, got %s", test.want, string(got))
			}
		})
	}
}

func TestError_UnmarshalJSON(t *testing.T) {
	tt := map[string]struct {
		input string
		want  any
	}{
		"Marshal Error": {
			`{"code":123","message":"message","operation":"op","error":"error","file_line":""}`,
			"invalid character",
		},
		"With Error": {
			`{"code":"internal","message":"message","operation":"op","error":"error","file_line":""}`,
			Error{
				Code:      "internal",
				Message:   "message",
				Operation: "op",
				Err:       errors.New("error"),
			},
		},
		"Without Error": {
			`{"code":"internal","message":"message","operation":"op","error":"","file_line":""}`,
			Error{
				Code:      "internal",
				Message:   "message",
				Operation: "op",
			},
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			e := &Error{}
			got := e.UnmarshalJSON([]byte(test.input))
			if got != nil {
				if !strings.Contains(got.Error(), fmt.Sprintf("%s", test.want)) {
					t.Fatalf("expecting %s to contain, got %s", test.want, got)
				}
				return
			}
			if !reflect.DeepEqual(test.want, *e) {
				t.Fatalf("expecting %+v, got %+v", test.want, e)
			}
		})
	}
}

func TestError_Scan(t *testing.T) {
	tt := map[string]struct {
		input any
		want  any
	}{
		"Success": {
			[]byte(`{"code": "code"}`),
			nil,
		},
		"Nil": {
			nil,
			nil,
		},
		"Unsupported Scan": {
			"wrong",
			"scan not supported for *errors.Error",
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			e := &Error{}
			got := e.Scan(test.input)
			if got != nil {
				if !strings.Contains(got.Error(), fmt.Sprintf("%s", test.want)) {
					t.Fatalf("expecting %s to contain, got %s", test.want, got)
				}
				return
			}
			if !reflect.DeepEqual(test.want, got) {
				t.Fatalf("expecting %+v, got %+v", test.want, e)
			}
		})
	}
}

func TestError_Value(t *testing.T) {
	e := &Error{}

	t.Run("Success", func(t *testing.T) {
		val, _ := e.Value(e)
		want := `{"code":"","message":"","operation":"","error":"","file_line":""}`
		got := val.([]byte)
		if !reflect.DeepEqual(want, string(got)) {
			t.Fatalf("expecting %+v, got %s", want, string(got))
		}
	})

	t.Run("Nil", func(t *testing.T) {
		got, _ := e.Value(nil)
		if !reflect.DeepEqual(nil, got) {
			t.Fatalf("expecting %+v, got %s", nil, got)
		}
	})
}
