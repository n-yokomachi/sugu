package object

import (
	"fmt"
	"math"
	"strings"
	"sugu/ast"
)

// ObjectType はオブジェクトの型を表す
type ObjectType string

const (
	NUMBER_OBJ   ObjectType = "NUMBER"
	STRING_OBJ   ObjectType = "STRING"
	BOOLEAN_OBJ  ObjectType = "BOOLEAN"
	NULL_OBJ     ObjectType = "NULL"
	RETURN_OBJ   ObjectType = "RETURN"
	ERROR_OBJ    ObjectType = "ERROR"
	FUNCTION_OBJ ObjectType = "FUNCTION"
	BUILTIN_OBJ  ObjectType = "BUILTIN"
	ARRAY_OBJ    ObjectType = "ARRAY"
)

// Object はすべてのオブジェクトの基底インターフェース
type Object interface {
	Type() ObjectType
	Inspect() string
}

// Number は数値を表す
type Number struct {
	Value float64
}

func (n *Number) Type() ObjectType { return NUMBER_OBJ }
func (n *Number) Inspect() string {
	// 整数の場合は小数点以下を表示しない
	// math.Truncを使用して正確に判定し、int64の範囲内かもチェック
	if math.Trunc(n.Value) == n.Value &&
		n.Value >= math.MinInt64 && n.Value <= math.MaxInt64 {
		return fmt.Sprintf("%d", int64(n.Value))
	}
	return fmt.Sprintf("%g", n.Value)
}

// String は文字列を表す
type String struct {
	Value string
}

func (s *String) Type() ObjectType { return STRING_OBJ }
func (s *String) Inspect() string  { return s.Value }

// Boolean は真偽値を表す
type Boolean struct {
	Value bool
}

func (b *Boolean) Type() ObjectType { return BOOLEAN_OBJ }
func (b *Boolean) Inspect() string {
	if b.Value {
		return "true"
	}
	return "false"
}

// Null はnull値を表す
type Null struct{}

func (n *Null) Type() ObjectType { return NULL_OBJ }
func (n *Null) Inspect() string  { return "null" }

// ReturnValue はreturn文の戻り値をラップする
type ReturnValue struct {
	Value Object
}

func (rv *ReturnValue) Type() ObjectType { return RETURN_OBJ }
func (rv *ReturnValue) Inspect() string  { return rv.Value.Inspect() }

// Error はエラーを表す
type Error struct {
	Message string
}

func (e *Error) Type() ObjectType { return ERROR_OBJ }
func (e *Error) Inspect() string  { return "ERROR: " + e.Message }

// Function は関数を表す
type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment
	Name       string
}

func (f *Function) Type() ObjectType { return FUNCTION_OBJ }
func (f *Function) Inspect() string {
	params := []string{}
	for _, p := range f.Parameters {
		params = append(params, p.String())
	}

	if f.Name != "" {
		return fmt.Sprintf("func %s(%s) => { ... }", f.Name, strings.Join(params, ", "))
	}
	return fmt.Sprintf("func(%s) => { ... }", strings.Join(params, ", "))
}

// BuiltinFunction は組み込み関数の型
type BuiltinFunction func(args ...Object) Object

// Builtin は組み込み関数を表す
type Builtin struct {
	Fn BuiltinFunction
}

func (b *Builtin) Type() ObjectType { return BUILTIN_OBJ }
func (b *Builtin) Inspect() string  { return "builtin function" }

// Array は配列を表す
type Array struct {
	Elements []Object
}

func (a *Array) Type() ObjectType { return ARRAY_OBJ }
func (a *Array) Inspect() string {
	elements := []string{}
	for _, e := range a.Elements {
		elements = append(elements, e.Inspect())
	}
	return "[" + strings.Join(elements, ", ") + "]"
}
