package lexer

import (
	"fmt"
	"monkey/token"
)

type Lexer struct {
	input       string
	position    int  // 現在読んでいる文字(ch)の位置
	readPositon int  // positionの次の位置
	ch          byte // 現在読んでいる文字
}

// Lexer のコンストラクタ
func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

// Lexer のメソッド関数
// 最初が小文字 → Lexerパッケージからのみ利用できる, 最初が大文字 → 他のパッケージでも使用できる
func (l *Lexer) readChar() {
	// 次の一文字が終端に到達したかどうかのチェック
	if l.readPositon >= len(l.input) {
		// ch = 0 はEOFを意味する
		l.ch = 0
	} else {
		// 現在はASCII文字だけを扱っているため、次の1バイトだけを読み込めばよい
		// UnicodeやUTF-8を扱う場合、次の一文字が複数のバイトで構成される可能性があるため別途処理が必要となる
		l.ch = l.input[l.readPositon]
	}

	l.position = l.readPositon
	l.readPositon += 1
}

func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token
	fmt.Printf("l.ch: ")
	fmt.Println(l.ch)
	switch l.ch {
	case '=':
		tok = newToken(token.ASSIGN, l.ch)
	case ';':
		tok = newToken(token.SEMICOLON, l.ch)
	case '(':
		tok = newToken(token.LPAREN, l.ch)
	case ')':
		tok = newToken(token.RPAREN, l.ch)
	case ',':
		tok = newToken(token.COMMA, l.ch)
	case '+':
		tok = newToken(token.PLUS, l.ch)
	case '{':
		tok = newToken(token.LBRACE, l.ch)
	case '}':
		tok = newToken(token.RBRACE, l.ch)
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		// letter: 英字
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			return tok
		} else {
			tok = newToken(token.ILLEGAL, l.ch)
		}
	}

	l.readChar()
	return tok
}

// 例：「let abc ...」があったとき、position = 'a'の位置、l.position = 'c'の位置となり、その範囲のabcを取得する
func (l *Lexer) readIdentifier() string {
	// 最初の基準となる位置を把握しておく[]
	position := l.position

	// 識別子を非英字になるまで読み進める
	for isLetter(l.ch) {
		l.readChar()
	}

	return l.input[position:l.position]
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}
