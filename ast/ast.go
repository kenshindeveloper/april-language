package ast

import "github.com/kenshindeveloper/april/token"

type Node interface {
	TokenLiteral() string //su uso es solo de prueba..
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
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

//Identifier es una
type Identifier struct {
	Token token.Token
	Value string
}

func (ident *Identifier) expressionNode() {}

func (ident *Identifier) TokenLiteral() string {
	return ident.Token.Literal
}

//!Estructura Identifier

//VarStatement es una estructura contenedora de la declaracion var ejemplo:var x int = 15;
type VarStatement struct {
	Token token.Token
	Name  *Identifier
	Value Expression
}

func (ls *VarStatement) statementNode() {}

func (ls *VarStatement) TokenLiteral() string {
	return ls.Token.Literal
}
