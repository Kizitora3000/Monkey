package parser

import (
	"monkey/ast"
	"monkey/lexer"
	"monkey/token"
)

type Parser struct {
	l *lexer.Lexer

	// lexerでは文字列を読んでいたが、今回はトークンを取得する
	curlToekn token.Token
	peekToken token.Token
}

// Parserを作るためのコンストラクタ
func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l}

	// 現在調べているトークンだけだと十分な情報が得られない場合があるので、次のトークンも調べるようにする
	p.nextToken()
	p.nextToken()

	return p
}

// peekToekenに今見ているトークンが入っている
// nextToken関数が呼ばれたらcurlTokenに今見ていたトークンを格納して、peekTokenは次のトークンを見るようにする
func (p *Parser) nextToken() {
	p.curlToekn = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) ParseProgram() *ast.Program {
	return nil
}
