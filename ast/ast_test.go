package ast

import (
	"testing"

	"sugu/token"
)

func TestString(t *testing.T) {
	program := &Program{
		Statements: []Statement{
			&VariableStatement{
				Token: token.Token{Type: token.MUT, Literal: "mut"},
				Name: &Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "x"},
					Value: "x",
				},
				Value: &NumberLiteral{
					Token: token.Token{Type: token.NUMBER, Literal: "10"},
					Value: "10",
				},
			},
		},
	}

	if program.String() != "mut x = 10;" {
		t.Errorf("program.String() wrong. got=%q", program.String())
	}
}

func TestIdentifier(t *testing.T) {
	ident := &Identifier{
		Token: token.Token{Type: token.IDENT, Literal: "myVar"},
		Value: "myVar",
	}

	if ident.String() != "myVar" {
		t.Errorf("ident.String() wrong. got=%q", ident.String())
	}

	if ident.TokenLiteral() != "myVar" {
		t.Errorf("ident.TokenLiteral() wrong. got=%q", ident.TokenLiteral())
	}
}

func TestNumberLiteral(t *testing.T) {
	num := &NumberLiteral{
		Token: token.Token{Type: token.NUMBER, Literal: "42"},
		Value: "42",
	}

	if num.String() != "42" {
		t.Errorf("num.String() wrong. got=%q", num.String())
	}
}

func TestStringLiteral(t *testing.T) {
	str := &StringLiteral{
		Token: token.Token{Type: token.STRING, Literal: "hello"},
		Value: "hello",
	}

	if str.String() != "\"hello\"" {
		t.Errorf("str.String() wrong. got=%q", str.String())
	}
}

func TestBooleanLiteral(t *testing.T) {
	trueLit := &BooleanLiteral{
		Token: token.Token{Type: token.TRUE, Literal: "true"},
		Value: true,
	}

	if trueLit.String() != "true" {
		t.Errorf("trueLit.String() wrong. got=%q", trueLit.String())
	}

	falseLit := &BooleanLiteral{
		Token: token.Token{Type: token.FALSE, Literal: "false"},
		Value: false,
	}

	if falseLit.String() != "false" {
		t.Errorf("falseLit.String() wrong. got=%q", falseLit.String())
	}
}

func TestNullLiteral(t *testing.T) {
	null := &NullLiteral{
		Token: token.Token{Type: token.NULL, Literal: "null"},
	}

	if null.String() != "null" {
		t.Errorf("null.String() wrong. got=%q", null.String())
	}
}

func TestPrefixExpression(t *testing.T) {
	prefix := &PrefixExpression{
		Token:    token.Token{Type: token.BANG, Literal: "!"},
		Operator: "!",
		Right: &BooleanLiteral{
			Token: token.Token{Type: token.TRUE, Literal: "true"},
			Value: true,
		},
	}

	if prefix.String() != "(!true)" {
		t.Errorf("prefix.String() wrong. got=%q", prefix.String())
	}
}

func TestInfixExpression(t *testing.T) {
	infix := &InfixExpression{
		Token: token.Token{Type: token.PLUS, Literal: "+"},
		Left: &NumberLiteral{
			Token: token.Token{Type: token.NUMBER, Literal: "5"},
			Value: "5",
		},
		Operator: "+",
		Right: &NumberLiteral{
			Token: token.Token{Type: token.NUMBER, Literal: "10"},
			Value: "10",
		},
	}

	if infix.String() != "(5 + 10)" {
		t.Errorf("infix.String() wrong. got=%q", infix.String())
	}
}

func TestCallExpression(t *testing.T) {
	call := &CallExpression{
		Token: token.Token{Type: token.LPAREN, Literal: "("},
		Function: &Identifier{
			Token: token.Token{Type: token.IDENT, Literal: "add"},
			Value: "add",
		},
		Arguments: []Expression{
			&NumberLiteral{
				Token: token.Token{Type: token.NUMBER, Literal: "1"},
				Value: "1",
			},
			&NumberLiteral{
				Token: token.Token{Type: token.NUMBER, Literal: "2"},
				Value: "2",
			},
		},
	}

	if call.String() != "add(1, 2)" {
		t.Errorf("call.String() wrong. got=%q", call.String())
	}
}

func TestReturnStatement(t *testing.T) {
	ret := &ReturnStatement{
		Token: token.Token{Type: token.RETURN, Literal: "return"},
		ReturnValue: &NumberLiteral{
			Token: token.Token{Type: token.NUMBER, Literal: "42"},
			Value: "42",
		},
	}

	if ret.String() != "return 42;" {
		t.Errorf("ret.String() wrong. got=%q", ret.String())
	}
}

func TestIfStatement(t *testing.T) {
	ifStmt := &IfStatement{
		Token: token.Token{Type: token.IF, Literal: "if"},
		Condition: &BooleanLiteral{
			Token: token.Token{Type: token.TRUE, Literal: "true"},
			Value: true,
		},
		Consequence: &BlockStatement{
			Token:      token.Token{Type: token.LBRACE, Literal: "{"},
			Statements: []Statement{},
		},
		Alternative: nil,
	}

	if ifStmt.String() != "if (true) {  }" {
		t.Errorf("ifStmt.String() wrong. got=%q", ifStmt.String())
	}
}

func TestFunctionLiteral(t *testing.T) {
	fn := &FunctionLiteral{
		Token: token.Token{Type: token.FUNC, Literal: "func"},
		Name: &Identifier{
			Token: token.Token{Type: token.IDENT, Literal: "add"},
			Value: "add",
		},
		Parameters: []*Identifier{
			{
				Token: token.Token{Type: token.IDENT, Literal: "a"},
				Value: "a",
			},
			{
				Token: token.Token{Type: token.IDENT, Literal: "b"},
				Value: "b",
			},
		},
		Body: &BlockStatement{
			Token:      token.Token{Type: token.LBRACE, Literal: "{"},
			Statements: []Statement{},
		},
	}

	if fn.String() != "func add(a, b) => {  }" {
		t.Errorf("fn.String() wrong. got=%q", fn.String())
	}
}

func TestWhileStatement(t *testing.T) {
	while := &WhileStatement{
		Token: token.Token{Type: token.WHILE, Literal: "while"},
		Condition: &BooleanLiteral{
			Token: token.Token{Type: token.TRUE, Literal: "true"},
			Value: true,
		},
		Body: &BlockStatement{
			Token:      token.Token{Type: token.LBRACE, Literal: "{"},
			Statements: []Statement{},
		},
	}

	if while.String() != "while (true) {  }" {
		t.Errorf("while.String() wrong. got=%q", while.String())
	}
}

func TestBreakStatement(t *testing.T) {
	brk := &BreakStatement{
		Token: token.Token{Type: token.BREAK, Literal: "break"},
	}

	if brk.String() != "break;" {
		t.Errorf("brk.String() wrong. got=%q", brk.String())
	}
}

func TestContinueStatement(t *testing.T) {
	cont := &ContinueStatement{
		Token: token.Token{Type: token.CONTINUE, Literal: "continue"},
	}

	if cont.String() != "continue;" {
		t.Errorf("cont.String() wrong. got=%q", cont.String())
	}
}
