package ast

import (
	"bytes"

	"github.com/kenshindeveloper/april/token"
)

type Node interface {
	TokenLiteral() string //su uso es solo de prueba..
	String() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

//***************************************************************************************
//***************************************************************************************
//***************************************************************************************

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

func (p *Program) String() string {
	var out bytes.Buffer

	for _, s := range p.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}

//***************************************************************************************
//***************************************************************************************
//***************************************************************************************

//Identifier es una
type Identifier struct {
	Token token.Token
	Value string
}

func (ident *Identifier) expressionNode() {}

func (ident *Identifier) TokenLiteral() string {
	return ident.Token.Literal
}

func (ident *Identifier) String() string {
	return ident.Value
}

//***************************************************************************************
//***************************************************************************************
//***************************************************************************************

//VarStatement es una estructura contenedora ejemplo:var x int = 15;
//OJO: falta declarar el attr de -> Type <-
type VarStatement struct {
	Token token.Token //var
	Name  *Identifier //x
	Value Expression  //15
}

func (varStmt *VarStatement) statementNode() {}

func (varStmt *VarStatement) TokenLiteral() string {
	return varStmt.Token.Literal
}

func (varStmt *VarStatement) String() string {
	var out bytes.Buffer

	out.WriteString(varStmt.TokenLiteral() + " ")
	out.WriteString(varStmt.Name.String())
	out.WriteString(" = ")

	if varStmt.Value != nil {
		out.WriteString(varStmt.Value.String())
	}

	out.WriteString(";")

	return out.String()
}

//***************************************************************************************
//***************************************************************************************
//***************************************************************************************

//ReturnStatement es una estructura que almacena un Statement
type ReturnStatement struct {
	Token       token.Token
	ReturnValue Expression
}

func (r *ReturnStatement) statementNode() {}

func (r *ReturnStatement) TokenLiteral() string {
	return r.Token.Literal
}

func (r *ReturnStatement) String() string {
	var out bytes.Buffer

	out.WriteString(r.TokenLiteral() + " ")

	if r.ReturnValue != nil {
		out.WriteString(r.ReturnValue.String())
	}

	out.WriteString(";")
	return out.String()
}

//***************************************************************************************
//***************************************************************************************
//***************************************************************************************

type ExpressionStatement struct {
	Token      token.Token //El primer token de la expression
	Expression Expression
}

func (expreStmt *ExpressionStatement) statementNode() {}

func (exprStmt *ExpressionStatement) TokenLiteral() string {
	return exprStmt.Token.Literal
}

func (exprStmt *ExpressionStatement) String() string {
	if exprStmt.Expression != nil {
		return exprStmt.Expression.String()
	}
	return ""
}

//***************************************************************************************
//***************************************************************************************
//***************************************************************************************

type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (i *IntegerLiteral) expressionNode() {}

func (i *IntegerLiteral) TokenLiteral() string {
	return i.Token.Literal
}

func (i *IntegerLiteral) String() string {
	return i.Token.Literal
}

//***************************************************************************************
//***************************************************************************************
//***************************************************************************************

type PrefixExpression struct {
	Token    token.Token
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) expressionNode() {}

func (pe *PrefixExpression) TokenLiteral() string {
	return pe.Token.Literal
}

func (pe *PrefixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String())
	out.WriteString(")")

	return out.String()
}

//***************************************************************************************
//***************************************************************************************
//***************************************************************************************

type InfixExpression struct {
	Token    token.Token
	Left     Expression
	Operator string
	Right    Expression
}

func (ie *InfixExpression) expressionNode() {}

func (ie *InfixExpression) TokenLiteral() string {
	return ie.Token.Literal
}

func (ie *InfixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString(" " + ie.Operator + " ")
	out.WriteString(ie.Right.String())
	out.WriteString(")")

	return out.String()
}
