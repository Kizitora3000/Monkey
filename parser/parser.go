package parser

import (
	"fmt"
	"monkey/ast"
	"monkey/lexer"
	"monkey/token"
)

type Parser struct {
	l *lexer.Lexer

	errors []string

	// lexerでは文字列を読んでいたが、今回はトークンを取得する
	curToken  token.Token
	peekToken token.Token
}

// Parserを作るためのコンストラクタ
func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}

	// 現在調べているトークンだけだと十分な情報が得られない場合があるので、次のトークンも調べるようにする
	/* 十分な情報を得られない例
	5; なのか 5 + 5; なのかを判別するとき．;があるから処理を終えるのか，+だから演算子に関連したパーサを呼び出すのか
	*/

	/* let a = 5 の例
	1回目のnextToken: p.curToken = nil, p.peekToken = p.peekToken() = let
	2回目のnextToken: p.curToken = let, p.peekToken = p.peekToken() = a
	*/
	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead",
		t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

// peekToekenに今見ているトークンが入っている
// nextToken関数が呼ばれたらcurTokenに今見ていたトークンを格納して、peekTokenは次のトークンを見るようにする
func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	// トークンがEOF(トークン列の最後)になるまで入力されたトークン列を読む
	for p.curToken.Type != token.EOF {
		stmt := p.parseStatement()

		// parseStatementで読み込んだステートメントがnil以外（事前に定義したletやreturnなど）であればStatementsに追加する
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}

	return program
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.LET:
		return p.parseLetStatement()
	default:
		return nil
	}
}

// 初めがletで，次が識別子，その次が=であることをチェック素すr
// 初めがletなのは既にparseStatementで確定させている
func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: p.curToken}

	// 識別子のチェック
	if !p.expectPeek(token.IDENT) {
		return nil
	}

	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	// 等号のチェック
	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	// TODO: セミコロンに遭遇するまで式を読み飛ばす（=より右側の式はまだ記述していないので，一旦飛ばす）
	if !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

// 現在のトークンが引数のトークンと一致しているかどうかの等号演算
func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}

// 次のトークンが引数のトークンと一致しているかどうかの等号演算
func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

// 次のトークンが想定しているトークンと一致していれば次のトークンを読み進める
// 例えば let x = 5; となっているところが let x 5;となっていれば return false を返す
func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		p.peekError(t)
		return false
	}
}
