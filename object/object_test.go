package object

import (
	"math"
	"sugu/ast"
	"sugu/token"
	"testing"
)

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

func TestFunctionInspect(t *testing.T) {
	// 名前付き関数のテスト
	namedFn := &Function{
		Name: "add",
		Parameters: []*ast.Identifier{
			{Token: token.Token{Type: token.IDENT, Literal: "a"}, Value: "a"},
			{Token: token.Token{Type: token.IDENT, Literal: "b"}, Value: "b"},
		},
		Body: &ast.BlockStatement{},
		Env:  nil,
	}
	expectedNamed := "func add(a, b) => { ... }"
	if namedFn.Inspect() != expectedNamed {
		t.Errorf("Function.Inspect() wrong for named function. got=%s, want=%s",
			namedFn.Inspect(), expectedNamed)
	}
	if namedFn.Type() != FUNCTION_OBJ {
		t.Errorf("Function.Type() wrong. got=%s", namedFn.Type())
	}

	// 無名関数のテスト
	anonFn := &Function{
		Name: "",
		Parameters: []*ast.Identifier{
			{Token: token.Token{Type: token.IDENT, Literal: "x"}, Value: "x"},
		},
		Body: &ast.BlockStatement{},
		Env:  nil,
	}
	expectedAnon := "func(x) => { ... }"
	if anonFn.Inspect() != expectedAnon {
		t.Errorf("Function.Inspect() wrong for anonymous function. got=%s, want=%s",
			anonFn.Inspect(), expectedAnon)
	}

	// パラメータなしの関数
	noParamFn := &Function{
		Name:       "greet",
		Parameters: []*ast.Identifier{},
		Body:       &ast.BlockStatement{},
		Env:        nil,
	}
	expectedNoParam := "func greet() => { ... }"
	if noParamFn.Inspect() != expectedNoParam {
		t.Errorf("Function.Inspect() wrong for no-param function. got=%s, want=%s",
			noParamFn.Inspect(), expectedNoParam)
	}
}

func TestNumberInspectLargeValues(t *testing.T) {
	tests := []struct {
		value    float64
		expected string
	}{
		// 大きな整数値（int64の範囲内）
		{1e15, "1000000000000000"},
		// 非常に大きな値（int64で正確に表現できない）
		{1e20, "1e+20"},
		{1e19, "1e+19"},
		// NaN と Infinity
		{math.Inf(1), "+Inf"},
		{math.Inf(-1), "-Inf"},
	}

	for _, tt := range tests {
		n := &Number{Value: tt.value}
		if n.Inspect() != tt.expected {
			t.Errorf("Number.Inspect() for %v wrong. got=%s, want=%s",
				tt.value, n.Inspect(), tt.expected)
		}
	}

	// NaNは特別なテスト（NaN != NaN なので）
	nanNum := &Number{Value: math.NaN()}
	if nanNum.Inspect() != "NaN" {
		t.Errorf("Number.Inspect() for NaN wrong. got=%s, want=NaN", nanNum.Inspect())
	}
}
