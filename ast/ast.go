package ast

import (
	"bytes"
	"strings"

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

func (es *ExpressionStatement) statementNode() {}

func (es *ExpressionStatement) TokenLiteral() string {
	return es.Token.Literal
}

func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
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

//***************************************************************************************
//***************************************************************************************
//***************************************************************************************

type Boolean struct {
	Token token.Token
	Value bool
}

func (b *Boolean) expressionNode() {}

func (b *Boolean) TokenLiteral() string {
	return b.Token.Literal
}

func (b *Boolean) String() string {
	return b.Token.Literal
}

//***************************************************************************************
//***************************************************************************************
//***************************************************************************************

type IfExpression struct {
	Token       token.Token
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (ife *IfExpression) expressionNode() {}

func (ife *IfExpression) TokenLiteral() string {
	return ife.Token.Literal
}

func (ife *IfExpression) String() string {
	var out bytes.Buffer

	out.WriteString("if ")
	out.WriteString(ife.Condition.String())
	out.WriteString(" ")
	out.WriteString(ife.Consequence.String())

	if ife.Alternative != nil {
		out.WriteString("else ")
		out.WriteString(ife.Alternative.String())
	}

	return out.String()
}

//***************************************************************************************
//***************************************************************************************
//***************************************************************************************

type BlockStatement struct {
	Token      token.Token
	Statements []Statement
}

func (bs *BlockStatement) statementNode() {}

func (bs *BlockStatement) TokenLiteral() string {
	return bs.Token.Literal
}

func (bs *BlockStatement) String() string {
	var out bytes.Buffer

	for _, s := range bs.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}

//***************************************************************************************
//***************************************************************************************
//***************************************************************************************

type FunctionLiteral struct {
	Token      token.Token //el token 'fn'
	Parameters []*Identifier
	Body       *BlockStatement
}

func (fl *FunctionLiteral) expressionNode() {}

func (fl *FunctionLiteral) TokenLiteral() string {
	return fl.Token.Literal
}

func (fl *FunctionLiteral) String() string {
	var out bytes.Buffer

	parameters := []string{}

	for _, data := range fl.Parameters {
		parameters = append(parameters, data.String())
	}

	out.WriteString(fl.TokenLiteral())
	out.WriteString("(")
	out.WriteString(strings.Join(parameters, ", "))
	out.WriteString(")")
	out.WriteString("{")
	out.WriteString(fl.Body.String())
	out.WriteString("}")

	return out.String()
}

//***************************************************************************************
//***************************************************************************************
//***************************************************************************************

type CallExpression struct {
	Token     token.Token //the '(' token
	Function  Expression  // Identifier or FunctionLiteral
	Arguments []Expression
}

func (ce *CallExpression) expressionNode() {}

func (ce *CallExpression) TokenLiteral() string {
	return ce.Token.Literal
}

func (ce *CallExpression) String() string {
	var out bytes.Buffer

	args := []string{}

	for _, para := range ce.Arguments {
		args = append(args, para.String())
	}

	out.WriteString(ce.Function.String())
	out.WriteString("(")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")

	return out.String()
}
