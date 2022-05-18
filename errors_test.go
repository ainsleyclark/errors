// Copyright 2022 Ainsley Clark. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package errors

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestError_Error(t *testing.T) {
	tt := map[string]struct {
		input *Error
		want  string
	}{
		"Normal": {
			&Error{Code: INTERNAL, Message: "test", Operation: "op", Err: fmt.Errorf("err")},
			"op: err",
		},
		"Nil Operation": {
			&Error{Code: INTERNAL, Message: "test", Operation: "", Err: fmt.Errorf("err")},
			"err",
		},
		"Nil Err": {
			&Error{Code: INTERNAL, Message: "test", Operation: "", Err: nil},
			"<internal> test",
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

func TestError_Code(t *testing.T) {
	tt := map[string]struct {
		input error
		want  string
	}{
		"Normal": {
			&Error{Code: INTERNAL, Message: "test", Operation: "op", Err: fmt.Errorf("err")},
			"internal",
		},
		"Nil Input": {
			nil,
			"",
		},
		"Nil Code": {
			&Error{Code: "", Message: "test", Operation: "op", Err: fmt.Errorf("err")},
			"internal",
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			got := Code(test.input)
			if !reflect.DeepEqual(test.want, got) {
				t.Fatalf("expecting %s, got %s", test.want, got)
			}
		})
	}
}

func Test_Message(t *testing.T) {
	tt := map[string]struct {
		input error
		want  string
	}{
		"Normal": {
			&Error{Code: INTERNAL, Message: "test", Operation: "op", Err: fmt.Errorf("err")},
			"test",
		},
		"Nil Input": {
			nil,
			"",
		},
		"Nil Message": {
			&Error{Code: "", Message: "", Operation: "op", Err: fmt.Errorf("err")},
			GlobalError,
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			got := Message(test.input)
			if !reflect.DeepEqual(test.want, got) {
				t.Fatalf("expecting %s, got %s", test.want, got)
			}
		})
	}
}

func TestError_ToError(t *testing.T) {
	tt := map[string]struct {
		input any
		want  *Error
	}{
		"Pointer": {
			&Error{Code: INTERNAL, Message: "test", Operation: "op", Err: fmt.Errorf("err")},
			&Error{Code: INTERNAL, Message: "test", Operation: "op", Err: fmt.Errorf("err")},
		},
		"Non Pointer": {
			Error{Code: INTERNAL, Message: "test", Operation: "op", Err: fmt.Errorf("err")},
			&Error{Code: INTERNAL, Message: "test", Operation: "op", Err: fmt.Errorf("err")},
		},
		"Error": {
			fmt.Errorf("err"),
			&Error{Err: fmt.Errorf("err")},
		},
		"String": {
			"err",
			&Error{Err: fmt.Errorf("err")},
		},
		"Default": {
			nil,
			nil,
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			got := ToError(test.input)
			if !reflect.DeepEqual(test.want, got) {
				t.Fatalf("expecting %s, got %s", test.want, got)
			}
		})
	}
}

//func TestNew(t *testing.T) {
//	want := &Error{
//		Code:      INVALID,
//		Message:   "message",
//		Operation: "op",
//		Err:       fmt.Errorf("error"),
//	}
//	got := New(fmt.Errorf("error"), "message", INVALID, "op")
//	if !reflect.DeepEqual(want, got) {
//		t.Fatalf("expecting %+v, got %+v", want, got)
//	}
//}

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
	e := new(fmt.Errorf("error"), "message", INTERNAL, "op")
	got := e.ProgramCounters()
	want := 100
	if !reflect.DeepEqual(len(got), want) {
		t.Fatalf("expecting %d, got %d", want, got)
	}
}

func TestError_RuntimeFrames(t *testing.T) {
	e := new(fmt.Errorf("error"), "message", INTERNAL, "op")
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
