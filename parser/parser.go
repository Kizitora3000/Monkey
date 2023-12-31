package parser

import (
	"fmt"
	"monkey/ast"
	"monkey/lexer"
	"monkey/token"
	"strconv"
)

// iotaを使って，LOWESTが1，EQUALSが2，...と値を入れていく
// この値は演算時の優先順位を示している．1が最も速い
const (
	_ int = iota
	LOWEST
	EQUALS      // ==
	LESSGREATER // > or <
	SUM         // +
	PRODUCT     // *
	PREFIX      // -X or !X
	CALL        // myFunction(X)
)

// 優先順位テーブル
// EQ(=)とNOT_EQ(!=)は同じ優先順位(EQUALS=2)など
var precedences = map[token.TokenType]int{
	token.EQ:       EQUALS,
	token.NOT_EQ:   EQUALS,
	token.LT:       LESSGREATER,
	token.GT:       LESSGREATER,
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.SLASH:    PRODUCT,
	token.ASTERISK: PRODUCT,
	token.LPAREN:   CALL,
}

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

type Parser struct {
	l *lexer.Lexer

	errors []string

	// lexerでは文字列を読んでいたが、今回はトークンを取得する
	curToken  token.Token
	peekToken token.Token

	// -5 や !hoge などの前置演算子をパースするために，トークンと対応したパース関数をmapで紐づける
	prefixParseFns map[token.TokenType]prefixParseFn

	// 5+3 や 2 > 1 などの中値演算子をパースするために，トークンと対応したパース関数をmapで紐づける
	infixParseFns map[token.TokenType]infixParseFn
}

// Parserを作るためのコンストラクタ
func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}

	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	p.registerPrefix(token.IDENT, p.parseIdentifier)
	p.registerPrefix(token.INT, p.parseIntegerLiteral)
	p.registerPrefix(token.BANG, p.parsePrefixExpression)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)
	p.registerPrefix(token.TRUE, p.parseBoolean)
	p.registerPrefix(token.FALSE, p.parseBoolean)
	p.registerPrefix(token.LPAREN, p.parseGroupedExpression)
	p.registerPrefix(token.IF, p.parseIfExpression)
	p.registerPrefix(token.FUNCTION, p.parseFunctionLiteral)

	p.infixParseFns = make(map[token.TokenType]infixParseFn)
	p.registerInfix(token.PLUS, p.parseInfixExpression)
	p.registerInfix(token.MINUS, p.parseInfixExpression)
	p.registerInfix(token.SLASH, p.parseInfixExpression)
	p.registerInfix(token.ASTERISK, p.parseInfixExpression)
	p.registerInfix(token.EQ, p.parseInfixExpression)
	p.registerInfix(token.NOT_EQ, p.parseInfixExpression)
	p.registerInfix(token.LT, p.parseInfixExpression)
	p.registerInfix(token.GT, p.parseInfixExpression)

	// 呼び出し式 add(2, 3) における中値演算子"("を登録する
	// これは add という識別子(関数を束縛している) と 2,3 という引数リストの2つの式を持つ必要があるので，中値演算子となる
	p.registerInfix(token.LPAREN, p.parseCallExpression)

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

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
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
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
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

	p.nextToken()

	stmt.Value = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.curToken}

	p.nextToken()

	stmt.ReturnValue = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.curToken}

	stmt.Expression = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) noPrefixParseFnError(t token.TokenType) {
	msg := fmt.Sprintf("no prefix parse function for %s found", t)
	p.errors = append(p.errors, msg)
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		p.noPrefixParseFnError(p.curToken.Type)
		return nil
	}
	leftExp := prefix()

	// 優先順位を考慮する処理
	// ref. 専攻科→コンパイラ(goodnote)
	for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}

		p.nextToken()

		leftExp = infix(leftExp)
	}

	return leftExp
}

// p.peekToken(一つ先)のトークンタイプに応じた優先順位を返す
func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}

	return LOWEST
}

// p.curToken(今見ているところ)のトークンタイプに応じた優先順位を返す
func (p *Parser) curPrecedence() int {
	if p, ok := precedences[p.curToken.Type]; ok {
		return p
	}

	return LOWEST
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}

	p.nextToken()

	expression.Right = p.parseExpression(PREFIX)

	return expression
}

// 引数: "left ast.Expression" → *ast.InfixExpressionノードを構築する際に"left ast.Expression"(中値演算子の左側の式(expression))を使う
func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     left,
	}

	// 現在の演算子の優先順位を得る
	precedence := p.curPrecedence()
	p.nextToken()

	// 次のトークンの優先順位を得る
	expression.Right = p.parseExpression(precedence)

	return expression
}

func (p *Parser) parseBoolean() ast.Expression {
	return &ast.Boolean{Token: p.curToken, Value: p.curTokenIs(token.TRUE)}
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken()

	exp := p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return exp
}

func (p *Parser) parseIfExpression() ast.Expression {
	expression := &ast.IfExpression{Token: p.curToken}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	// p.curToekenのトークンがtoken.LBRACEになるまでトークンを読み進める
	p.nextToken()
	expression.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	// "{"までたどり着いたら，ブロック内の内容をパースする
	expression.Consequence = p.parseBlockStatement()

	// elseが存在する場合のみ，else内のブロックも評価する
	// "peekTokenIs"を使用しているように，elseを省略していても動作するようになっている
	if p.peekTokenIs(token.ELSE) {
		// peekToeknがelseなので，まずは一個進める
		p.nextToken()

		// curTokenがelseのときに，次のトークンが"("じゃなかったらおかしいのでチェック
		if !p.expectPeek(token.LBRACE) {
			return nil
		}

		expression.Alternative = p.parseBlockStatement()
	}

	return expression
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{Token: p.curToken}
	block.Statements = []ast.Statement{}

	p.nextToken()

	// "}" に到達するか EOF (終端記号)に到達するまで読み進める
	// "}" → 正常終了, "EOF" → 正しい構文ではない
	for !p.curTokenIs(token.RBRACE) && !p.curTokenIs(token.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.nextToken()
	}

	return block
}

func (p *Parser) parseFunctionLiteral() ast.Expression {
	lit := &ast.FunctionLiteral{Token: p.curToken}

	// fnトークンの後が"("かどうかチェック
	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	lit.Parameters = p.parseFunctionParameters()

	// 引数"(a, b, c)"を読み込んだ後が"{"かどうかチェック
	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	lit.Body = p.parseBlockStatement()

	return lit
}

func (p *Parser) parseFunctionParameters() []*ast.Identifier {
	identifiers := []*ast.Identifier{}

	// 引数が空の場合，curToken="("において，次のトークンはpeekToken=")"となるのでそこで終わらせる
	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		return identifiers
	}

	p.nextToken()

	// 第一引数を読み取る
	ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	identifiers = append(identifiers, ident)

	// 第二引数以降はカンマごとに区切られているので，カンマごとに読み取っていく
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

func (p *Parser) parseCallExpression(function ast.Expression) ast.Expression {
	// CallExpression ノードの作成
	exp := &ast.CallExpression{Token: p.curToken, Function: function}

	// CallExpression ノードの Arguments に 関数呼び出しの引数部分を格納
	exp.Arguments = p.parseCallArguments()
	return exp
}

// ここの処理は parseFunctionParameters関数とほぼ同じ
// 違いは2つある．
// 1つ目：引数にあたる部分に識別子や関数リテラルなどいろんなものが入るので，parseExpressionを使っていパン化しているところ
// 2つ目：返り値が []*ast.Identifier ではなく []ast.Expression となる
// 関数リテラルの定義時は識別子だけだが，関数呼び出しでは識別子だけでなく，引数に関数リテラルを入れることもできる
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

func (p *Parser) parseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{Token: p.curToken}

	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}

	lit.Value = value

	return lit
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

func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}
