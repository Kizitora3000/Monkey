package token

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

const (
	ILLEGAL = "ILLEGAL" // 未知のト－クン・文字であることを意味する
	EOF     = "EOF"     // ファイル終端

	// 識別子(変数名), リテラル
	IDENT = "IDENT"
	INT   = "INT"

	// 演算子
	ASSIGN = "="
	PLUS   = "+"

	// デリミタ(区切り文字)
	COMMA     = ","
	SEMICOLON = ";"

	LPAREN = "("
	RPAREN = ")"
	LBRACE = "{"
	RBRACE = "}"

	// キーワード(予約語)
	FUNCTION = "FUNCTION"
	LET      = "LET"
)

// 予約語と識別子（変数名, 関数名, etc.）
var keywords = map[string]TokenType{
	"fn":  FUNCTION,
	"let": LET,
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {

	}
}
