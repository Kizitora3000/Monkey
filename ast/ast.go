package ast

import (
	"monkey/token"
)

type Node interface {
	TokenLiteral() string
}

type Statement interface {
	Node
	statementNode() // この関数はダミーメソッドで、Statementの中でExpressionを使った場合にコンパイラがエラーを出力できるようにするために記述している
}

type Expression interface {
	Node
	expressionNode() // この関数はダミーメソッドで、Expressionの中でStatementを使った場合にコンパイラがエラーを出力できるようにするために記述している
}

type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

// let文の定義
type LetStatement struct {
	Token token.Token // <expression>のトークン
	Name  *Identifier // <expression>に変数名が入るけ０素
	Value Expression  // <expression>に評価した結果の値が入るケース
}

func (ls *LetStatement) statementNode()       {}
func (ls *LetStatement) TokenLiteral() string { return ls.Token.Literal }

type ReturnStatement struct {
	Token       token.Token // 'return' トークン
	ReturnValue Expression
}

func (rs *ReturnStatement) statementNode()       {}
func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Literal }

// Monkeyにおける識別子は値を生成すると定義している (つまり、識別子はStatementではなくExpression)
type Identifier struct {
	Token token.Token
	Value string
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
