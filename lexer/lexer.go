package lexer

import (
	"sugu/token"
)

type Lexer struct {
	input        string
	position     int    // 現在の位置（現在の文字を指す）
	readPosition int    // 次の位置（現在の文字の次を指す）
	ch           byte   // 現在検査中の文字
	stringError  string // 文字列パース中のエラー
	line         int    // 現在の行番号（1から始まる）
	column       int    // 現在の列番号（1から始まる）
}

func New(input string) *Lexer {
	l := &Lexer{input: input, line: 1, column: 0}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition++

	// 行番号・列番号を更新
	if l.ch == '\n' {
		l.line++
		l.column = 0
	} else {
		l.column++
	}
}

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}
	return l.input[l.readPosition]
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhitespace()

	// トークンの開始位置を記録
	startLine := l.line
	startColumn := l.column

	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			l.readChar()
			tok = l.newTokenWithLiteral(token.EQ, "==", startLine, startColumn)
		} else if l.peekChar() == '>' {
			l.readChar()
			tok = l.newTokenWithLiteral(token.ARROW, "=>", startLine, startColumn)
		} else {
			tok = l.newToken(token.ASSIGN, l.ch)
		}
	case '+':
		if l.peekChar() == '+' {
			l.readChar()
			tok = l.newTokenWithLiteral(token.PLUS_PLUS, "++", startLine, startColumn)
		} else if l.peekChar() == '=' {
			l.readChar()
			tok = l.newTokenWithLiteral(token.PLUS_ASSIGN, "+=", startLine, startColumn)
		} else {
			tok = l.newToken(token.PLUS, l.ch)
		}
	case '-':
		if l.peekChar() == '-' {
			l.readChar()
			tok = l.newTokenWithLiteral(token.MINUS_MINUS, "--", startLine, startColumn)
		} else if l.peekChar() == '=' {
			l.readChar()
			tok = l.newTokenWithLiteral(token.MINUS_ASSIGN, "-=", startLine, startColumn)
		} else {
			tok = l.newToken(token.MINUS, l.ch)
		}
	case '*':
		if l.peekChar() == '=' {
			l.readChar()
			tok = l.newTokenWithLiteral(token.ASTERISK_ASSIGN, "*=", startLine, startColumn)
		} else {
			tok = l.newToken(token.ASTERISK, l.ch)
		}
	case '/':
		if l.peekChar() == '/' {
			l.skipComment()
			return l.NextToken()
		} else if l.peekChar() == '=' {
			l.readChar()
			tok = l.newTokenWithLiteral(token.SLASH_ASSIGN, "/=", startLine, startColumn)
		} else {
			tok = l.newToken(token.SLASH, l.ch)
		}
	case '%':
		if l.peekChar() == '=' {
			l.readChar()
			tok = l.newTokenWithLiteral(token.PERCENT_ASSIGN, "%=", startLine, startColumn)
		} else {
			tok = l.newToken(token.PERCENT, l.ch)
		}
	case '!':
		if l.peekChar() == '=' {
			l.readChar()
			tok = l.newTokenWithLiteral(token.NOT_EQ, "!=", startLine, startColumn)
		} else {
			tok = l.newToken(token.BANG, l.ch)
		}
	case '<':
		if l.peekChar() == '=' {
			l.readChar()
			tok = l.newTokenWithLiteral(token.LT_EQ, "<=", startLine, startColumn)
		} else {
			tok = l.newToken(token.LT, l.ch)
		}
	case '>':
		if l.peekChar() == '=' {
			l.readChar()
			tok = l.newTokenWithLiteral(token.GT_EQ, ">=", startLine, startColumn)
		} else {
			tok = l.newToken(token.GT, l.ch)
		}
	case '&':
		if l.peekChar() == '&' {
			l.readChar()
			tok = l.newTokenWithLiteral(token.AND, "&&", startLine, startColumn)
		} else {
			tok = l.newToken(token.ILLEGAL, l.ch)
		}
	case '|':
		if l.peekChar() == '|' {
			l.readChar()
			tok = l.newTokenWithLiteral(token.OR, "||", startLine, startColumn)
		} else {
			tok = l.newToken(token.ILLEGAL, l.ch)
		}
	case ',':
		tok = l.newToken(token.COMMA, l.ch)
	case ';':
		tok = l.newToken(token.SEMICOLON, l.ch)
	case ':':
		tok = l.newToken(token.COLON, l.ch)
	case '(':
		tok = l.newToken(token.LPAREN, l.ch)
	case ')':
		tok = l.newToken(token.RPAREN, l.ch)
	case '{':
		tok = l.newToken(token.LBRACE, l.ch)
	case '}':
		tok = l.newToken(token.RBRACE, l.ch)
	case '[':
		tok = l.newToken(token.LBRACKET, l.ch)
	case ']':
		tok = l.newToken(token.RBRACKET, l.ch)
	case '"':
		l.stringError = "" // エラーをリセット
		literal := l.readString()
		if l.stringError != "" {
			tok = l.newTokenWithLiteral(token.ILLEGAL, l.stringError, startLine, startColumn)
		} else {
			tok = l.newTokenWithLiteral(token.STRING, literal, startLine, startColumn)
		}
	case 0:
		tok = l.newTokenWithLiteral(token.EOF, "", startLine, startColumn)
	default:
		if isLetter(l.ch) {
			literal := l.readIdentifier()
			tok = l.newTokenWithLiteral(token.LookupIdent(literal), literal, startLine, startColumn)
			return tok
		} else if isDigit(l.ch) {
			literal := l.readNumber()
			tok = l.newTokenWithLiteral(token.NUMBER, literal, startLine, startColumn)
			return tok
		} else {
			tok = l.newToken(token.ILLEGAL, l.ch)
		}
	}

	l.readChar()
	return tok
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

func (l *Lexer) skipComment() {
	// "//" を読み飛ばす
	l.readChar() // 最初の '/'
	l.readChar() // 2番目の '/'

	// 複数行コメント: //-- ... --//
	if l.ch == '-' && l.peekChar() == '-' {
		l.readChar() // '-'
		l.readChar() // '-'
		for {
			if l.ch == 0 {
				return
			}
			if l.ch == '-' && l.peekChar() == '-' {
				l.readChar() // '-'
				l.readChar() // '-'
				if l.ch == '/' && l.peekChar() == '/' {
					l.readChar() // '/'
					l.readChar() // '/'
					return
				}
			}
			l.readChar()
		}
	}

	// 単一行コメント: // ...
	for l.ch != '\n' && l.ch != 0 {
		l.readChar()
	}
}

func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) || isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	// 小数点の処理
	if l.ch == '.' && isDigit(l.peekChar()) {
		l.readChar() // '.'
		for isDigit(l.ch) {
			l.readChar()
		}
	}
	return l.input[position:l.position]
}

func (l *Lexer) readString() string {
	var result []byte
	l.readChar() // 開始の '"' をスキップ

	for l.ch != '"' && l.ch != 0 {
		if l.ch == '\\' {
			l.readChar()
			switch l.ch {
			case 'n':
				result = append(result, '\n')
			case 't':
				result = append(result, '\t')
			case 'r':
				result = append(result, '\r')
			case '"':
				result = append(result, '"')
			case '\\':
				result = append(result, '\\')
			case 0:
				// バックスラッシュ直後にEOF
				l.stringError = "unexpected end of string after \\"
				return string(result)
			default:
				// 不明なエスケープシーケンス
				l.stringError = "unknown escape sequence: \\" + string(l.ch)
				return string(result)
			}
		} else {
			result = append(result, l.ch)
		}
		l.readChar()
	}

	return string(result)
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func (l *Lexer) newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch), Line: l.line, Column: l.column}
}

func (l *Lexer) newTokenWithLiteral(tokenType token.TokenType, literal string, line, column int) token.Token {
	return token.Token{Type: tokenType, Literal: literal, Line: line, Column: column}
}
