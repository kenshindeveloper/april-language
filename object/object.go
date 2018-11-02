package object

import (
	"bytes"
	"fmt"
	"hash/fnv"
	"os"
	"strings"

	"github.com/kenshindeveloper/april/ast"
)

const (
	ERROR_OBJ    = "ERROR"
	INTEGER_OBJ  = "INTEGER"
	BOOLEAN_OBJ  = "BOOLEAN"
	DOUBLE_OBJ   = "DOUBLE"
	STRING_OBJ   = "STRING"
	NIL_OBJ      = "NIL"
	BUILTIN_OBJ  = "BUILTIN"
	RETURN_OBJ   = "RETURN"
	BREAK_OBJ    = "BREAK"
	CLOSURE_OBJ  = "CLOSURE"
	FUNCTION_OBJ = "FUNCTION"
	LIST_OBJ     = "LIST"
	HASH_OBJ     = "HASH"
	STREAM_OBJ   = "STREAM"
	STRUCT_OBJ   = "STRUCT"
)

type ObjectType string

type Object interface {
	Type() ObjectType
	Inspect() string
}

type Hashable interface {
	HashKey() HashKey
}

func GetType(name ObjectType) string {
	switch name {
	case ERROR_OBJ:
		return "error"
	case INTEGER_OBJ:
		return "int"
	case BOOLEAN_OBJ:
		return "bool"
	case DOUBLE_OBJ:
		return "double"
	case STRING_OBJ:
		return "string"
	case NIL_OBJ:
		return "null"
	case RETURN_OBJ:
		return "return"
	case BREAK_OBJ:
		return "break"
	case BUILTIN_OBJ:
		return "function"
	case CLOSURE_OBJ:
		return "closure"
	case LIST_OBJ:
		return "list"
	case HASH_OBJ:
		return "map"
	case STREAM_OBJ:
		return "stream"
	default:
		return "null"
	}
}

//***************************************************************************************
//***************************************************************************************
//***************************************************************************************

type Integer struct {
	Value int64
}

func (i *Integer) Type() ObjectType {
	return INTEGER_OBJ
}

func (i *Integer) Inspect() string {
	return fmt.Sprintf("%d", i.Value)
}

//***************************************************************************************
//***************************************************************************************
//***************************************************************************************

type Boolean struct {
	Value bool
}

func (b *Boolean) Type() ObjectType {
	return BOOLEAN_OBJ
}

func (b *Boolean) Inspect() string {
	return fmt.Sprintf("%t", b.Value)
}

//***************************************************************************************
//***************************************************************************************
//***************************************************************************************

type String struct {
	Value string
}

func (s *String) Type() ObjectType {
	return STRING_OBJ
}

func (s *String) Inspect() string {
	return "'" + s.Value + "'"
}

//***************************************************************************************
//***************************************************************************************
//***************************************************************************************

type Double struct {
	Value float64
}

func (d *Double) Type() ObjectType {
	return DOUBLE_OBJ
}

func (d *Double) Inspect() string {
	return fmt.Sprintf("%g", d.Value)
}

//***************************************************************************************
//***************************************************************************************
//***************************************************************************************

type Nil struct{}

func (n *Nil) Type() ObjectType {
	return NIL_OBJ
}

func (n *Nil) Inspect() string {
	return "null"
}

//***************************************************************************************
//***************************************************************************************
//***************************************************************************************

type Error struct {
	Message string
}

func (e *Error) Type() ObjectType {
	return ERROR_OBJ
}

func (e *Error) Inspect() string {
	return "Error: " + e.Message
}

//***************************************************************************************
//***************************************************************************************
//***************************************************************************************

type BuiltinFunction func(args ...Object) Object

type Builtin struct {
	Fn BuiltinFunction
}

func (b *Builtin) Type() ObjectType {
	return BUILTIN_OBJ
}

func (b *Builtin) Inspect() string {
	return "builtin function"
}

//***************************************************************************************
//***************************************************************************************
//***************************************************************************************

type ReturnStatement struct {
	Value Object
}

func (rs *ReturnStatement) Type() ObjectType {
	return RETURN_OBJ
}

func (rs *ReturnStatement) Inspect() string {
	return rs.Value.Inspect()
}

//***************************************************************************************
//***************************************************************************************
//***************************************************************************************

type BreakStatement struct {
}

func (b *BreakStatement) Type() ObjectType {
	return BREAK_OBJ
}

func (b *BreakStatement) Inspect() string {
	return "break"
}

//***************************************************************************************
//***************************************************************************************
//***************************************************************************************

type FunctionClosure struct {
	Parameters []*ast.FunctionParameters
	Return     *ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment
}

func (fc *FunctionClosure) Type() ObjectType {
	return CLOSURE_OBJ
}

func (fc *FunctionClosure) Inspect() string {
	var out bytes.Buffer

	params := []string{}
	for _, p := range fc.Parameters {
		params = append(params, p.String())
	}

	out.WriteString("fn")
	out.WriteString("(")
	out.WriteString(strings.Join(params, ","))
	out.WriteString(")")
	out.WriteString("{ ")
	out.WriteString(fc.Body.String())
	out.WriteString(" }")

	return out.String()
}

//***************************************************************************************
//***************************************************************************************
//***************************************************************************************

type List struct {
	Elements []Object
}

func (l *List) Type() ObjectType {
	return LIST_OBJ
}

func (l *List) Inspect() string {
	var out bytes.Buffer

	elements := []string{}
	for _, element := range l.Elements {
		elements = append(elements, element.Inspect())
	}

	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")

	return out.String()
}

//***************************************************************************************
//***************************************************************************************
//***************************************************************************************

type HashKey struct {
	Type  ObjectType
	Value uint64
}

func (b *Boolean) HashKey() HashKey {
	var value uint64

	if b.Value {
		value = 1
	} else {
		value = 0
	}

	return HashKey{Type: b.Type(), Value: value}
}

func (i *Integer) HashKey() HashKey {
	return HashKey{Type: i.Type(), Value: uint64(i.Value)}
}

func (d *Double) HashKey() HashKey {
	return HashKey{Type: d.Type(), Value: uint64(d.Value)}
}

func (s *String) HashKey() HashKey {
	h := fnv.New64a()
	h.Write([]byte(s.Value))
	return HashKey{Type: s.Type(), Value: h.Sum64()}
}

//***************************************************************************************
//***************************************************************************************
//***************************************************************************************

type HashPair struct {
	Key   Object
	Value Object
}

type Hash struct {
	Pairs map[HashKey]HashPair
}

func (h *Hash) Type() ObjectType {
	return HASH_OBJ
}

func (h *Hash) Inspect() string {
	var out bytes.Buffer

	pairs := []string{}
	for _, pair := range h.Pairs {
		pairs = append(pairs, fmt.Sprintf("%s: %s", pair.Key.Inspect(), pair.Value.Inspect()))
	}

	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")

	return out.String()
}

//***************************************************************************************
//***************************************************************************************
//***************************************************************************************

type Function struct {
	Name       *ast.Identifier
	Parameters []*ast.FunctionParameters
	Return     *ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment
}

func (fc *Function) Type() ObjectType {
	return CLOSURE_OBJ
}

func (fc *Function) Inspect() string {
	var out bytes.Buffer

	params := []string{}
	for _, p := range fc.Parameters {
		params = append(params, p.String())
	}

	out.WriteString("fn")
	out.WriteString("(")
	out.WriteString(strings.Join(params, ","))
	out.WriteString(")")
	out.WriteString("{ ")
	out.WriteString(fc.Body.String())
	out.WriteString(" }")

	return out.String()
}

//***************************************************************************************
//***************************************************************************************
//***************************************************************************************

type Stream struct {
	FILE *os.File
}

func (s *Stream) Type() ObjectType {
	return STREAM_OBJ
}

func (s *Stream) Inspect() string {
	if s.FILE != nil {
		return s.FILE.Name()
	}

	return "null"
}

//***************************************************************************************
//***************************************************************************************
//***************************************************************************************

type Struct struct {
	Env *Environment
}

func (s *Struct) Type() ObjectType {
	return STRUCT_OBJ
}

func (s *Struct) Inspect() string {
	return "struct"
}
