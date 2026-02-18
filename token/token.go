package token

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
	Line    int
	Column  int
}

const (
	// 特殊トークン
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	// 識別子とリテラル
	IDENT  = "IDENT"  // 変数名、関数名
	NUMBER = "NUMBER" // 数値（整数・小数）
	STRING = "STRING" // 文字列

	// 演算子
	ASSIGN   = "="
	PLUS     = "+"
	MINUS    = "-"
	ASTERISK = "*"
	SLASH    = "/"
	PERCENT  = "%"
	BANG     = "!"

	// 後置演算子
	PLUS_PLUS   = "++"
	MINUS_MINUS = "--"

	// 複合代入演算子
	PLUS_ASSIGN     = "+="
	MINUS_ASSIGN    = "-="
	ASTERISK_ASSIGN = "*="
	SLASH_ASSIGN    = "/="
	PERCENT_ASSIGN  = "%="

	// 比較演算子
	EQ     = "=="
	NOT_EQ = "!="
	LT     = "<"
	GT     = ">"
	LT_EQ  = "<="
	GT_EQ  = ">="

	// 論理演算子
	AND = "&&"
	OR  = "||"

	// デリミタ
	COMMA     = ","
	SEMICOLON = ";"
	COLON     = ":"
	ARROW     = "=>"

	LPAREN   = "("
	RPAREN   = ")"
	LBRACE   = "{"
	RBRACE   = "}"
	LBRACKET = "["
	RBRACKET = "]"

	// キーワード
	MUT      = "MUT"
	CONST    = "CONST"
	FUNC     = "FUNC"
	RETURN   = "RETURN"
	IF       = "IF"
	ELSE     = "ELSE"
	SWITCH   = "SWITCH"
	CASE     = "CASE"
	DEFAULT  = "DEFAULT"
	WHILE    = "WHILE"
	FOR      = "FOR"
	BREAK    = "BREAK"
	CONTINUE = "CONTINUE"
	TRUE     = "TRUE"
	FALSE    = "FALSE"
	NULL     = "NULL"
	TRY      = "TRY"
	CATCH    = "CATCH"
	THROW    = "THROW"
)

var keywords = map[string]TokenType{
	"mut":      MUT,
	"const":    CONST,
	"func":     FUNC,
	"return":   RETURN,
	"if":       IF,
	"else":     ELSE,
	"switch":   SWITCH,
	"case":     CASE,
	"default":  DEFAULT,
	"while":    WHILE,
	"for":      FOR,
	"break":    BREAK,
	"continue": CONTINUE,
	"true":     TRUE,
	"false":    FALSE,
	"null":     NULL,
	"try":      TRY,
	"catch":    CATCH,
	"throw":    THROW,
}

// LookupIdent は識別子がキーワードかどうかを判定する
func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
