package parser

import (
	"sugu/ast"
	"sugu/lexer"
	"testing"
)

func TestVariableStatements(t *testing.T) {
	input := `
mut x = 5;
const y = 10;
`

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 2 {
		t.Fatalf("program.Statements does not contain 2 statements. got=%d",
			len(program.Statements))
	}

	tests := []struct {
		expectedIdentifier string
	}{
		{"x"},
		{"y"},
	}

	for i, tt := range tests {
		stmt := program.Statements[i]
		if !testVariableStatement(t, stmt, tt.expectedIdentifier) {
			return
		}
	}
}

func testVariableStatement(t *testing.T, s ast.Statement, name string) bool {
	if s.TokenLiteral() != "mut" && s.TokenLiteral() != "const" {
		t.Errorf("s.TokenLiteral not 'mut' or 'const'. got=%q", s.TokenLiteral())
		return false
	}

	varStmt, ok := s.(*ast.VariableStatement)
	if !ok {
		t.Errorf("s not *ast.VariableStatement. got=%T", s)
		return false
	}

	if varStmt.Name.Value != name {
		t.Errorf("varStmt.Name.Value not '%s'. got=%s", name, varStmt.Name.Value)
		return false
	}

	if varStmt.Name.TokenLiteral() != name {
		t.Errorf("varStmt.Name.TokenLiteral() not '%s'. got=%s",
			name, varStmt.Name.TokenLiteral())
		return false
	}

	return true
}

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(errors))
	for _, msg := range errors {
		t.Errorf("parser error: %q", msg)
	}
	t.FailNow()
}

func TestReturnStatements(t *testing.T) {
	input := `
return 5;
return 10;
`

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 2 {
		t.Fatalf("program.Statements does not contain 2 statements. got=%d",
			len(program.Statements))
	}

	for _, stmt := range program.Statements {
		returnStmt, ok := stmt.(*ast.ReturnStatement)
		if !ok {
			t.Errorf("stmt not *ast.ReturnStatement. got=%T", stmt)
			continue
		}
		if returnStmt.TokenLiteral() != "return" {
			t.Errorf("returnStmt.TokenLiteral not 'return', got %q",
				returnStmt.TokenLiteral())
		}
	}
}

func TestIdentifierExpression(t *testing.T) {
	input := "foobar;"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d",
			len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	ident, ok := stmt.Expression.(*ast.Identifier)
	if !ok {
		t.Fatalf("exp not *ast.Identifier. got=%T", stmt.Expression)
	}
	if ident.Value != "foobar" {
		t.Errorf("ident.Value not %s. got=%s", "foobar", ident.Value)
	}
	if ident.TokenLiteral() != "foobar" {
		t.Errorf("ident.TokenLiteral not %s. got=%s", "foobar",
			ident.TokenLiteral())
	}
}

func TestNumberLiteralExpression(t *testing.T) {
	input := "5;"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d",
			len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	literal, ok := stmt.Expression.(*ast.NumberLiteral)
	if !ok {
		t.Fatalf("exp not *ast.NumberLiteral. got=%T", stmt.Expression)
	}
	if literal.Value != "5" {
		t.Errorf("literal.Value not %s. got=%s", "5", literal.Value)
	}
}

func TestPrefixExpressions(t *testing.T) {
	prefixTests := []struct {
		input    string
		operator string
		value    string
	}{
		{"!5;", "!", "5"},
		{"-15;", "-", "15"},
		{"!true;", "!", "true"},
		{"!false;", "!", "false"},
	}

	for _, tt := range prefixTests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
				1, len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
				program.Statements[0])
		}

		exp, ok := stmt.Expression.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("stmt is not ast.PrefixExpression. got=%T", stmt.Expression)
		}
		if exp.Operator != tt.operator {
			t.Fatalf("exp.Operator is not '%s'. got=%s",
				tt.operator, exp.Operator)
		}
	}
}

func TestInfixExpressions(t *testing.T) {
	infixTests := []struct {
		input      string
		leftValue  string
		operator   string
		rightValue string
	}{
		{"5 + 5;", "5", "+", "5"},
		{"5 - 5;", "5", "-", "5"},
		{"5 * 5;", "5", "*", "5"},
		{"5 / 5;", "5", "/", "5"},
		{"5 > 5;", "5", ">", "5"},
		{"5 < 5;", "5", "<", "5"},
		{"5 == 5;", "5", "==", "5"},
		{"5 != 5;", "5", "!=", "5"},
	}

	for _, tt := range infixTests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
				1, len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
				program.Statements[0])
		}

		exp, ok := stmt.Expression.(*ast.InfixExpression)
		if !ok {
			t.Fatalf("exp is not ast.InfixExpression. got=%T", stmt.Expression)
		}

		if exp.Operator != tt.operator {
			t.Errorf("exp.Operator is not '%s'. got=%s",
				tt.operator, exp.Operator)
		}
	}
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"-a * b",
			"((-a) * b)",
		},
		{
			"!-a",
			"(!(-a))",
		},
		{
			"a + b + c",
			"((a + b) + c)",
		},
		{
			"a + b - c",
			"((a + b) - c)",
		},
		{
			"a * b * c",
			"((a * b) * c)",
		},
		{
			"a * b / c",
			"((a * b) / c)",
		},
		{
			"a + b / c",
			"(a + (b / c))",
		},
		{
			"a + b * c + d / e - f",
			"(((a + (b * c)) + (d / e)) - f)",
		},
		{
			"3 + 4; -5 * 5",
			"(3 + 4)((-5) * 5)",
		},
		{
			"5 > 4 == 3 < 4",
			"((5 > 4) == (3 < 4))",
		},
		{
			"5 < 4 != 3 > 4",
			"((5 < 4) != (3 > 4))",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
		{
			"true",
			"true",
		},
		{
			"false",
			"false",
		},
		{
			"3 > 5 == false",
			"((3 > 5) == false)",
		},
		{
			"3 < 5 == true",
			"((3 < 5) == true)",
		},
		{
			"1 + (2 + 3) + 4",
			"((1 + (2 + 3)) + 4)",
		},
		{
			"(5 + 5) * 2",
			"((5 + 5) * 2)",
		},
		{
			"2 / (5 + 5)",
			"(2 / (5 + 5))",
		},
		{
			"-(5 + 5)",
			"(-(5 + 5))",
		},
		{
			"!(true == true)",
			"(!(true == true))",
		},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		actual := program.String()
		if actual != tt.expected {
			t.Errorf("expected=%q, got=%q", tt.expected, actual)
		}
	}
}

func TestIfExpression(t *testing.T) {
	input := `if (x < y) { x }`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
			1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.IfStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.IfStatement. got=%T",
			program.Statements[0])
	}

	if stmt.Condition == nil {
		t.Fatalf("stmt.Condition is nil")
	}

	if stmt.Consequence == nil {
		t.Fatalf("stmt.Consequence is nil")
	}
}

func TestFunctionLiteralParsing(t *testing.T) {
	input := `func add(x, y) => { return x + y; }`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
			1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	function, ok := stmt.Expression.(*ast.FunctionLiteral)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.FunctionLiteral. got=%T",
			stmt.Expression)
	}

	if len(function.Parameters) != 2 {
		t.Fatalf("function literal parameters wrong. want 2, got=%d\n",
			len(function.Parameters))
	}

	if function.Parameters[0].Value != "x" {
		t.Fatalf("parameter[0] is not 'x'. got=%q", function.Parameters[0].Value)
	}

	if function.Parameters[1].Value != "y" {
		t.Fatalf("parameter[1] is not 'y'. got=%q", function.Parameters[1].Value)
	}

	if len(function.Body.Statements) != 1 {
		t.Fatalf("function.Body.Statements has not 1 statements. got=%d\n",
			len(function.Body.Statements))
	}
}

func TestCallExpressionParsing(t *testing.T) {
	input := "add(1, 2 * 3, 4 + 5);"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
			1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("stmt is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.CallExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.CallExpression. got=%T",
			stmt.Expression)
	}

	if len(exp.Arguments) != 3 {
		t.Fatalf("wrong length of arguments. got=%d", len(exp.Arguments))
	}
}

func TestWhileStatement(t *testing.T) {
	input := `while (x < 10) { x; }`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement. got=%d",
			len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.WhileStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.WhileStatement. got=%T",
			program.Statements[0])
	}

	if stmt.Condition == nil {
		t.Fatalf("stmt.Condition is nil")
	}

	if stmt.Body == nil {
		t.Fatalf("stmt.Body is nil")
	}
}

func TestElseIfStatement(t *testing.T) {
	input := `
if (x < 0) {
	outln("negative");
} else if (x == 0) {
	outln("zero");
} else if (x > 0) {
	outln("positive");
} else {
	outln("unknown");
}
`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement. got=%d",
			len(program.Statements))
	}

	ifStmt, ok := program.Statements[0].(*ast.IfStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.IfStatement. got=%T",
			program.Statements[0])
	}
	testElseIfChain(t, ifStmt)
}

func testElseIfChain(t *testing.T, ifStmt *ast.IfStatement) {
	// 最初のif
	if ifStmt.Consequence == nil {
		t.Fatalf("ifStmt.Consequence is nil")
	}

	// else if (x == 0)
	if ifStmt.Alternative == nil {
		t.Fatalf("ifStmt.Alternative is nil (first else if)")
	}

	if len(ifStmt.Alternative.Statements) != 1 {
		t.Fatalf("else branch does not contain 1 statement. got=%d",
			len(ifStmt.Alternative.Statements))
	}

	secondIf, ok := ifStmt.Alternative.Statements[0].(*ast.IfStatement)
	if !ok {
		t.Fatalf("first else branch is not IfStatement. got=%T",
			ifStmt.Alternative.Statements[0])
	}

	// else if (x > 0)
	if secondIf.Alternative == nil {
		t.Fatalf("secondIf.Alternative is nil (second else if)")
	}

	thirdIf, ok := secondIf.Alternative.Statements[0].(*ast.IfStatement)
	if !ok {
		t.Fatalf("second else branch is not IfStatement. got=%T",
			secondIf.Alternative.Statements[0])
	}

	// final else
	if thirdIf.Alternative == nil {
		t.Fatalf("thirdIf.Alternative is nil (final else)")
	}
}

func TestArrayLiteralParsing(t *testing.T) {
	input := "[1, 2 * 2, 3 + 3]"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	array, ok := stmt.Expression.(*ast.ArrayLiteral)
	if !ok {
		t.Fatalf("exp is not ast.ArrayLiteral. got=%T", stmt.Expression)
	}

	if len(array.Elements) != 3 {
		t.Fatalf("len(array.Elements) not 3. got=%d", len(array.Elements))
	}

	testNumberLiteral(t, array.Elements[0], "1")

	infix, ok := array.Elements[1].(*ast.InfixExpression)
	if !ok {
		t.Fatalf("array.Elements[1] is not ast.InfixExpression. got=%T", array.Elements[1])
	}
	if infix.Operator != "*" {
		t.Fatalf("infix.Operator is not '*'. got=%s", infix.Operator)
	}

	infix, ok = array.Elements[2].(*ast.InfixExpression)
	if !ok {
		t.Fatalf("array.Elements[2] is not ast.InfixExpression. got=%T", array.Elements[2])
	}
	if infix.Operator != "+" {
		t.Fatalf("infix.Operator is not '+'. got=%s", infix.Operator)
	}
}

func TestEmptyArrayLiteral(t *testing.T) {
	input := "[]"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	array, ok := stmt.Expression.(*ast.ArrayLiteral)
	if !ok {
		t.Fatalf("exp is not ast.ArrayLiteral. got=%T", stmt.Expression)
	}

	if len(array.Elements) != 0 {
		t.Fatalf("len(array.Elements) not 0. got=%d", len(array.Elements))
	}
}

func TestIndexExpressionParsing(t *testing.T) {
	input := "myArray[1 + 1]"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	indexExp, ok := stmt.Expression.(*ast.IndexExpression)
	if !ok {
		t.Fatalf("exp is not ast.IndexExpression. got=%T", stmt.Expression)
	}

	ident, ok := indexExp.Left.(*ast.Identifier)
	if !ok {
		t.Fatalf("indexExp.Left is not ast.Identifier. got=%T", indexExp.Left)
	}
	if ident.Value != "myArray" {
		t.Fatalf("ident.Value is not 'myArray'. got=%s", ident.Value)
	}

	infix, ok := indexExp.Index.(*ast.InfixExpression)
	if !ok {
		t.Fatalf("indexExp.Index is not ast.InfixExpression. got=%T", indexExp.Index)
	}
	if infix.Operator != "+" {
		t.Fatalf("infix.Operator is not '+'. got=%s", infix.Operator)
	}
}

func testNumberLiteral(t *testing.T, exp ast.Expression, value string) {
	num, ok := exp.(*ast.NumberLiteral)
	if !ok {
		t.Fatalf("exp is not ast.NumberLiteral. got=%T", exp)
	}
	if num.Value != value {
		t.Fatalf("num.Value is not %s. got=%s", value, num.Value)
	}
}

func TestParserErrorsWithPosition(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"mut x",
			"line 1, column 6: expected next token to be =, got EOF instead",
		},
		{
			"if x { }",
			"line 1, column 4: expected next token to be (, got IDENT instead",
		},
		{
			"@",
			"line 1, column 1: no prefix parse function for ILLEGAL found",
		},
		{
			"mut x = 10;\n@",
			"line 2, column 1: no prefix parse function for ILLEGAL found",
		},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		p.ParseProgram()

		errors := p.Errors()
		if len(errors) == 0 {
			t.Errorf("expected parser errors for input %q", tt.input)
			continue
		}

		if errors[0] != tt.expected {
			t.Errorf("wrong error message.\nexpected=%q\ngot=%q", tt.expected, errors[0])
		}
	}
}
