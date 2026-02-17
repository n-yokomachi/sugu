package evaluator

import (
	"fmt"
	"os"
	"sugu/lexer"
	"sugu/object"
	"sugu/parser"
	"testing"
)

func TestEvalNumberExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{"5", 5},
		{"10", 10},
		{"-5", -5},
		{"-10", -10},
		{"5 + 5 + 5 + 5 - 10", 10},
		{"2 * 2 * 2 * 2 * 2", 32},
		{"-50 + 100 + -50", 0},
		{"5 * 2 + 10", 20},
		{"5 + 2 * 10", 25},
		{"20 + 2 * -10", 0},
		{"50 / 2 * 2 + 10", 60},
		{"2 * (5 + 10)", 30},
		{"3 * 3 * 3 + 10", 37},
		{"3 * (3 * 3) + 10", 37},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
		{"10 % 3", 1},
		{"10 % 5", 0},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testNumberObject(t, evaluated, tt.expected)
	}
}

func TestEvalBooleanExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
		{"1 < 2", true},
		{"1 > 2", false},
		{"1 < 1", false},
		{"1 > 1", false},
		{"1 == 1", true},
		{"1 != 1", false},
		{"1 == 2", false},
		{"1 != 2", true},
		{"1 <= 1", true},
		{"1 >= 1", true},
		{"1 <= 2", true},
		{"2 >= 1", true},
		{"true == true", true},
		{"false == false", true},
		{"true == false", false},
		{"true != false", true},
		{"false != true", true},
		{"(1 < 2) == true", true},
		{"(1 < 2) == false", false},
		{"(1 > 2) == true", false},
		{"(1 > 2) == false", true},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func TestBangOperator(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"!true", false},
		{"!false", true},
		{"!5", false},
		{"!!true", true},
		{"!!false", false},
		{"!!5", true},
		{"!null", true},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func TestStringLiteral(t *testing.T) {
	input := `"Hello World!"`

	evaluated := testEval(input)
	str, ok := evaluated.(*object.String)
	if !ok {
		t.Fatalf("object is not String. got=%T (%+v)", evaluated, evaluated)
	}

	if str.Value != "Hello World!" {
		t.Errorf("String has wrong value. got=%q", str.Value)
	}
}

func TestStringConcatenation(t *testing.T) {
	input := `"Hello" + " " + "World!"`

	evaluated := testEval(input)
	str, ok := evaluated.(*object.String)
	if !ok {
		t.Fatalf("object is not String. got=%T (%+v)", evaluated, evaluated)
	}

	if str.Value != "Hello World!" {
		t.Errorf("String has wrong value. got=%q", str.Value)
	}
}

func TestNullLiteral(t *testing.T) {
	input := "null"

	evaluated := testEval(input)
	if evaluated != NULL {
		t.Errorf("object is not NULL. got=%T (%+v)", evaluated, evaluated)
	}
}

func TestIfElseExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"if (true) { 10 }", 10},
		{"if (false) { 10 }", nil},
		{"if (1) { 10 }", 10},
		{"if (1 < 2) { 10 }", 10},
		{"if (1 > 2) { 10 }", nil},
		{"if (1 > 2) { 10 } else { 20 }", 20},
		{"if (1 < 2) { 10 } else { 20 }", 10},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			testNumberObject(t, evaluated, float64(integer))
		} else {
			testNullObject(t, evaluated)
		}
	}
}

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{"return 10;", 10},
		{"return 10; 9;", 10},
		{"return 2 * 5; 9;", 10},
		{"9; return 2 * 5; 9;", 10},
		{`
if (10 > 1) {
  if (10 > 1) {
    return 10;
  }
  return 1;
}
`, 10},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testNumberObject(t, evaluated, tt.expected)
	}
}

func TestErrorHandling(t *testing.T) {
	tests := []struct {
		input           string
		expectedMessage string
	}{
		{
			"5 + true;",
			"type mismatch: NUMBER + BOOLEAN",
		},
		{
			"5 + true; 5;",
			"type mismatch: NUMBER + BOOLEAN",
		},
		{
			"-true",
			"unknown operator: -BOOLEAN",
		},
		{
			"true + false;",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"5; true + false; 5",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"if (10 > 1) { true + false; }",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"foobar",
			"line 1, column 1: identifier not found: foobar",
		},
		{
			`"Hello" - "World"`,
			"unknown operator: STRING - STRING",
		},
		{
			"10 / 0",
			"division by zero",
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)

		errObj, ok := evaluated.(*object.Error)
		if !ok {
			t.Errorf("no error object returned. got=%T(%+v)", evaluated, evaluated)
			continue
		}

		if errObj.Message != tt.expectedMessage {
			t.Errorf("wrong error message. expected=%q, got=%q", tt.expectedMessage, errObj.Message)
		}
	}
}

func TestVariableStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{"mut a = 5; a;", 5},
		{"mut a = 5 * 5; a;", 25},
		{"mut a = 5; mut b = a; b;", 5},
		{"mut a = 5; mut b = a; mut c = a + b + 5; c;", 15},
		{"const PI = 3.14; PI;", 3.14},
	}

	for _, tt := range tests {
		testNumberObject(t, testEval(tt.input), tt.expected)
	}
}

func TestFunctionObject(t *testing.T) {
	input := "func(x) => { x + 2; };"

	evaluated := testEval(input)
	fn, ok := evaluated.(*object.Function)
	if !ok {
		t.Fatalf("object is not Function. got=%T (%+v)", evaluated, evaluated)
	}

	if len(fn.Parameters) != 1 {
		t.Fatalf("function has wrong parameters. Parameters=%+v", fn.Parameters)
	}

	if fn.Parameters[0].String() != "x" {
		t.Fatalf("parameter is not 'x'. got=%q", fn.Parameters[0])
	}

	expectedBody := "{ (x + 2) }"
	if fn.Body.String() != expectedBody {
		t.Fatalf("body is not %q. got=%q", expectedBody, fn.Body.String())
	}
}

func TestFunctionApplication(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{"func identity(x) => { x; }; identity(5);", 5},
		{"func identity(x) => { return x; }; identity(5);", 5},
		{"func double(x) => { x * 2; }; double(5);", 10},
		{"func add(x, y) => { x + y; }; add(5, 5);", 10},
		{"func add(x, y) => { x + y; }; add(5 + 5, add(5, 5));", 20},
		{"func(x) => { x; }(5)", 5},
	}

	for _, tt := range tests {
		testNumberObject(t, testEval(tt.input), tt.expected)
	}
}

func TestClosures(t *testing.T) {
	input := `
func newAdder(x) => {
  func(y) => { x + y; };
};

mut addTwo = newAdder(2);
addTwo(2);`

	testNumberObject(t, testEval(input), 4)
}

func TestWhileStatement(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{`
mut x = 0;
while (x < 5) {
  x = x + 1;
}
x;
`, 5},
		{`
mut sum = 0;
mut i = 1;
while (i <= 10) {
  sum = sum + i;
  i = i + 1;
}
sum;
`, 55},
	}

	for _, tt := range tests {
		testNumberObject(t, testEval(tt.input), tt.expected)
	}
}

func TestBreakStatement(t *testing.T) {
	input := `
mut x = 0;
while (true) {
  x = x + 1;
  if (x == 5) {
    break;
  }
}
x;
`
	testNumberObject(t, testEval(input), 5)
}

func TestContinueStatement(t *testing.T) {
	input := `
mut sum = 0;
mut i = 0;
while (i < 10) {
  i = i + 1;
  if (i % 2 == 0) {
    continue;
  }
  sum = sum + i;
}
sum;
`
	// 1 + 3 + 5 + 7 + 9 = 25
	testNumberObject(t, testEval(input), 25)
}

func TestForStatement(t *testing.T) {
	input := `
mut sum = 0;
for (mut i = 1; i <= 10; i = i + 1) {
  sum = sum + i;
}
sum;
`
	testNumberObject(t, testEval(input), 55)
}

func TestSwitchStatement(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{`
mut x = 2;
mut result = 0;
switch (x) {
  case 1: {
    result = 10;
  }
  case 2: {
    result = 20;
  }
  case 3: {
    result = 30;
  }
}
result;
`, 20},
		{`
mut x = 5;
mut result = 0;
switch (x) {
  case 1: {
    result = 10;
  }
  default: {
    result = 99;
  }
}
result;
`, 99},
	}

	for _, tt := range tests {
		testNumberObject(t, testEval(tt.input), tt.expected)
	}
}

func TestLogicalOperators(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"true && true", true},
		{"true && false", false},
		{"false && true", false},
		{"false && false", false},
		{"true || true", true},
		{"true || false", true},
		{"false || true", true},
		{"false || false", false},
		// 短絡評価
		{"false && (1 / 0)", false}, // エラーにならない
		{"true || (1 / 0)", true},   // エラーにならない
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		switch expected := tt.expected.(type) {
		case bool:
			testBooleanObject(t, evaluated, expected)
		}
	}
}

func TestBuiltinFunctions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`type(1)`, "NUMBER"},
		{`type("hello")`, "STRING"},
		{`type(true)`, "BOOLEAN"},
		{`type(null)`, "NULL"},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)

		switch expected := tt.expected.(type) {
		case string:
			str, ok := evaluated.(*object.String)
			if !ok {
				t.Errorf("object is not String. got=%T (%+v)", evaluated, evaluated)
				continue
			}
			if str.Value != expected {
				t.Errorf("wrong value. expected=%q, got=%q", expected, str.Value)
			}
		}
	}
}

func TestBuiltinLen(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`len("")`, 0},
		{`len("four")`, 4},
		{`len("hello world")`, 11},
		{`len([1, 2, 3])`, 3},
		{`len([])`, 0},
		{`len({"a": 1, "b": 2})`, 2},
		{`len({})`, 0},
		{`len(1)`, "argument to `len` not supported, got NUMBER"},
		{`len("one", "two")`, "wrong number of arguments. got=2, want=1"},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)

		switch expected := tt.expected.(type) {
		case int:
			testNumberObject(t, evaluated, float64(expected))
		case string:
			errObj, ok := evaluated.(*object.Error)
			if !ok {
				t.Errorf("object is not Error. got=%T (%+v)", evaluated, evaluated)
				continue
			}
			if errObj.Message != expected {
				t.Errorf("wrong error message. expected=%q, got=%q", expected, errObj.Message)
			}
		}
	}
}

func TestBuiltinPush(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`push([1, 2], 3)`, []int{1, 2, 3}},
		{`push([], 1)`, []int{1}},
		{`mut arr = [1]; push(arr, 2); arr`, []int{1}}, // 元の配列は変更されない
		{`push(1, 1)`, "argument to `push` must be ARRAY, got NUMBER"},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)

		switch expected := tt.expected.(type) {
		case []int:
			arr, ok := evaluated.(*object.Array)
			if !ok {
				t.Errorf("object is not Array. got=%T (%+v)", evaluated, evaluated)
				continue
			}
			if len(arr.Elements) != len(expected) {
				t.Errorf("wrong num of elements. got=%d, want=%d", len(arr.Elements), len(expected))
				continue
			}
			for i, exp := range expected {
				testNumberObject(t, arr.Elements[i], float64(exp))
			}
		case string:
			errObj, ok := evaluated.(*object.Error)
			if !ok {
				t.Errorf("object is not Error. got=%T (%+v)", evaluated, evaluated)
				continue
			}
			if errObj.Message != expected {
				t.Errorf("wrong error message. expected=%q, got=%q", expected, errObj.Message)
			}
		}
	}
}

func TestBuiltinFirstLastRest(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`first([1, 2, 3])`, 1},
		{`first([])`, nil},
		{`last([1, 2, 3])`, 3},
		{`last([])`, nil},
		{`rest([1, 2, 3])`, []int{2, 3}},
		{`rest([1])`, []int{}},
		{`rest([])`, nil},
		{`first(1)`, "argument to `first` must be ARRAY, got NUMBER"},
		{`last(1)`, "argument to `last` must be ARRAY, got NUMBER"},
		{`rest(1)`, "argument to `rest` must be ARRAY, got NUMBER"},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)

		switch expected := tt.expected.(type) {
		case int:
			testNumberObject(t, evaluated, float64(expected))
		case nil:
			testNullObject(t, evaluated)
		case []int:
			arr, ok := evaluated.(*object.Array)
			if !ok {
				t.Errorf("object is not Array. got=%T (%+v)", evaluated, evaluated)
				continue
			}
			if len(arr.Elements) != len(expected) {
				t.Errorf("wrong num of elements. got=%d, want=%d", len(arr.Elements), len(expected))
				continue
			}
			for i, exp := range expected {
				testNumberObject(t, arr.Elements[i], float64(exp))
			}
		case string:
			errObj, ok := evaluated.(*object.Error)
			if !ok {
				t.Errorf("object is not Error. got=%T (%+v)", evaluated, evaluated)
				continue
			}
			if errObj.Message != expected {
				t.Errorf("wrong error message. expected=%q, got=%q", expected, errObj.Message)
			}
		}
	}
}

func TestBuiltinKeysValues(t *testing.T) {
	// keys と values の順序は不定なので、要素数のみチェック
	tests := []struct {
		input        string
		expectedLen  int
		expectedType string
	}{
		{`len(keys({"a": 1, "b": 2}))`, 2, "number"},
		{`len(keys({}))`, 0, "number"},
		{`len(values({"a": 1, "b": 2}))`, 2, "number"},
		{`len(values({}))`, 0, "number"},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testNumberObject(t, evaluated, float64(tt.expectedLen))
	}
}

func TestBuiltinKeysValuesErrors(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`keys(1)`, "argument to `keys` must be MAP, got NUMBER"},
		{`values(1)`, "argument to `values` must be MAP, got NUMBER"},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		errObj, ok := evaluated.(*object.Error)
		if !ok {
			t.Errorf("object is not Error. got=%T (%+v)", evaluated, evaluated)
			continue
		}
		if errObj.Message != tt.expected {
			t.Errorf("wrong error message. expected=%q, got=%q", tt.expected, errObj.Message)
		}
	}
}

// ヘルパー関数

func testEval(input string) object.Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	env := object.NewEnvironment()

	return Eval(program, env)
}

func testNumberObject(t *testing.T, obj object.Object, expected float64) bool {
	result, ok := obj.(*object.Number)
	if !ok {
		t.Errorf("object is not Number. got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value. got=%f, want=%f", result.Value, expected)
		return false
	}
	return true
}

func testBooleanObject(t *testing.T, obj object.Object, expected bool) bool {
	result, ok := obj.(*object.Boolean)
	if !ok {
		t.Errorf("object is not Boolean. got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value. got=%t, want=%t", result.Value, expected)
		return false
	}
	return true
}

func testNullObject(t *testing.T, obj object.Object) bool {
	if obj != NULL {
		t.Errorf("object is not NULL. got=%T (%+v)", obj, obj)
		return false
	}
	return true
}

func TestArrayLiterals(t *testing.T) {
	input := "[1, 2 * 2, 3 + 3]"

	evaluated := testEval(input)
	result, ok := evaluated.(*object.Array)
	if !ok {
		t.Fatalf("object is not Array. got=%T (%+v)", evaluated, evaluated)
	}

	if len(result.Elements) != 3 {
		t.Fatalf("array has wrong num of elements. got=%d", len(result.Elements))
	}

	testNumberObject(t, result.Elements[0], 1)
	testNumberObject(t, result.Elements[1], 4)
	testNumberObject(t, result.Elements[2], 6)
}

func TestArrayIndexExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{
			"[1, 2, 3][0]",
			1,
		},
		{
			"[1, 2, 3][1]",
			2,
		},
		{
			"[1, 2, 3][2]",
			3,
		},
		{
			"mut i = 0; [1][i];",
			1,
		},
		{
			"[1, 2, 3][1 + 1];",
			3,
		},
		{
			"mut myArray = [1, 2, 3]; myArray[2];",
			3,
		},
		{
			"mut myArray = [1, 2, 3]; myArray[0] + myArray[1] + myArray[2];",
			6,
		},
		{
			"mut myArray = [1, 2, 3]; mut i = myArray[0]; myArray[i];",
			2,
		},
		{
			"[1, 2, 3][3]",
			nil,
		},
		{
			"[1, 2, 3][-1]",
			nil,
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			testNumberObject(t, evaluated, float64(integer))
		} else {
			testNullObject(t, evaluated)
		}
	}
}

func TestStringIndexExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{
			`"hello"[0]`,
			"h",
		},
		{
			`"hello"[4]`,
			"o",
		},
		{
			`"hello"[5]`,
			nil,
		},
		{
			`"hello"[-1]`,
			nil,
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		str, ok := tt.expected.(string)
		if ok {
			result, ok := evaluated.(*object.String)
			if !ok {
				t.Errorf("object is not String. got=%T (%+v)", evaluated, evaluated)
				continue
			}
			if result.Value != str {
				t.Errorf("object has wrong value. got=%q, want=%q", result.Value, str)
			}
		} else {
			testNullObject(t, evaluated)
		}
	}
}

func TestArrayInspect(t *testing.T) {
	input := "[1, 2, 3]"

	evaluated := testEval(input)
	result, ok := evaluated.(*object.Array)
	if !ok {
		t.Fatalf("object is not Array. got=%T (%+v)", evaluated, evaluated)
	}

	expected := "[1, 2, 3]"
	if result.Inspect() != expected {
		t.Errorf("Array.Inspect() wrong. got=%q, want=%q", result.Inspect(), expected)
	}
}

func TestMapLiterals(t *testing.T) {
	input := `{"one": 1, "two": 2, "three": 3}`

	evaluated := testEval(input)
	result, ok := evaluated.(*object.Map)
	if !ok {
		t.Fatalf("Eval didn't return Map. got=%T (%+v)", evaluated, evaluated)
	}

	expected := map[object.HashKey]float64{
		(&object.String{Value: "one"}).HashKey():   1,
		(&object.String{Value: "two"}).HashKey():   2,
		(&object.String{Value: "three"}).HashKey(): 3,
	}

	if len(result.Pairs) != len(expected) {
		t.Fatalf("Map has wrong num of pairs. got=%d", len(result.Pairs))
	}

	for expectedKey, expectedValue := range expected {
		pair, ok := result.Pairs[expectedKey]
		if !ok {
			t.Errorf("no pair for given key in Pairs")
		}

		testNumberObject(t, pair.Value, expectedValue)
	}
}

func TestMapIndexExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{
			`{"foo": 5}["foo"]`,
			5,
		},
		{
			`{"foo": 5}["bar"]`,
			nil,
		},
		{
			`mut key = "foo"; {"foo": 5}[key]`,
			5,
		},
		{
			`{}["foo"]`,
			nil,
		},
		{
			`{5: 5}[5]`,
			5,
		},
		{
			`{true: 5}[true]`,
			5,
		},
		{
			`{false: 5}[false]`,
			5,
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			testNumberObject(t, evaluated, float64(integer))
		} else {
			testNullObject(t, evaluated)
		}
	}
}

func TestMapWithExpressionKeys(t *testing.T) {
	input := `
mut key = "name";
{key: "value"}[key]
`
	evaluated := testEval(input)
	result, ok := evaluated.(*object.String)
	if !ok {
		t.Fatalf("object is not String. got=%T (%+v)", evaluated, evaluated)
	}

	if result.Value != "value" {
		t.Errorf("String has wrong value. got=%q, want=%q", result.Value, "value")
	}
}

func TestMapUnhashableKey(t *testing.T) {
	input := `{[1, 2]: "array"}`

	evaluated := testEval(input)
	errObj, ok := evaluated.(*object.Error)
	if !ok {
		t.Fatalf("no error object returned. got=%T (%+v)", evaluated, evaluated)
	}

	expectedMessage := "unusable as hash key: ARRAY"
	if errObj.Message != expectedMessage {
		t.Errorf("wrong error message. expected=%q, got=%q", expectedMessage, errObj.Message)
	}
}

func TestErrorMessagesWithPosition(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"foobar",
			"line 1, column 1: identifier not found: foobar",
		},
		{
			"mut x = 10;\nfoobar",
			"line 2, column 1: identifier not found: foobar",
		},
		{
			"const x = 10;\nx = 20;",
			"line 2, column 3: cannot reassign to const variable: x",
		},
		{
			"y = 10;",
			"line 1, column 1: identifier not found: y",
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)

		errObj, ok := evaluated.(*object.Error)
		if !ok {
			t.Errorf("input %q: no error object returned. got=%T(%+v)", tt.input, evaluated, evaluated)
			continue
		}

		if errObj.Message != tt.expected {
			t.Errorf("input %q: wrong error message.\nexpected=%q\ngot=%q", tt.input, tt.expected, errObj.Message)
		}
	}
}

func TestBuiltinPop(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`pop([1, 2, 3])`, []int{1, 2}},
		{`pop([1])`, []int{}},
		{`pop([])`, nil},
		{`mut arr = [1, 2]; pop(arr); arr`, []int{1, 2}}, // 元の配列は変更されない
		{`pop(1)`, "argument to `pop` must be ARRAY, got NUMBER"},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)

		switch expected := tt.expected.(type) {
		case []int:
			arr, ok := evaluated.(*object.Array)
			if !ok {
				t.Errorf("object is not Array. got=%T (%+v)", evaluated, evaluated)
				continue
			}
			if len(arr.Elements) != len(expected) {
				t.Errorf("wrong num of elements. got=%d, want=%d", len(arr.Elements), len(expected))
				continue
			}
			for i, exp := range expected {
				testNumberObject(t, arr.Elements[i], float64(exp))
			}
		case nil:
			testNullObject(t, evaluated)
		case string:
			errObj, ok := evaluated.(*object.Error)
			if !ok {
				t.Errorf("object is not Error. got=%T (%+v)", evaluated, evaluated)
				continue
			}
			if errObj.Message != expected {
				t.Errorf("wrong error message. expected=%q, got=%q", expected, errObj.Message)
			}
		}
	}
}

func TestModuloWithFloats(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		// 整数同士の剰余（Req 1.2）
		{"10 % 3", 1},
		{"10 % 5", 0},
		{"0 % 5", 0},
		// 浮動小数点同士の剰余（Req 1.1）
		{"5.5 % 2.0", 1.5},
		{"10.5 % 3", 1.5},
		{"7.5 % 2.5", 0},
		// 負の数の剰余
		{"-10 % 3", -1},
		{"-5.5 % 2.0", -1.5},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testNumberObject(t, evaluated, tt.expected)
	}
}

func TestModuloByZero(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		// ゼロ除算エラー（Req 1.3）
		{"10 % 0", "division by zero"},
		{"5.5 % 0", "division by zero"},
		{"0 % 0", "division by zero"},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		errObj, ok := evaluated.(*object.Error)
		if !ok {
			t.Errorf("input %q: no error object returned. got=%T(%+v)", tt.input, evaluated, evaluated)
			continue
		}
		if errObj.Message != tt.expected {
			t.Errorf("input %q: wrong error message. expected=%q, got=%q", tt.input, tt.expected, errObj.Message)
		}
	}
}

func TestUnicodeStringLen(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{`len("hello")`, 5},
		{`len("あいう")`, 3},
		{`len("hello世界")`, 7},
		{`len("")`, 0},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testNumberObject(t, evaluated, float64(tt.expected))
	}
}

func TestUnicodeStringIndex(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`"hello"[0]`, "h"},
		{`"あいう"[0]`, "あ"},
		{`"あいう"[1]`, "い"},
		{`"あいう"[2]`, "う"},
		{`"hello世界"[5]`, "世"},
		{`"hello世界"[6]`, "界"},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		result, ok := evaluated.(*object.String)
		if !ok {
			t.Errorf("object is not String. got=%T (%+v)", evaluated, evaluated)
			continue
		}
		if result.Value != tt.expected {
			t.Errorf("wrong value. expected=%q, got=%q", tt.expected, result.Value)
		}
	}
}

func TestMapInspectOrder(t *testing.T) {
	// マップのInspect()が安定した順序を返すことを確認
	input := `{"b": 2, "a": 1, "c": 3}`

	// 複数回実行しても同じ結果になることを確認
	for i := 0; i < 5; i++ {
		evaluated := testEval(input)
		result, ok := evaluated.(*object.Map)
		if !ok {
			t.Fatalf("object is not Map. got=%T (%+v)", evaluated, evaluated)
		}

		inspect := result.Inspect()
		expected := "{a: 1, b: 2, c: 3}"
		if inspect != expected {
			t.Errorf("iteration %d: Map.Inspect() wrong. got=%q, want=%q", i, inspect, expected)
		}
	}
}

func TestArrayIndexAssignment(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		// 配列要素への代入
		{`mut arr = [1, 2, 3]; arr[0] = 10; arr[0]`, 10},
		{`mut arr = [1, 2, 3]; arr[1] = 20; arr[1]`, 20},
		{`mut arr = [1, 2, 3]; arr[2] = 30; arr[2]`, 30},
		// 配列全体の確認
		{`mut arr = [1, 2, 3]; arr[0] = 10; arr[1] + arr[2]`, 5},
		// 式をインデックスに使用
		{`mut arr = [1, 2, 3]; mut i = 1; arr[i] = 100; arr[1]`, 100},
		// 代入式の戻り値
		{`mut arr = [1, 2, 3]; arr[0] = 99`, 99},
		// 範囲外アクセスはエラー
		{`mut arr = [1, 2, 3]; arr[3] = 10`, "array index out of bounds: 3 (length: 3)"},
		{`mut arr = [1, 2, 3]; arr[-1] = 10`, "array index out of bounds: -1 (length: 3)"},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)

		switch expected := tt.expected.(type) {
		case int:
			testNumberObject(t, evaluated, float64(expected))
		case string:
			errObj, ok := evaluated.(*object.Error)
			if !ok {
				t.Errorf("input %q: object is not Error. got=%T (%+v)", tt.input, evaluated, evaluated)
				continue
			}
			if errObj.Message != expected {
				t.Errorf("input %q: wrong error message. expected=%q, got=%q", tt.input, expected, errObj.Message)
			}
		}
	}
}

func TestMapIndexAssignment(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		// マップ要素への代入（既存キー）
		{`mut map = {"a": 1, "b": 2}; map["a"] = 10; map["a"]`, 10},
		// マップ要素への代入（新規キー）
		{`mut map = {"a": 1}; map["b"] = 2; map["b"]`, 2},
		// 数値キー
		{`mut map = {1: "one"}; map[1] = "ONE"; map[1]`, "ONE"},
		// 代入式の戻り値
		{`mut map = {}; map["key"] = "value"`, "value"},
		// 式をキーに使用
		{`mut map = {}; mut key = "test"; map[key] = 123; map["test"]`, 123},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)

		switch expected := tt.expected.(type) {
		case int:
			testNumberObject(t, evaluated, float64(expected))
		case string:
			str, ok := evaluated.(*object.String)
			if ok {
				if str.Value != expected {
					t.Errorf("input %q: wrong string value. expected=%q, got=%q", tt.input, expected, str.Value)
				}
			} else {
				errObj, ok := evaluated.(*object.Error)
				if !ok {
					t.Errorf("input %q: object is not String or Error. got=%T (%+v)", tt.input, evaluated, evaluated)
					continue
				}
				if errObj.Message != expected {
					t.Errorf("input %q: wrong error message. expected=%q, got=%q", tt.input, expected, errObj.Message)
				}
			}
		}
	}
}

func TestIndexAssignmentToImmutableString(t *testing.T) {
	input := `mut s = "hello"; s[0] = "H"`

	evaluated := testEval(input)
	errObj, ok := evaluated.(*object.Error)
	if !ok {
		t.Fatalf("no error object returned. got=%T (%+v)", evaluated, evaluated)
	}

	expected := "index assignment not supported: STRING"
	if errObj.Message != expected {
		t.Errorf("wrong error message. expected=%q, got=%q", expected, errObj.Message)
	}
}

func TestIndexAssignmentToConstVariable(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		// const配列への要素代入
		{`const arr = [1, 2, 3]; arr[0] = 10`, "cannot modify const variable: arr"},
		// constマップへの要素代入
		{`const map = {"a": 1}; map["a"] = 10`, "cannot modify const variable: map"},
		// constマップへの新規キー追加
		{`const map = {"a": 1}; map["b"] = 2`, "cannot modify const variable: map"},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		errObj, ok := evaluated.(*object.Error)
		if !ok {
			t.Errorf("input %q: no error object returned. got=%T (%+v)", tt.input, evaluated, evaluated)
			continue
		}
		if errObj.Message != tt.expected {
			t.Errorf("input %q: wrong error message. expected=%q, got=%q", tt.input, tt.expected, errObj.Message)
		}
	}
}

func TestTryCatchStatement(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		// throw をキャッチ
		{
			`try { throw "error"; } catch (e) { e }`,
			"error",
		},
		// throw せずに正常終了
		{
			`try { 42 } catch (e) { e }`,
			42,
		},
		// throw した値を処理
		{
			`try { throw 123; } catch (e) { e + 1 }`,
			124,
		},
		// ネストした try/catch
		{
			`try {
				try {
					throw "inner";
				} catch (e) {
					throw "outer";
				}
			} catch (e) {
				e
			}`,
			"outer",
		},
		// 関数内での throw
		{
			`func throwError() => { throw "func error"; };
			try { throwError(); } catch (e) { e }`,
			"func error",
		},
		// catch 内で変数にアクセス
		{
			`mut result = "none";
			try { throw "caught"; } catch (e) { result = e; }
			result`,
			"caught",
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)

		switch expected := tt.expected.(type) {
		case int:
			testNumberObject(t, evaluated, float64(expected))
		case string:
			str, ok := evaluated.(*object.String)
			if !ok {
				t.Errorf("input %q: object is not String. got=%T (%+v)", tt.input, evaluated, evaluated)
				continue
			}
			if str.Value != expected {
				t.Errorf("input %q: wrong value. expected=%q, got=%q", tt.input, expected, str.Value)
			}
		}
	}
}

func TestUncaughtThrow(t *testing.T) {
	input := `throw "uncaught error"`

	evaluated := testEval(input)
	errObj, ok := evaluated.(*object.Error)
	if !ok {
		t.Fatalf("no error object returned. got=%T (%+v)", evaluated, evaluated)
	}

	expected := `uncaught exception: uncaught error`
	if errObj.Message != expected {
		t.Errorf("wrong error message. expected=%q, got=%q", expected, errObj.Message)
	}
}

func TestTryCatchWithBuiltinError(t *testing.T) {
	// 組み込み関数からのエラーもキャッチできる
	input := `try { 10 / 0 } catch (e) { "caught: " + e }`

	evaluated := testEval(input)
	str, ok := evaluated.(*object.String)
	if !ok {
		t.Fatalf("object is not String. got=%T (%+v)", evaluated, evaluated)
	}

	expected := "caught: division by zero"
	if str.Value != expected {
		t.Errorf("wrong value. expected=%q, got=%q", expected, str.Value)
	}
}

func TestBuiltinInt(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		// 数値から整数
		{`int(3.7)`, 3},
		{`int(3.2)`, 3},
		{`int(-3.7)`, -3},
		{`int(42)`, 42},
		{`int(0)`, 0},
		// 文字列から整数
		{`int("42")`, 42},
		{`int("3.14")`, 3},
		{`int("-10")`, -10},
		// 真偽値から整数
		{`int(true)`, 1},
		{`int(false)`, 0},
		// エラーケース
		{`int("hello")`, "cannot convert \"hello\" to int"},
		{`int(null)`, "cannot convert NULL to int"},
		{`int([1, 2])`, "cannot convert ARRAY to int"},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)

		switch expected := tt.expected.(type) {
		case int:
			testNumberObject(t, evaluated, float64(expected))
		case string:
			errObj, ok := evaluated.(*object.Error)
			if !ok {
				t.Errorf("input %q: object is not Error. got=%T (%+v)", tt.input, evaluated, evaluated)
				continue
			}
			if errObj.Message != expected {
				t.Errorf("input %q: wrong error message. expected=%q, got=%q", tt.input, expected, errObj.Message)
			}
		}
	}
}

func TestBuiltinFloat(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		// 数値はそのまま
		{`float(42)`, 42.0},
		{`float(3.14)`, 3.14},
		// 文字列から浮動小数点数
		{`float("3.14")`, 3.14},
		{`float("42")`, 42.0},
		{`float("-10.5")`, -10.5},
		// 真偽値から浮動小数点数
		{`float(true)`, 1.0},
		{`float(false)`, 0.0},
		// エラーケース
		{`float("hello")`, "cannot convert \"hello\" to float"},
		{`float(null)`, "cannot convert NULL to float"},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)

		switch expected := tt.expected.(type) {
		case float64:
			testNumberObject(t, evaluated, expected)
		case string:
			errObj, ok := evaluated.(*object.Error)
			if !ok {
				t.Errorf("input %q: object is not Error. got=%T (%+v)", tt.input, evaluated, evaluated)
				continue
			}
			if errObj.Message != expected {
				t.Errorf("input %q: wrong error message. expected=%q, got=%q", tt.input, expected, errObj.Message)
			}
		}
	}
}

func TestBuiltinString(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`string(42)`, "42"},
		{`string(3.14)`, "3.14"},
		{`string(true)`, "true"},
		{`string(false)`, "false"},
		{`string(null)`, "null"},
		{`string("hello")`, "hello"},
		{`string([1, 2, 3])`, "[1, 2, 3]"},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		str, ok := evaluated.(*object.String)
		if !ok {
			t.Errorf("input %q: object is not String. got=%T (%+v)", tt.input, evaluated, evaluated)
			continue
		}
		if str.Value != tt.expected {
			t.Errorf("input %q: wrong value. expected=%q, got=%q", tt.input, tt.expected, str.Value)
		}
	}
}

func TestBuiltinBool(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		// 数値
		{`bool(0)`, false},
		{`bool(1)`, true},
		{`bool(-1)`, true},
		{`bool(3.14)`, true},
		// 文字列
		{`bool("")`, false},
		{`bool("hello")`, true},
		// 真偽値
		{`bool(true)`, true},
		{`bool(false)`, false},
		// null
		{`bool(null)`, false},
		// 配列
		{`bool([])`, false},
		{`bool([1])`, true},
		{`bool([1, 2, 3])`, true},
		// マップ
		{`bool({})`, false},
		{`bool({"a": 1})`, true},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func TestBuiltinFileIO(t *testing.T) {
	// テスト用の一時ファイルを作成
	testFile := "test_file_io.txt"
	testContent := "Hello, Sugu!"

	// writeFile のテスト
	writeInput := fmt.Sprintf(`writeFile("%s", "%s")`, testFile, testContent)
	evaluated := testEval(writeInput)
	if evaluated != TRUE {
		t.Errorf("writeFile failed. got=%T (%+v)", evaluated, evaluated)
	}

	// fileExists のテスト（存在する場合）
	existsInput := fmt.Sprintf(`fileExists("%s")`, testFile)
	evaluated = testEval(existsInput)
	testBooleanObject(t, evaluated, true)

	// readFile のテスト
	readInput := fmt.Sprintf(`readFile("%s")`, testFile)
	evaluated = testEval(readInput)
	str, ok := evaluated.(*object.String)
	if !ok {
		t.Errorf("readFile failed. got=%T (%+v)", evaluated, evaluated)
	} else if str.Value != testContent {
		t.Errorf("readFile wrong content. expected=%q, got=%q", testContent, str.Value)
	}

	// appendFile のテスト
	appendContent := " Appended!"
	appendInput := fmt.Sprintf(`appendFile("%s", "%s")`, testFile, appendContent)
	evaluated = testEval(appendInput)
	if evaluated != TRUE {
		t.Errorf("appendFile failed. got=%T (%+v)", evaluated, evaluated)
	}

	// 追記後のファイル内容を確認
	evaluated = testEval(readInput)
	str, ok = evaluated.(*object.String)
	if !ok {
		t.Errorf("readFile after append failed. got=%T (%+v)", evaluated, evaluated)
	} else if str.Value != testContent+appendContent {
		t.Errorf("readFile after append wrong content. expected=%q, got=%q", testContent+appendContent, str.Value)
	}

	// ファイルを削除
	os.Remove(testFile)

	// fileExists のテスト（存在しない場合）
	evaluated = testEval(existsInput)
	testBooleanObject(t, evaluated, false)

	// readFile のエラーテスト（存在しないファイル）
	readInput = `readFile("nonexistent_file.txt")`
	evaluated = testEval(readInput)
	errObj, ok := evaluated.(*object.Error)
	if !ok {
		t.Errorf("readFile should return error for nonexistent file. got=%T (%+v)", evaluated, evaluated)
	} else if errObj.Message == "" {
		t.Errorf("readFile error message should not be empty")
	}
}

func TestDeleteBuiltin(t *testing.T) {
	// 正常系: 存在するキーの削除（Req 2.1）
	tests := []struct {
		input    string
		expected interface{} // bool or string(error)
	}{
		// キーが存在する場合は削除して true を返す
		{`mut m = {"a": 1, "b": 2}; delete(m, "a")`, true},
		// 数値キーの削除
		{`mut m = {1: "one", 2: "two"}; delete(m, 1)`, true},
		// ブーリアンキーの削除
		{`mut m = {true: "yes"}; delete(m, true)`, true},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected.(bool))
	}

	// 削除後にキーが存在しないことを確認
	input := `mut m = {"a": 1, "b": 2}; delete(m, "a"); m["a"]`
	evaluated := testEval(input)
	testNullObject(t, evaluated)

	// 削除後に他のキーが残っていることを確認
	input = `mut m = {"a": 1, "b": 2}; delete(m, "a"); m["b"]`
	evaluated = testEval(input)
	testNumberObject(t, evaluated, 2)
}

func TestDeleteBuiltinNotFound(t *testing.T) {
	// 存在しないキーの場合は false を返す（Req 2.2）
	tests := []struct {
		input    string
		expected bool
	}{
		{`mut m = {"a": 1}; delete(m, "z")`, false},
		{`mut m = {}; delete(m, "a")`, false},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func TestDeleteBuiltinErrors(t *testing.T) {
	// 型エラー（Req 2.4）
	tests := []struct {
		input    string
		expected string
	}{
		// 第1引数がマップ以外
		{`delete([1,2,3], 0)`, "argument to `delete` must be MAP, got ARRAY"},
		{`delete("hello", 0)`, "argument to `delete` must be MAP, got STRING"},
		{`delete(42, 0)`, "argument to `delete` must be MAP, got NUMBER"},
		// 引数の数が不正
		{`delete({"a": 1})`, "wrong number of arguments. got=1, want=2"},
		{`delete({"a": 1}, "a", "b")`, "wrong number of arguments. got=3, want=2"},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		errObj, ok := evaluated.(*object.Error)
		if !ok {
			t.Errorf("input %q: no error object returned. got=%T(%+v)", tt.input, evaluated, evaluated)
			continue
		}
		if errObj.Message != tt.expected {
			t.Errorf("input %q: wrong error message. expected=%q, got=%q", tt.input, tt.expected, errObj.Message)
		}
	}
}

func TestSplitBuiltin(t *testing.T) {
	// 正常系（Req 3.1）
	tests := []struct {
		input    string
		expected []string
	}{
		{`split("a,b,c", ",")`, []string{"a", "b", "c"}},
		{`split("hello world", " ")`, []string{"hello", "world"}},
		{`split("abc", "")`, []string{"a", "b", "c"}},
		{`split("", ",")`, []string{""}},
		// マルチバイト文字
		{`split("あ,い,う", ",")`, []string{"あ", "い", "う"}},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		arr, ok := evaluated.(*object.Array)
		if !ok {
			t.Errorf("input %q: object is not Array. got=%T (%+v)", tt.input, evaluated, evaluated)
			continue
		}
		if len(arr.Elements) != len(tt.expected) {
			t.Errorf("input %q: wrong num of elements. got=%d, want=%d", tt.input, len(arr.Elements), len(tt.expected))
			continue
		}
		for i, exp := range tt.expected {
			str, ok := arr.Elements[i].(*object.String)
			if !ok {
				t.Errorf("input %q: element %d is not String. got=%T", tt.input, i, arr.Elements[i])
				continue
			}
			if str.Value != exp {
				t.Errorf("input %q: element %d wrong. expected=%q, got=%q", tt.input, i, exp, str.Value)
			}
		}
	}

	// 型エラー（Req 3.9）
	errTests := []struct {
		input    string
		expected string
	}{
		{`split(123, ",")`, "argument to `split` must be STRING, got NUMBER"},
		{`split("a", 1)`, "second argument to `split` must be STRING, got NUMBER"},
		{`split("a")`, "wrong number of arguments. got=1, want=2"},
	}
	for _, tt := range errTests {
		evaluated := testEval(tt.input)
		errObj, ok := evaluated.(*object.Error)
		if !ok {
			t.Errorf("input %q: no error. got=%T(%+v)", tt.input, evaluated, evaluated)
			continue
		}
		if errObj.Message != tt.expected {
			t.Errorf("input %q: wrong error. expected=%q, got=%q", tt.input, tt.expected, errObj.Message)
		}
	}
}

func TestJoinBuiltin(t *testing.T) {
	// 正常系（Req 3.2）
	tests := []struct {
		input    string
		expected string
	}{
		{`join(["a", "b", "c"], ",")`, "a,b,c"},
		{`join(["hello", "world"], " ")`, "hello world"},
		{`join(["x"], "-")`, "x"},
		{`join([], ",")`, ""},
		// マルチバイト文字
		{`join(["あ", "い", "う"], "・")`, "あ・い・う"},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		str, ok := evaluated.(*object.String)
		if !ok {
			t.Errorf("input %q: object is not String. got=%T (%+v)", tt.input, evaluated, evaluated)
			continue
		}
		if str.Value != tt.expected {
			t.Errorf("input %q: wrong value. expected=%q, got=%q", tt.input, tt.expected, str.Value)
		}
	}

	// 型エラー（Req 3.9）
	errTests := []struct {
		input    string
		expected string
	}{
		{`join("abc", ",")`, "argument to `join` must be ARRAY, got STRING"},
		{`join(["a"], 1)`, "second argument to `join` must be STRING, got NUMBER"},
		{`join(["a"])`, "wrong number of arguments. got=1, want=2"},
	}
	for _, tt := range errTests {
		evaluated := testEval(tt.input)
		errObj, ok := evaluated.(*object.Error)
		if !ok {
			t.Errorf("input %q: no error. got=%T(%+v)", tt.input, evaluated, evaluated)
			continue
		}
		if errObj.Message != tt.expected {
			t.Errorf("input %q: wrong error. expected=%q, got=%q", tt.input, tt.expected, errObj.Message)
		}
	}
}

func TestTrimBuiltin(t *testing.T) {
	// 正常系（Req 3.3）
	tests := []struct {
		input    string
		expected string
	}{
		{`trim("  hello  ")`, "hello"},
		{`trim("hello")`, "hello"},
		{`trim("  ")`, ""},
		{`trim("")`, ""},
		{`trim("\thello\n")`, "hello"},
		// マルチバイト文字
		{`trim("  こんにちは  ")`, "こんにちは"},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		str, ok := evaluated.(*object.String)
		if !ok {
			t.Errorf("input %q: object is not String. got=%T (%+v)", tt.input, evaluated, evaluated)
			continue
		}
		if str.Value != tt.expected {
			t.Errorf("input %q: wrong value. expected=%q, got=%q", tt.input, tt.expected, str.Value)
		}
	}

	// 型エラー
	errTests := []struct {
		input    string
		expected string
	}{
		{`trim(123)`, "argument to `trim` must be STRING, got NUMBER"},
		{`trim()`, "wrong number of arguments. got=0, want=1"},
	}
	for _, tt := range errTests {
		evaluated := testEval(tt.input)
		errObj, ok := evaluated.(*object.Error)
		if !ok {
			t.Errorf("input %q: no error. got=%T(%+v)", tt.input, evaluated, evaluated)
			continue
		}
		if errObj.Message != tt.expected {
			t.Errorf("input %q: wrong error. expected=%q, got=%q", tt.input, tt.expected, errObj.Message)
		}
	}
}

func TestReplaceBuiltin(t *testing.T) {
	// 正常系（Req 3.4）
	tests := []struct {
		input    string
		expected string
	}{
		{`replace("hello world", "world", "Go")`, "hello Go"},
		{`replace("aaa", "a", "b")`, "bbb"},
		{`replace("hello", "xyz", "abc")`, "hello"},
		{`replace("", "a", "b")`, ""},
		// マルチバイト文字
		{`replace("こんにちは世界", "世界", "Sugu")`, "こんにちはSugu"},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		str, ok := evaluated.(*object.String)
		if !ok {
			t.Errorf("input %q: object is not String. got=%T (%+v)", tt.input, evaluated, evaluated)
			continue
		}
		if str.Value != tt.expected {
			t.Errorf("input %q: wrong value. expected=%q, got=%q", tt.input, tt.expected, str.Value)
		}
	}

	// 型エラー（Req 3.9）
	errTests := []struct {
		input    string
		expected string
	}{
		{`replace(123, "a", "b")`, "argument to `replace` must be STRING, got NUMBER"},
		{`replace("a", 1, "b")`, "second argument to `replace` must be STRING, got NUMBER"},
		{`replace("a", "b", 1)`, "third argument to `replace` must be STRING, got NUMBER"},
		{`replace("a", "b")`, "wrong number of arguments. got=2, want=3"},
	}
	for _, tt := range errTests {
		evaluated := testEval(tt.input)
		errObj, ok := evaluated.(*object.Error)
		if !ok {
			t.Errorf("input %q: no error. got=%T(%+v)", tt.input, evaluated, evaluated)
			continue
		}
		if errObj.Message != tt.expected {
			t.Errorf("input %q: wrong error. expected=%q, got=%q", tt.input, tt.expected, errObj.Message)
		}
	}
}
