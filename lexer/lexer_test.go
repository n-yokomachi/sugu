package lexer

import (
	"testing"

	"sugu/token"
)

func TestNextToken(t *testing.T) {
	input := `mut x = 10;
const PI = 3.14;

func add(a, b) => {
	return a + b;
}

if (x > 0) {
	x = x - 1;
}
`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.MUT, "mut"},
		{token.IDENT, "x"},
		{token.ASSIGN, "="},
		{token.NUMBER, "10"},
		{token.SEMICOLON, ";"},
		{token.CONST, "const"},
		{token.IDENT, "PI"},
		{token.ASSIGN, "="},
		{token.NUMBER, "3.14"},
		{token.SEMICOLON, ";"},
		{token.FUNC, "func"},
		{token.IDENT, "add"},
		{token.LPAREN, "("},
		{token.IDENT, "a"},
		{token.COMMA, ","},
		{token.IDENT, "b"},
		{token.RPAREN, ")"},
		{token.ARROW, "=>"},
		{token.LBRACE, "{"},
		{token.RETURN, "return"},
		{token.IDENT, "a"},
		{token.PLUS, "+"},
		{token.IDENT, "b"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.IF, "if"},
		{token.LPAREN, "("},
		{token.IDENT, "x"},
		{token.GT, ">"},
		{token.NUMBER, "0"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.IDENT, "x"},
		{token.ASSIGN, "="},
		{token.IDENT, "x"},
		{token.MINUS, "-"},
		{token.NUMBER, "1"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.EOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestOperators(t *testing.T) {
	input := `+ - * / % == != < > <= >= && || !`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.PLUS, "+"},
		{token.MINUS, "-"},
		{token.ASTERISK, "*"},
		{token.SLASH, "/"},
		{token.PERCENT, "%"},
		{token.EQ, "=="},
		{token.NOT_EQ, "!="},
		{token.LT, "<"},
		{token.GT, ">"},
		{token.LT_EQ, "<="},
		{token.GT_EQ, ">="},
		{token.AND, "&&"},
		{token.OR, "||"},
		{token.BANG, "!"},
		{token.EOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestString(t *testing.T) {
	input := `"hello" "world"`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.STRING, "hello"},
		{token.STRING, "world"},
		{token.EOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestKeywords(t *testing.T) {
	input := `mut const func return if else switch case default while for break continue true false null`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.MUT, "mut"},
		{token.CONST, "const"},
		{token.FUNC, "func"},
		{token.RETURN, "return"},
		{token.IF, "if"},
		{token.ELSE, "else"},
		{token.SWITCH, "switch"},
		{token.CASE, "case"},
		{token.DEFAULT, "default"},
		{token.WHILE, "while"},
		{token.FOR, "for"},
		{token.BREAK, "break"},
		{token.CONTINUE, "continue"},
		{token.TRUE, "true"},
		{token.FALSE, "false"},
		{token.NULL, "null"},
		{token.EOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestComments(t *testing.T) {
	input := `mut x = 10; // これはコメント
mut y = 20;
//-- これは
複数行コメント --//
mut z = 30;`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.MUT, "mut"},
		{token.IDENT, "x"},
		{token.ASSIGN, "="},
		{token.NUMBER, "10"},
		{token.SEMICOLON, ";"},
		{token.MUT, "mut"},
		{token.IDENT, "y"},
		{token.ASSIGN, "="},
		{token.NUMBER, "20"},
		{token.SEMICOLON, ";"},
		{token.MUT, "mut"},
		{token.IDENT, "z"},
		{token.ASSIGN, "="},
		{token.NUMBER, "30"},
		{token.SEMICOLON, ";"},
		{token.EOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestStringEscapeSequences(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`"hello\nworld"`, "hello\nworld"},
		{`"hello\tworld"`, "hello\tworld"},
		{`"hello\rworld"`, "hello\rworld"},
		{`"hello\"world"`, "hello\"world"},
		{`"hello\\world"`, "hello\\world"},
		{`"line1\nline2\nline3"`, "line1\nline2\nline3"},
		{`"tab\there"`, "tab\there"},
		{`"quote: \"test\""`, "quote: \"test\""},
	}

	for _, tt := range tests {
		l := New(tt.input)
		tok := l.NextToken()

		if tok.Type != token.STRING {
			t.Errorf("token type wrong. expected=STRING, got=%s", tok.Type)
		}

		if tok.Literal != tt.expected {
			t.Errorf("string literal wrong. input=%s, expected=%q, got=%q",
				tt.input, tt.expected, tok.Literal)
		}
	}
}
