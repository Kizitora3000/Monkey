package token

type TokenType string

type Token struct {
	Type    TokenType // トークンが数(INT)なのか？変数(IDENT)なのか？キーワード(FUNCTION, LET, etc.)なのか、という種類を示す
	Literal string    // トークンの内容が入る。数ならその値、変数なら変数名が入る。キーワードはTokenTypeと同じ
}

const (
	ILLEGAL = "ILLEGAL" // 未知のト－クン・文字であることを意味する
	EOF     = "EOF"     // ファイル終端

	// 識別子(変数名), リテラル
	IDENT = "IDENT"
	INT   = "INT"

	// 演算子
	ASSIGN   = "="
	PLUS     = "+"
	MINUS    = "-"
	BANG     = "!"
	ASTERISK = "*"
	SLASH    = "/"

	LT = "<"
	GT = ">"

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
	TRUE     = "TRUE"
	FALSE    = "FALSE"
	IF       = "IF"
	ELSE     = "ELSE"
	RETURN   = "RETURN"
	EQ       = "=="
	NOT_EQ   = "!="
)

// 予約語の定義
var keywords = map[string]TokenType{
	"fn":     FUNCTION,
	"let":    LET,
	"true":   TRUE,
	"false":  FALSE,
	"if":     IF,
	"else":   ELSE,
	"return": RETURN,
}

// 予約語と識別子（変数名, 関数名, etc.）の識別を行う
// 予約語ならそのトークンを、そうでなければ識別子を意味する"IDENT"を返す
func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
