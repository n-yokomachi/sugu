package parser

import (
	"fmt"
	"sugu/ast"
	"sugu/lexer"
	"sugu/token"
)

// Parser はトークン列からASTを構築する
type Parser struct {
	l      *lexer.Lexer
	errors []string

	curToken  token.Token
	peekToken token.Token
}

// New は新しいParserを作成する
func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}
	// 2つトークンを読み込んで curToken と peekToken をセット
	p.nextToken()
	p.nextToken()
	return p
}

// Errors はパースエラーを返す
func (p *Parser) Errors() []string {
	return p.errors
}

// nextToken はトークンを1つ進める
func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

// curTokenIs は現在のトークンが指定された型かチェック
func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}

// peekTokenIs は次のトークンが指定された型かチェック
func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

// expectPeek は次のトークンが期待された型かチェックし、そうなら進める
func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	}
	p.peekError(t)
	return false
}

// peekError はエラーメッセージを記録
func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead",
		t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

// ParseProgram はプログラム全体をパース
func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for !p.curTokenIs(token.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}

	return program
}

// parseStatement は文をパース
func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.MUT, token.CONST:
		return p.parseVariableStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	case token.IF:
		return p.parseIfStatement()
	case token.WHILE:
		return p.parseWhileStatement()
	case token.FOR:
		return p.parseForStatement()
	case token.SWITCH:
		return p.parseSwitchStatement()
	case token.BREAK:
		return p.parseBreakStatement()
	case token.CONTINUE:
		return p.parseContinueStatement()
	default:
		return p.parseExpressionStatement()
	}
}

// parseVariableStatement は変数宣言をパース
func (p *Parser) parseVariableStatement() *ast.VariableStatement {
	stmt := &ast.VariableStatement{Token: p.curToken}

	if !p.expectPeek(token.IDENT) {
		return nil
	}

	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	p.nextToken()

	stmt.Value = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

// parseReturnStatement はreturn文をパース
func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.curToken}

	p.nextToken()

	stmt.ReturnValue = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

// parseExpressionStatement は式文をパース
func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.curToken}

	stmt.Expression = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

// 演算子の優先順位
const (
	_ int = iota
	LOWEST
	ASSIGN      // =
	LOGICAL_OR  // ||
	LOGICAL_AND // &&
	EQUALS      // ==
	LESSGREATER // > or <
	SUM         // +
	PRODUCT     // *
	PREFIX      // -X or !X
	CALL        // myFunction(X)
	INDEX       // array[index]
)

// 優先順位テーブル
var precedences = map[token.TokenType]int{
	token.ASSIGN:   ASSIGN,
	token.OR:       LOGICAL_OR,
	token.AND:      LOGICAL_AND,
	token.EQ:       EQUALS,
	token.NOT_EQ:   EQUALS,
	token.LT:       LESSGREATER,
	token.GT:       LESSGREATER,
	token.LT_EQ:    LESSGREATER,
	token.GT_EQ:    LESSGREATER,
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.ASTERISK: PRODUCT,
	token.SLASH:    PRODUCT,
	token.PERCENT:  PRODUCT,
	token.LPAREN:   CALL,
	token.LBRACKET: INDEX,
}

// peekPrecedence は次のトークンの優先順位を返す
func (p *Parser) peekPrecedence() int {
	if prec, ok := precedences[p.peekToken.Type]; ok {
		return prec
	}
	return LOWEST
}

// curPrecedence は現在のトークンの優先順位を返す
func (p *Parser) curPrecedence() int {
	if prec, ok := precedences[p.curToken.Type]; ok {
		return prec
	}
	return LOWEST
}

// parseExpression は式をパース（Pratt parsing）
func (p *Parser) parseExpression(precedence int) ast.Expression {
	// 前置式のパース
	var leftExp ast.Expression

	switch p.curToken.Type {
	case token.IDENT:
		leftExp = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	case token.NUMBER:
		leftExp = &ast.NumberLiteral{Token: p.curToken, Value: p.curToken.Literal}
	case token.STRING:
		leftExp = &ast.StringLiteral{Token: p.curToken, Value: p.curToken.Literal}
	case token.TRUE, token.FALSE:
		leftExp = &ast.BooleanLiteral{Token: p.curToken, Value: p.curTokenIs(token.TRUE)}
	case token.NULL:
		leftExp = &ast.NullLiteral{Token: p.curToken}
	case token.BANG, token.MINUS:
		leftExp = p.parsePrefixExpression()
	case token.LPAREN:
		leftExp = p.parseGroupedExpression()
	case token.FUNC:
		leftExp = p.parseFunctionLiteral()
	case token.LBRACKET:
		leftExp = p.parseArrayLiteral()
	case token.LBRACE:
		leftExp = p.parseMapLiteral()
	default:
		msg := fmt.Sprintf("no prefix parse function for %s found", p.curToken.Type)
		p.errors = append(p.errors, msg)
		return nil
	}

	// 中置式のパース
	for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() {
		switch p.peekToken.Type {
		case token.ASSIGN:
			// 代入式は左辺が識別子の場合のみ有効
			ident, ok := leftExp.(*ast.Identifier)
			if !ok {
				return leftExp
			}
			p.nextToken()
			leftExp = p.parseAssignExpression(ident)
		case token.PLUS, token.MINUS, token.ASTERISK, token.SLASH, token.PERCENT,
			token.EQ, token.NOT_EQ, token.LT, token.GT, token.LT_EQ, token.GT_EQ,
			token.AND, token.OR:
			p.nextToken()
			leftExp = p.parseInfixExpression(leftExp)
		case token.LPAREN:
			p.nextToken()
			leftExp = p.parseCallExpression(leftExp)
		case token.LBRACKET:
			p.nextToken()
			leftExp = p.parseIndexExpression(leftExp)
		default:
			return leftExp
		}
	}

	return leftExp
}

// parsePrefixExpression は前置演算子式をパース
func (p *Parser) parsePrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}

	p.nextToken()
	expression.Right = p.parseExpression(PREFIX)

	return expression
}

// parseAssignExpression は代入式をパース
func (p *Parser) parseAssignExpression(name *ast.Identifier) ast.Expression {
	expression := &ast.AssignExpression{
		Token: p.curToken,
		Name:  name,
	}

	p.nextToken()
	expression.Value = p.parseExpression(ASSIGN)

	return expression
}

// parseInfixExpression は中置演算子式をパース
func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     left,
	}

	precedence := p.curPrecedence()
	p.nextToken()
	expression.Right = p.parseExpression(precedence)

	return expression
}

// parseGroupedExpression はグループ化された式をパース
func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken()

	exp := p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return exp
}

// parseCallExpression は関数呼び出し式をパース
func (p *Parser) parseCallExpression(function ast.Expression) ast.Expression {
	exp := &ast.CallExpression{Token: p.curToken, Function: function}
	exp.Arguments = p.parseCallArguments()
	return exp
}

// parseCallArguments は関数の引数リストをパース
func (p *Parser) parseCallArguments() []ast.Expression {
	args := []ast.Expression{}

	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		return args
	}

	p.nextToken()
	args = append(args, p.parseExpression(LOWEST))

	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		args = append(args, p.parseExpression(LOWEST))
	}

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return args
}

// parseIfStatement はif文をパース
func (p *Parser) parseIfStatement() *ast.IfStatement {
	stmt := &ast.IfStatement{Token: p.curToken}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	p.nextToken()
	stmt.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	stmt.Consequence = p.parseBlockStatement()

	if p.peekTokenIs(token.ELSE) {
		p.nextToken()

		// else if をサポート
		if p.peekTokenIs(token.IF) {
			p.nextToken()
			ifStmt := p.parseIfStatement()
			if ifStmt == nil {
				return nil
			}
			stmt.Alternative = wrapIfInBlock(ifStmt)
		} else {
			if !p.expectPeek(token.LBRACE) {
				return nil
			}
			stmt.Alternative = p.parseBlockStatement()
		}
	}

	return stmt
}

// wrapIfInBlock はIfStatementをBlockStatementでラップする
func wrapIfInBlock(ifStmt *ast.IfStatement) *ast.BlockStatement {
	return &ast.BlockStatement{
		Token: ifStmt.Token,
		Statements: []ast.Statement{ifStmt},
	}
}

// parseBlockStatement はブロック文をパース
func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{Token: p.curToken}
	block.Statements = []ast.Statement{}

	p.nextToken()

	for !p.curTokenIs(token.RBRACE) && !p.curTokenIs(token.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.nextToken()
	}

	return block
}

// parseWhileStatement はwhile文をパース
func (p *Parser) parseWhileStatement() *ast.WhileStatement {
	stmt := &ast.WhileStatement{Token: p.curToken}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	p.nextToken()
	stmt.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	stmt.Body = p.parseBlockStatement()

	return stmt
}

// parseForStatement はfor文をパース
func (p *Parser) parseForStatement() *ast.ForStatement {
	stmt := &ast.ForStatement{Token: p.curToken}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	// 初期化式（オプショナル）
	if !p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
		stmt.Init = p.parseStatement()
	} else {
		p.nextToken()
	}

	// 条件式（オプショナル）
	if !p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
		stmt.Condition = p.parseExpression(LOWEST)
	}

	if !p.expectPeek(token.SEMICOLON) {
		return nil
	}

	// 更新式（オプショナル）
	if !p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		stmt.Update = p.parseExpression(LOWEST)
	}

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	stmt.Body = p.parseBlockStatement()

	return stmt
}

// parseSwitchStatement はswitch文をパース
func (p *Parser) parseSwitchStatement() *ast.SwitchStatement {
	stmt := &ast.SwitchStatement{Token: p.curToken}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	p.nextToken()
	stmt.Value = p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	p.nextToken()

	stmt.Cases = []*ast.CaseClause{}

	for !p.curTokenIs(token.RBRACE) && !p.curTokenIs(token.EOF) {
		if p.curTokenIs(token.CASE) {
			caseClause := &ast.CaseClause{Token: p.curToken}

			p.nextToken()
			caseClause.Value = p.parseExpression(LOWEST)

			if !p.expectPeek(token.COLON) {
				return nil
			}

			if !p.expectPeek(token.LBRACE) {
				return nil
			}

			caseClause.Body = p.parseBlockStatement()

			stmt.Cases = append(stmt.Cases, caseClause)

			p.nextToken()
		} else if p.curTokenIs(token.DEFAULT) {
			if !p.expectPeek(token.COLON) {
				return nil
			}

			if !p.expectPeek(token.LBRACE) {
				return nil
			}

			stmt.Default = p.parseBlockStatement()

			p.nextToken()
		} else {
			p.nextToken()
		}
	}

	return stmt
}

// parseBreakStatement はbreak文をパース
func (p *Parser) parseBreakStatement() *ast.BreakStatement {
	stmt := &ast.BreakStatement{Token: p.curToken}

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

// parseContinueStatement はcontinue文をパース
func (p *Parser) parseContinueStatement() *ast.ContinueStatement {
	stmt := &ast.ContinueStatement{Token: p.curToken}

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

// parseFunctionLiteral は関数リテラルをパース
func (p *Parser) parseFunctionLiteral() ast.Expression {
	lit := &ast.FunctionLiteral{Token: p.curToken}

	// 関数名（オプショナル）
	if p.peekTokenIs(token.IDENT) {
		p.nextToken()
		lit.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	lit.Parameters = p.parseFunctionParameters()

	if !p.expectPeek(token.ARROW) {
		return nil
	}

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	lit.Body = p.parseBlockStatement()

	return lit
}

// parseFunctionParameters は関数のパラメータリストをパース
func (p *Parser) parseFunctionParameters() []*ast.Identifier {
	identifiers := []*ast.Identifier{}

	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		return identifiers
	}

	p.nextToken()

	ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	identifiers = append(identifiers, ident)

	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
		identifiers = append(identifiers, ident)
	}

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return identifiers
}

// parseArrayLiteral は配列リテラルをパース
func (p *Parser) parseArrayLiteral() ast.Expression {
	array := &ast.ArrayLiteral{Token: p.curToken}
	array.Elements = p.parseExpressionList(token.RBRACKET)
	return array
}

// parseExpressionList は式のリストをパース（配列要素や関数引数に使用）
func (p *Parser) parseExpressionList(end token.TokenType) []ast.Expression {
	list := []ast.Expression{}

	if p.peekTokenIs(end) {
		p.nextToken()
		return list
	}

	p.nextToken()
	list = append(list, p.parseExpression(LOWEST))

	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		list = append(list, p.parseExpression(LOWEST))
	}

	if !p.expectPeek(end) {
		return nil
	}

	return list
}

// parseIndexExpression はインデックスアクセス式をパース
func (p *Parser) parseIndexExpression(left ast.Expression) ast.Expression {
	exp := &ast.IndexExpression{Token: p.curToken, Left: left}

	p.nextToken()
	exp.Index = p.parseExpression(LOWEST)

	if !p.expectPeek(token.RBRACKET) {
		return nil
	}

	return exp
}

// parseMapLiteral はマップリテラルをパース
func (p *Parser) parseMapLiteral() ast.Expression {
	mapLit := &ast.MapLiteral{Token: p.curToken}
	mapLit.Pairs = make(map[ast.Expression]ast.Expression)

	for !p.peekTokenIs(token.RBRACE) {
		p.nextToken()
		key := p.parseExpression(LOWEST)

		if !p.expectPeek(token.COLON) {
			return nil
		}

		p.nextToken()
		value := p.parseExpression(LOWEST)

		mapLit.Pairs[key] = value

		if !p.peekTokenIs(token.RBRACE) && !p.expectPeek(token.COMMA) {
			return nil
		}
	}

	if !p.expectPeek(token.RBRACE) {
		return nil
	}

	return mapLit
}
