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

func TestEmptyInput(t *testing.T) {
	l := New("")
	tok := l.NextToken()

	if tok.Type != token.EOF {
		t.Errorf("empty input should return EOF, got=%s", tok.Type)
	}
}

func TestUnterminatedString(t *testing.T) {
	l := New(`"hello`)
	tok := l.NextToken()

	// 未終端の文字列はSTRINGトークンとして扱われる（内容は"hello"まで）
	if tok.Type != token.STRING {
		t.Errorf("unterminated string should return STRING, got=%s", tok.Type)
	}
	if tok.Literal != "hello" {
		t.Errorf("unterminated string literal wrong. got=%q, want=%q", tok.Literal, "hello")
	}
}

func TestUnterminatedMultilineComment(t *testing.T) {
	l := New(`mut x = 10; //-- 未終端のコメント`)

	// mut
	tok := l.NextToken()
	if tok.Type != token.MUT {
		t.Errorf("expected MUT, got=%s", tok.Type)
	}

	// x
	tok = l.NextToken()
	if tok.Type != token.IDENT {
		t.Errorf("expected IDENT, got=%s", tok.Type)
	}

	// =
	tok = l.NextToken()
	if tok.Type != token.ASSIGN {
		t.Errorf("expected ASSIGN, got=%s", tok.Type)
	}

	// 10
	tok = l.NextToken()
	if tok.Type != token.NUMBER {
		t.Errorf("expected NUMBER, got=%s", tok.Type)
	}

	// ;
	tok = l.NextToken()
	if tok.Type != token.SEMICOLON {
		t.Errorf("expected SEMICOLON, got=%s", tok.Type)
	}

	// コメントはスキップされてEOFになるはず
	tok = l.NextToken()
	if tok.Type != token.EOF {
		t.Errorf("expected EOF after unterminated comment, got=%s", tok.Type)
	}
}

func TestNegativeNumberTokens(t *testing.T) {
	// 負の数は MINUS + NUMBER として解析される
	l := New(`-10`)

	tok := l.NextToken()
	if tok.Type != token.MINUS {
		t.Errorf("expected MINUS, got=%s", tok.Type)
	}
	if tok.Literal != "-" {
		t.Errorf("expected '-', got=%s", tok.Literal)
	}

	tok = l.NextToken()
	if tok.Type != token.NUMBER {
		t.Errorf("expected NUMBER, got=%s", tok.Type)
	}
	if tok.Literal != "10" {
		t.Errorf("expected '10', got=%s", tok.Literal)
	}
}

func TestWhitespaceOnly(t *testing.T) {
	l := New("   \t\n\r  ")
	tok := l.NextToken()

	if tok.Type != token.EOF {
		t.Errorf("whitespace-only input should return EOF, got=%s", tok.Type)
	}
}

func TestIllegalCharacters(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"@", "@"},
		{"#", "#"},
		{"$", "$"},
		{"~", "~"},
	}

	for _, tt := range tests {
		l := New(tt.input)
		tok := l.NextToken()

		if tok.Type != token.ILLEGAL {
			t.Errorf("illegal character %q should return ILLEGAL, got=%s", tt.input, tok.Type)
		}
		if tok.Literal != tt.expected {
			t.Errorf("illegal literal wrong. got=%s, want=%s", tok.Literal, tt.expected)
		}
	}
}

func TestUnknownEscapeSequence(t *testing.T) {
	tests := []struct {
		input           string
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{`"hello\xworld"`, token.ILLEGAL, "unknown escape sequence: \\x"},
		{`"test\avalue"`, token.ILLEGAL, "unknown escape sequence: \\a"},
		{`"foo\bbar"`, token.ILLEGAL, "unknown escape sequence: \\b"},
	}

	for _, tt := range tests {
		l := New(tt.input)
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Errorf("input=%q: expected type %s, got=%s", tt.input, tt.expectedType, tok.Type)
		}
		if tok.Literal != tt.expectedLiteral {
			t.Errorf("input=%q: expected literal %q, got=%q", tt.input, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestEscapeSequenceAtEOF(t *testing.T) {
	// バックスラッシュ直後にEOF
	l := New(`"hello\`)
	tok := l.NextToken()

	if tok.Type != token.ILLEGAL {
		t.Errorf("expected ILLEGAL for escape at EOF, got=%s", tok.Type)
	}
	if tok.Literal != "unexpected end of string after \\" {
		t.Errorf("unexpected literal: %q", tok.Literal)
	}
}
