package ast

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/kenshindeveloper/april/token"
)

type Node interface {
	TokenLiteral() string
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
	result := ""

	if len(p.Statements) > 0 {
		return p.Statements[0].String()
	}

	return result
}

func (p *Program) String() string {
	var out bytes.Buffer

	for _, stmt := range p.Statements {
		out.WriteString(stmt.String())
	}

	return out.String()
}

//***************************************************************************************
//***************************************************************************************
//***************************************************************************************

type Identifier struct {
	Token token.Token
	Name  string
	Line  int
}

func (i *Identifier) expressionNode() {}

func (i *Identifier) TokenLiteral() string {
	return i.Token.Literal
}

func (i *Identifier) String() string {
	return i.Token.Literal
}

//***************************************************************************************
//***************************************************************************************
//***************************************************************************************

type Nil struct {
	Token token.Token
	Name  string
	Line  int
}

func (n *Nil) expressionNode() {}

func (n *Nil) TokenLiteral() string {
	return n.Token.Literal
}

func (n *Nil) String() string {
	return n.Token.Literal
}

//***************************************************************************************
//***************************************************************************************
//***************************************************************************************

type VarStatement struct {
	Token token.Token //var
	Name  *Identifier //x
	Type  *Identifier // int
	Value Expression  // 5
	Line  int
}

func (vs *VarStatement) statementNode() {}

func (vs *VarStatement) TokenLiteral() string {
	return vs.Token.Literal
}

func (vs *VarStatement) String() string {
	var out bytes.Buffer

	out.WriteString(vs.Token.Literal)
	out.WriteString(" " + vs.Name.String() + " ")
	out.WriteString(vs.Type.String())

	if vs.Value != nil {
		out.WriteString(" = ")
		out.WriteString(vs.Value.String())
		out.WriteString(";")
	}

	return out.String()
}

//***************************************************************************************
//***************************************************************************************
//***************************************************************************************

type GlobalStatement struct {
	Token token.Token //var
	Name  *Identifier //x
	Type  *Identifier // int
	Value Expression  // 5
	Line  int
}

func (gs *GlobalStatement) statementNode() {}

func (gs *GlobalStatement) TokenLiteral() string {
	return gs.Token.Literal
}

func (gs *GlobalStatement) String() string {
	var out bytes.Buffer

	out.WriteString(gs.Token.Literal)
	out.WriteString(" " + gs.Name.String() + " ")
	out.WriteString(gs.Type.String())

	if gs.Value != nil {
		out.WriteString(" = ")
		out.WriteString(gs.Value.String())
		out.WriteString(";")
	}

	return out.String()
}

//***************************************************************************************
//***************************************************************************************
//***************************************************************************************

type ExpressionStatement struct {
	Token      token.Token
	Expression Expression
	Line       int
}

func (es *ExpressionStatement) statementNode() {}

func (es *ExpressionStatement) TokenLiteral() string {
	return es.Token.Literal
}

func (es *ExpressionStatement) String() string {
	var out bytes.Buffer

	if es.Expression != nil {
		out.WriteString(es.Expression.String())
	}
	return out.String()
}

//***************************************************************************************
//***************************************************************************************
//***************************************************************************************
type Integer struct {
	Token token.Token
	Value int64
	Line  int
}

func (i *Integer) expressionNode() {}

func (i *Integer) TokenLiteral() string {
	return i.Token.Literal
}

func (i *Integer) String() string {
	return i.Token.Literal
}

//***************************************************************************************
//***************************************************************************************
//***************************************************************************************

type Boolean struct {
	Token token.Token
	Value bool
	Line  int
}

func (b *Boolean) expressionNode() {}

func (b *Boolean) TokenLiteral() string {
	return b.Token.Literal
}

func (b *Boolean) String() string {
	return fmt.Sprintf("%t", b.Value)
}

//***************************************************************************************
//***************************************************************************************
//***************************************************************************************

type String struct {
	Token token.Token
	Value string
	Line  int
}

func (s *String) expressionNode() {}

func (s *String) TokenLiteral() string {
	return s.Token.Literal
}

func (s *String) String() string {
	return s.Value
}

//***************************************************************************************
//***************************************************************************************
//***************************************************************************************

type Double struct {
	Token token.Token
	Value float64
	Line  int
}

func (d *Double) expressionNode() {}

func (d *Double) TokenLiteral() string {
	return d.Token.Literal
}

func (d *Double) String() string {
	return fmt.Sprintf("%g", d.Value)
}

//***************************************************************************************
//***************************************************************************************
//***************************************************************************************

type PrefixExpression struct {
	Token    token.Token
	Operator string
	Right    Expression
	Line     int
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
	Line     int
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

type ImplicitDeclarationExpression struct {
	Token token.Token
	Left  *Identifier
	Right Expression
	Line  int
}

func (ide *ImplicitDeclarationExpression) expressionNode() {}

func (ide *ImplicitDeclarationExpression) TokenLiteral() string {
	return ide.Token.Literal
}

func (ide *ImplicitDeclarationExpression) String() string {
	var out bytes.Buffer

	out.WriteString(ide.Left.String())
	out.WriteString(" := ")
	out.WriteString(ide.Right.String())

	return out.String()
}

//***************************************************************************************
//***************************************************************************************
//***************************************************************************************

type AssignExpression struct {
	Token token.Token
	Left  Expression
	Right Expression
	Line  int
}

func (ae *AssignExpression) expressionNode() {}

func (ae *AssignExpression) TokenLiteral() string {
	return ae.Token.Literal
}

func (ae *AssignExpression) String() string {
	var out bytes.Buffer

	out.WriteString(ae.Left.String())
	out.WriteString(" = ")
	out.WriteString(ae.Right.String())

	return out.String()
}

//***************************************************************************************
//***************************************************************************************
//***************************************************************************************

type AssignOperationExpression struct {
	Token    token.Token
	Left     *Identifier
	Operator string
	Right    Expression
	Line     int
}

func (aoe *AssignOperationExpression) expressionNode() {}

func (aoe *AssignOperationExpression) TokenLiteral() string {
	return aoe.Token.Literal
}

func (aoe *AssignOperationExpression) String() string {
	var out bytes.Buffer

	out.WriteString(aoe.Left.String())
	out.WriteString(" " + aoe.Operator + " ")
	out.WriteString(aoe.Right.String())

	return out.String()
}

//***************************************************************************************
//***************************************************************************************
//***************************************************************************************

type ReturnStatement struct {
	Token      token.Token
	Expression Expression
	Line       int
}

func (rs *ReturnStatement) statementNode() {}

func (rs *ReturnStatement) TokenLiteral() string {
	return rs.Token.Literal
}

func (rs *ReturnStatement) String() string {
	var out bytes.Buffer

	if rs.Expression != nil {
		out.WriteString(rs.Expression.String())
	}

	return out.String()
}

//***************************************************************************************
//***************************************************************************************
//***************************************************************************************

type BreakStatement struct {
	Token token.Token
	Line  int
}

func (b *BreakStatement) statementNode() {}

func (b *BreakStatement) TokenLiteral() string {
	return b.Token.Literal
}

func (b *BreakStatement) String() string {
	var out bytes.Buffer

	out.WriteString(b.TokenLiteral())

	return out.String()
}

//***************************************************************************************
//***************************************************************************************
//***************************************************************************************

//0:fn - 1:if - 2:for
type BlockStatement struct {
	Token      token.Token
	Statements []Statement
	Type       int
	Line       int
}

func (bs *BlockStatement) statementNode() {}

func (bs *BlockStatement) TokenLiteral() string {
	return bs.Token.Literal
}

func (bs *BlockStatement) String() string {
	var out bytes.Buffer

	for _, data := range bs.Statements {
		out.WriteString(data.String())
	}

	return out.String()
}

//***************************************************************************************
//***************************************************************************************
//***************************************************************************************
type IfExpression struct {
	Token       token.Token
	Codition    Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
	Line        int
}

func (ie *IfExpression) expressionNode() {}

func (ie *IfExpression) TokenLiteral() string {
	return ie.Token.Literal
}

func (ie *IfExpression) String() string {
	var out bytes.Buffer

	out.WriteString("if")
	out.WriteString(ie.Codition.String())
	out.WriteString(" ")
	out.WriteString(ie.Consequence.String())

	if ie.Alternative != nil {
		out.WriteString("else")
		out.WriteString(ie.Alternative.String())
	}

	return out.String()
}

//***************************************************************************************
//***************************************************************************************
//***************************************************************************************

type FunctionClosure struct {
	Token      token.Token
	Parameters []*FunctionParameters
	Body       *BlockStatement
	Type       *Identifier
	Line       int
}

func (fc *FunctionClosure) expressionNode() {}

func (fc *FunctionClosure) TokenLiteral() string {
	return fc.Token.Literal
}

func (fc *FunctionClosure) String() string {
	var out bytes.Buffer

	params := []string{}
	for _, p := range fc.Parameters {
		params = append(params, p.String())
	}

	out.WriteString(fc.TokenLiteral())
	out.WriteString("(")
	out.WriteString(strings.Join(params, ","))
	out.WriteString(")")
	out.WriteString(fc.Type.String())
	out.WriteString(fc.Body.String())

	return out.String()
}

//***************************************************************************************
//***************************************************************************************
//***************************************************************************************

type CallExpression struct {
	Token     token.Token //'('
	Function  Expression  //Identifier o FunctionClosure
	Arguments []Expression
	Line      int
}

func (ce *CallExpression) expressionNode() {}

func (ce *CallExpression) TokenLiteral() string {
	return ce.Token.Literal
}
func (ce *CallExpression) String() string {
	var out bytes.Buffer

	args := []string{}

	for _, arg := range ce.Arguments {
		args = append(args, arg.String())
	}

	out.WriteString(ce.Function.String())
	out.WriteString("(")
	out.WriteString(strings.Join(args, ","))
	out.WriteString(")")

	return out.String()
}

//***************************************************************************************
//***************************************************************************************
//***************************************************************************************

type List struct {
	Token    token.Token
	Elements []Expression
	Line     int
}

func (l *List) expressionNode() {}

func (l *List) TokenLiteral() string {
	return l.Token.Literal
}

func (l *List) String() string {
	var out bytes.Buffer

	elements := []string{}

	for _, element := range l.Elements {
		elements = append(elements, element.String())
	}

	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")

	return out.String()
}

//***************************************************************************************
//***************************************************************************************
//***************************************************************************************

type IndexExpression struct {
	Token token.Token
	Left  Expression
	Index Expression
	Line  int
}

func (ie *IndexExpression) expressionNode() {}

func (ie *IndexExpression) TokenLiteral() string {
	return ie.Token.Literal
}

func (ie *IndexExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString("[")
	out.WriteString(ie.Index.String())
	out.WriteString("])")

	return out.String()
}

//***************************************************************************************
//***************************************************************************************
//***************************************************************************************

type Hash struct {
	Token token.Token
	Pairs map[Expression]Expression
	Line  int
}

func (h *Hash) expressionNode() {}

func (h *Hash) TokenLiteral() string {
	return h.Token.Literal
}

func (h *Hash) String() string {
	var out bytes.Buffer

	pairs := []string{}
	for key, value := range h.Pairs {
		pairs = append(pairs, key.String()+":"+value.String())
	}

	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")

	return out.String()
}

//***************************************************************************************
//***************************************************************************************
//***************************************************************************************

type PostfixExpression struct {
	Token    token.Token
	Left     *Identifier
	Operator string
	Line     int
}

func (pe *PostfixExpression) expressionNode() {}

func (pe *PostfixExpression) TokenLiteral() string {
	return pe.Token.Literal
}

func (pe *PostfixExpression) String() string {
	var out bytes.Buffer

	out.WriteString(pe.Left.String())
	out.WriteString(pe.Operator)

	return out.String()
}

//***************************************************************************************
//***************************************************************************************
//***************************************************************************************

type ForStatement struct {
	Token       token.Token
	Declaration Expression
	Condition   Expression
	Operation   Expression
	Body        *BlockStatement
	Line        int
}

func (fs *ForStatement) statementNode() {}

func (fs *ForStatement) TokenLiteral() string {
	return fs.Token.Literal
}

func (fs *ForStatement) String() string {
	var out bytes.Buffer

	out.WriteString("for")
	out.WriteString(" (")
	out.WriteString(fs.Condition.String())
	out.WriteString(" ) ")
	out.WriteString(" {")
	out.WriteString(fs.Body.String())
	out.WriteString(" }")

	return out.String()
}

//***************************************************************************************
//***************************************************************************************
//***************************************************************************************

type FunctionParameters struct {
	Token token.Token
	Name  *Identifier
	Type  *Identifier
	Line  int
}

func (fp *FunctionParameters) expressionNode() {}

func (fp *FunctionParameters) TokenLiteral() string {
	return fp.Token.Literal
}

func (fp *FunctionParameters) String() string {
	var out bytes.Buffer

	out.WriteString(fp.Name.Name)
	out.WriteString(":")
	out.WriteString(fp.Type.Name)

	return out.String()
}

//***************************************************************************************
//***************************************************************************************
//***************************************************************************************

type Function struct {
	Token      token.Token
	Name       *Identifier
	Parameters []*FunctionParameters
	Type       *Identifier
	Body       *BlockStatement
	Line       int
}

func (fn *Function) statementNode() {}

func (fn *Function) TokenLiteral() string {
	return fn.Token.Literal
}

func (fn *Function) String() string {
	var out bytes.Buffer

	params := []string{}
	for _, p := range fn.Parameters {
		params = append(params, p.String())
	}

	out.WriteString(fn.TokenLiteral())
	out.WriteString(fn.Name.String())
	out.WriteString(" (")
	out.WriteString(strings.Join(params, ","))
	out.WriteString(") ")
	out.WriteString(fn.Type.String())
	out.WriteString(fn.Body.String())

	return out.String()
}

//***************************************************************************************
//***************************************************************************************
//***************************************************************************************

type Stream struct {
	Token token.Token
	FILE  *os.File
	Line  int
}

func (s *Stream) expressionNode() {}

func (s *Stream) TokenLiteral() string {
	return s.Token.Literal
}

func (s *Stream) String() string {
	return s.FILE.Name()
}

//***************************************************************************************
//***************************************************************************************
//***************************************************************************************

type Struct struct {
	Token   token.Token
	Element map[*Identifier]*Identifier
	Line    int
}

func (s *Struct) expressionNode() {}

func (s *Struct) TokenLiteral() string {
	return s.Token.Literal
}

func (s *Struct) String() string {
	var out bytes.Buffer

	elements := []string{}
	for key, value := range s.Element {
		elements = append(elements, key.String()+":"+value.String())
	}

	out.WriteString("{")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("}")

	return out.String()
}
