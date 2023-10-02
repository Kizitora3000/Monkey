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
		tok = newToken(token.LBACE, l.ch)
	case '}':
		tok = newToken(token.RBACE, l.ch)
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	}

	l.readChar()
	return tok
}
