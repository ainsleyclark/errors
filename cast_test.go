// Copyright 2022 Ainsley Clark. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package errors

import (
	"fmt"
	"reflect"
	"testing"
)

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
