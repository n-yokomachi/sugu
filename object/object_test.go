package object

import "testing"

func TestNumberInspect(t *testing.T) {
	tests := []struct {
		value    float64
		expected string
	}{
		{42, "42"},
		{3.14, "3.14"},
		{0, "0"},
		{-10, "-10"},
		{100.5, "100.5"},
	}

	for _, tt := range tests {
		n := &Number{Value: tt.value}
		if n.Inspect() != tt.expected {
			t.Errorf("Number.Inspect() wrong. got=%s, want=%s", n.Inspect(), tt.expected)
		}
		if n.Type() != NUMBER_OBJ {
			t.Errorf("Number.Type() wrong. got=%s, want=%s", n.Type(), NUMBER_OBJ)
		}
	}
}

func TestStringInspect(t *testing.T) {
	s := &String{Value: "hello world"}
	if s.Inspect() != "hello world" {
		t.Errorf("String.Inspect() wrong. got=%s", s.Inspect())
	}
	if s.Type() != STRING_OBJ {
		t.Errorf("String.Type() wrong. got=%s", s.Type())
	}
}

func TestBooleanInspect(t *testing.T) {
	tests := []struct {
		value    bool
		expected string
	}{
		{true, "true"},
		{false, "false"},
	}

	for _, tt := range tests {
		b := &Boolean{Value: tt.value}
		if b.Inspect() != tt.expected {
			t.Errorf("Boolean.Inspect() wrong. got=%s, want=%s", b.Inspect(), tt.expected)
		}
		if b.Type() != BOOLEAN_OBJ {
			t.Errorf("Boolean.Type() wrong. got=%s", b.Type())
		}
	}
}

func TestNullInspect(t *testing.T) {
	n := &Null{}
	if n.Inspect() != "null" {
		t.Errorf("Null.Inspect() wrong. got=%s", n.Inspect())
	}
	if n.Type() != NULL_OBJ {
		t.Errorf("Null.Type() wrong. got=%s", n.Type())
	}
}

func TestErrorInspect(t *testing.T) {
	e := &Error{Message: "something went wrong"}
	if e.Inspect() != "ERROR: something went wrong" {
		t.Errorf("Error.Inspect() wrong. got=%s", e.Inspect())
	}
	if e.Type() != ERROR_OBJ {
		t.Errorf("Error.Type() wrong. got=%s", e.Type())
	}
}

func TestReturnValueInspect(t *testing.T) {
	rv := &ReturnValue{Value: &Number{Value: 42}}
	if rv.Inspect() != "42" {
		t.Errorf("ReturnValue.Inspect() wrong. got=%s", rv.Inspect())
	}
	if rv.Type() != RETURN_OBJ {
		t.Errorf("ReturnValue.Type() wrong. got=%s", rv.Type())
	}
}

func TestBuiltinInspect(t *testing.T) {
	b := &Builtin{Fn: func(args ...Object) Object { return nil }}
	if b.Inspect() != "builtin function" {
		t.Errorf("Builtin.Inspect() wrong. got=%s", b.Inspect())
	}
	if b.Type() != BUILTIN_OBJ {
		t.Errorf("Builtin.Type() wrong. got=%s", b.Type())
	}
}
