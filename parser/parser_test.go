package parser

import (
	"fmt"
	"testing"

	"github.com/kenshindeveloper/april/ast"
	"github.com/kenshindeveloper/april/lexer"
)

func TestStructExpression(t *testing.T) {
	input := `{x:int, y:string}`

	l := lexer.New(input)
	p := New(l)
	program := p.ParserProgram()
	checkParserErrors(t, p)
	if program == nil {
		t.Fatalf("program is null")
	}

	if len(program.Statements) != 1 {
		t.Fatalf("len(program.Statements) is not equal '%d'. got='%d'", 1, len(program.Statements))
	}

	_, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not equal '*ast.ExpressionStatement'. got='%T'", program.Statements[0])
	}
}

func TestGlobalStatement(t *testing.T) {
	tests := []struct {
		input string
	}{
		{"global t:int = 15;"},
		{`global t:string = "15";`},
		{"global t:double = 1;"},
		{"global t:bool = false;"},
	}

	for _, data := range tests {
		l := lexer.New(data.input)
		p := New(l)
		program := p.ParserProgram()
		checkParserErrors(t, p)
		if program == nil {
			t.Fatalf("program is null")
		}

		if len(program.Statements) != 1 {
			t.Fatalf("len(program.Statements) is not equal '%d'. got='%d'", 1, len(program.Statements))
		}

		_, ok := program.Statements[0].(*ast.GlobalStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not equal '*ast.ExpressionStatement'. got='%T'", program.Statements[0])
		}
	}

}

func TestNilExpression(t *testing.T) {
	input := "nil"

	l := lexer.New(input)
	p := New(l)
	program := p.ParserProgram()
	checkParserErrors(t, p)
	if program == nil {
		t.Fatalf("program is null")
	}

	if len(program.Statements) != 1 {
		t.Fatalf("len(program.Statements) is not equal '%d'. got='%d'", 1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not equal '*ast.ExpressionStatement'. got='%T'", program.Statements[0])
	}

	n, ok := stmt.Expression.(*ast.Nil)
	if !ok {
		t.Fatalf("stmt.Expression is not equal '*ast.Nil'. got='%T'", stmt.Expression)
	}

	if n.Name != "nil" {
		t.Fatalf("n.Name is not equal to 'nil'. got='%s'", n.Name)
	}

}

func TestFunctionStatement(t *testing.T) {

	tests := []struct {
		input string
	}{
		{"fn foo(){}"},
		{"fn foo(a:int) {}"},
		{"fn foo(a:int) string{}"},
		{"fn foo(a:int, b:double) double{}"},
		{"fn foo(a:string, b:string, c:string) string{}"},
	}

	for _, data := range tests {
		l := lexer.New(data.input)
		p := New(l)
		program := p.ParserProgram()
		checkParserErrors(t, p)
		if program == nil {
			t.Fatalf("program is null")
		}

		if len(program.Statements) != 1 {
			t.Fatalf("len(program.Statements) is not equal '%d'. got='%d'", 1, len(program.Statements))
		}

		_, ok := program.Statements[0].(*ast.Function)
		if !ok {
			t.Fatalf("program.Statements[0] is not equal '*ast.Function'. got='%T'", program.Statements[0])
		}
	}
}

func TestFunctionClosure(t *testing.T) {
	tests := []struct {
		input          string
		expectedParams []string
	}{
		{"x := fn() {};", []string{}},
		{"x := fn(x:int) {};", []string{"x"}},
		{"x := fn(x:int, y:int, z:int) int {};", []string{"x", "y", "z"}},
		{"x := fn(x:int, y:double, z:bool, w:map, v:list, k:int) double {};", []string{"x", "y", "z", "w", "v", "k"}},
	}

	for _, data := range tests {
		l := lexer.New(data.input)
		p := New(l)
		program := p.ParserProgram()
		checkParserErrors(t, p)
		if program == nil {
			t.Fatalf("program is null.")
		}

		if len(program.Statements) != 1 {
			t.Fatalf("len(program.Statement) is not equal '%d'. got='%d'", 1, len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("stmt is not equal '*ast.ExpressionStatement'. got='%T'", program.Statements[0])
		}

		decla, ok := stmt.Expression.(*ast.ImplicitDeclarationExpression)
		if !ok {
			t.Fatalf("stmt is not equal '*ast.ImplicitDeclarationExpression'. got='%T'", decla)
		}

		function, ok := decla.Right.(*ast.FunctionClosure)
		if !ok {
			t.Fatalf("stmt is not equal '*ast.FunctionClosure'. got='%T'", function)
		}

		// if len(function.Parameters) != len(data.expectedParams) {
		// 	t.Fatalf("length parameters wrong. want %d, got=%d", len(function.Parameters), len(data.expectedParams))
		// }

		// for i, ident := range data.expectedParams {
		// 	if function.Parameters[i].Name != ident {
		// 		t.Fatalf("function.Paramentes[%d] is not equal '%s'. got='%s'", i, ident, function.Parameters[i])
		// 	}
		// }
	}
}

func TestParsingImportExpression(t *testing.T) {
	input := `
		import "path/test.april"
	`

	l := lexer.New(input)
	p := New(l)
	program := p.ParserProgram()
	if program == nil {
		t.Fatalf("program is null")
	}

	if len(program.Statements) != 1 {
		t.Fatalf("len(program.Statements) is not equal '%d'. got='%d'", 1, len(program.Statements))
	}
}

func TestParsingForExpression(t *testing.T) {
	tests := []struct {
		input string
	}{
		{"for (x := 1; x < 10; x++) {}"},
		{"for (true) {}"},
	}

	for _, data := range tests {
		l := lexer.New(data.input)
		p := New(l)
		program := p.ParserProgram()
		if program == nil {
			t.Fatalf("program is null")
		}

		if len(program.Statements) != 1 {
			t.Fatalf("len(program.Statements) is not equal '%d'. got='%d'", 1, len(program.Statements))
		}

		_, ok := program.Statements[0].(*ast.ForStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not equal '*ast.ForStatement'. got='%T'", program.Statements[0])
		}
	}
}

func TestParsingPostfix(t *testing.T) {
	tests := []struct {
		input    string
		ident    string
		operator string
	}{
		{"x++;", "x", "++"},
		{"x--;", "x", "--"},
	}

	for _, data := range tests {
		l := lexer.New(data.input)
		p := New(l)
		program := p.ParserProgram()
		if program == nil {
			t.Fatalf("program is null")
		}

		if len(program.Statements) != 1 {
			t.Fatalf("len(program.Statements) is not equal '%d'. got='%d'", 1, len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not equal '*ast.ExpressionStatement'. got='%T'", program.Statements[0])
		}

		postfix, ok := stmt.Expression.(*ast.PostfixExpression)
		if !ok {
			t.Fatalf("stmt.Expression is not equal '*ast.PostfixExpression'. got='%T'", stmt.Expression)
		}

		ident := postfix.Left
		if ident.Name != data.ident {
			t.Fatalf("ident.Name is not equal '%s'. got='%s'", data.ident, ident.Name)
		}

		if postfix.Operator != data.operator {
			t.Fatalf("postfix.Operator is not equal '%s'. got='%s'", data.operator, postfix.Operator)
		}
	}

}

func TestParsingHashExpression(t *testing.T) {
	input := `{ "one": 1, "two": 2, "three": 3 }`

	l := lexer.New(input)
	p := New(l)
	program := p.ParserProgram()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	hash, ok := stmt.Expression.(*ast.Hash)
	if !ok {
		t.Fatalf("expression is nor '*ast.Hash'. got='%T'", stmt.Expression)
	}

	if len(hash.Pairs) != 3 {
		t.Fatalf("hash.Pairs has wrong length. got='%d'", len(hash.Pairs))
	}

	expected := map[string]int64{
		"one":   1,
		"two":   2,
		"three": 3,
	}

	for key, value := range hash.Pairs {
		literal, ok := key.(*ast.String)
		if !ok {
			t.Errorf("key is not 'ast.String'. got='%T'", key)
		}

		expecValue := expected[literal.String()]
		testIntegerLiteral(t, value, expecValue)
	}
}

func TestIndexExpression(t *testing.T) {
	input := "miLista[1 + 1]"

	l := lexer.New(input)
	p := New(l)
	program := p.ParserProgram()
	checkParserErrors(t, p)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	ie, ok := stmt.Expression.(*ast.IndexExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not equal '*ast.IndexExpression'. got='%T'", stmt.Expression)
	}

	if !testIdentifier(t, ie.Left, "miLista") {
		return
	}

	if !testInfixExpression(t, ie.Index, 1, "+", 1) {
		return
	}
}

func TestParsingList(t *testing.T) {
	input := "[1, 2 * 2, 3 + 3]"

	l := lexer.New(input)
	p := New(l)
	program := p.ParserProgram()
	checkParserErrors(t, p)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	list, ok := stmt.Expression.(*ast.List)
	if !ok {
		t.Fatalf("stmt.Expression is not equal '*ast.Literal'. got='%T'", stmt.Expression)
	}
	if len(list.Elements) != 3 {
		t.Fatalf("len(list.Elements) is not equal '%d'. got='%d'", 3, len(list.Elements))
	}

	testIntegerLiteral(t, list.Elements[0], 1)
	testInfixExpression(t, list.Elements[1], 2, "*", 2)
	testInfixExpression(t, list.Elements[2], 3, "+", 3)
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"-a * b", "((-a) * b)"},
		{"not -a", "(not(-a))"},
		{"a + b + c", "((a + b) + c)"},
		{"a + b - c", "((a + b) - c)"},
		{"a * b * c", "((a * b) * c)"},
		{"a * b / c", "((a * b) / c)"},
		{"a + b / c", "(a + (b / c))"},
		{"a + b * c + d / e - f", "(((a + (b * c)) + (d / e)) - f)"},
		{"3 + 4; -5 * 5", "(3 + 4)((-5) * 5)"},
		{"5 > 4 == 3 < 4", "((5 > 4) == (3 < 4))"},
		{"5 < 4 != 3 > 4", "((5 < 4) != (3 > 4))"},
		{"3 + 4 * 5 == 3 * 1 + 4 * 5", "((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))"},
		{"3 + 4 * 5 == 3 * 1 + 4 * 5", "((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))"},
		{"true", "true"},
		{"false", "false"},
		{"3 > 5 == false", "((3 > 5) == false)"},
		{"3 < 5 == true", "((3 < 5) == true)"},
		{"1 + (2 + 3) + 4", "((1 + (2 + 3)) + 4)"},
		{"(5 + 5) * 2", "((5 + 5) * 2)"},
		{"2 / (5 + 5)", "(2 / (5 + 5))"},
		{"-(5 + 5)", "(-(5 + 5))"},
		{"not (true == true)", "(not(true == true))"},
		{"a + add(b * c) + d", "((a + add((b * c))) + d)"},
		{"add(a, b, 1, 2 * 3, 4 + 5, add(6, 7 * 8))", "add(a,b,1,(2 * 3),(4 + 5),add(6,(7 * 8)))"},
		{"add(a + b + c * d / f + g)", "add((((a + b) + ((c * d) / f)) + g))"},
		{"a * [1, 2, 3, 4][b * c] * d", "((a * ([1, 2, 3, 4][(b * c)])) * d)"},
		{"add(a * b[2], b[1], 2 * [1, 2][1])", "add((a * (b[2])),(b[1]),(2 * ([1, 2][1])))"},

		{"x += 1", "x += 1"},
		{"x -= 1", "x -= 1"},
		{"x *= 1", "x *= 1"},
		{"x /= 1", "x /= 1"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParserProgram()
		checkParserErrors(t, p)

		actual := program.String()
		if actual != tt.expected {
			t.Errorf("expected=%q, got=%q", tt.expected, actual)
		}
	}
}

func TestCallExpressionParsing(t *testing.T) {
	input := `add(1, 2 * 3, 4 + 5);`

	l := lexer.New(input)
	p := New(l)
	program := p.ParserProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not containt %d. got=%d", 1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not *ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	expr, ok := stmt.Expression.(*ast.CallExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not *ast.CallExpression. got=%T", stmt.Expression)
	}

	if !testIdentifier(t, expr.Function, "add") {
		return
	}

	if len(expr.Arguments) != 3 {
		t.Fatalf("wrong length of arguments. got=%d", len(expr.Arguments))
	}

	testLiteralExpression(t, expr.Arguments[0], 1)
	testInfixExpression(t, expr.Arguments[1], 2, "*", 3)
	testInfixExpression(t, expr.Arguments[2], 4, "+", 5)
}

func TestIfExpression(t *testing.T) {
	input := `if (x < y) {x;} else {y;}`

	l := lexer.New(input)
	p := New(l)
	program := p.ParserProgram()
	if program == nil {
		t.Fatalf("program is null.")
	}

	if len(program.Statements) != 1 {
		t.Fatalf("len(program.Statement) is not equal '%d'. got='%d'", 1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("stmt is not equal '*ast.ExpressionStatement'. got='%T'", program.Statements[0])
	}

	expr, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("stmt is not equal '*ast.IfExpression'. got='%T'", expr)
	}
}

func TestImplicitDeclaration(t *testing.T) {
	tests := []struct {
		input         string
		nameExprected string
		expected      interface{}
	}{
		{"x := 6;", "x", 6},
		{"x := 156;", "x", 156},
		{"foo := true;", "foo", true},
		{"foo := false;", "foo", false},
	}

	for _, data := range tests {
		l := lexer.New(data.input)
		p := New(l)
		program := p.ParserProgram()
		if len(p.Error()) > 0 {
			t.Fatalf("Errors in ParserProgram.")
		}

		if program == nil {
			t.Fatalf("ParserProgram is null")
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statement is not '*ast.ExpressionStatement'. got='%T'", program.Statements)
		}

		ide, ok := stmt.Expression.(*ast.ImplicitDeclarationExpression)
		if !ok {
			t.Fatalf("stmt is not '*ast.ImplicitDeclarationExpression'. got='%T'", stmt)
		}

		if ide.Left.Name != data.nameExprected {
			t.Fatalf("ide.Left.Name is nor equal '%s'. got='%s'", data.nameExprected, ide.Left.Name)
		}

	}
}

func TestVarStatement(t *testing.T) {
	input := `
		var x:int = 15;
		var foo:int = 99;
		var y:int;
	`
	expected := []string{"x", "foo", "y"}

	l := lexer.New(input)
	p := New(l)
	program := p.ParserProgram()
	if program == nil {
		t.Fatalf("ParserProgram is equal nil.")
	}

	if len(program.Statements) != len(expected) {
		t.Fatalf("len(program.Statement) is not equal '%d'. got='%d'", len(expected), len(program.Statements))
	}

	for i, response := range expected {
		vs, ok := program.Statements[i].(*ast.VarStatement)
		if !ok {
			t.Fatalf("program.Statement is not equal *ast.VarStatement. got='%T'", program.Statements[i])
		}

		if vs.Name.Name != response {
			t.Fatalf("vs.Name is not '%s'. got='%s'", response, vs.Name.Name)
		}

	}

}

func TestInfixExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"5 + 6 * 9 + 2", "((5 + (6 * 9)) + 2)"},
		{"1 + (2 + 3) + 4", "((1 + (2 + 3)) + 4)"},
		{"(5 + 5) * 2", "((5 + 5) * 2)"},
		{"2 / (5 + 5)", "(2 / (5 + 5))"},
	}

	for _, out := range tests {
		l := lexer.New(out.input)
		p := New(l)
		program := p.ParserProgram()
		if program == nil {
			t.Errorf("ParserProgram is nil.")
		}

		if len(program.Statements) != 1 {
			t.Errorf("len(program.Statements) is not '%d'. got=%d", 1, len(program.Statements))
		}

		if program.Statements[0].String() != out.expected {
			t.Errorf("program.Statements[0].String() is no equal '%s'. got='%s'", out.expected, program.Statements[0].String())
		}
	}
}

func TestSimpleInfixExpresions(t *testing.T) {

	tests := []struct {
		input    string
		left     int64
		operator string
		right    int64
	}{
		{"8 + 6", 8, "+", 6},
		{"5 + 3", 5, "+", 3},
		{"8 - 9", 8, "-", 9},
		{"723 * 123", 723, "*", 123},
		{"1455 / 565", 1455, "/", 565},
	}

	for _, out := range tests {
		l := lexer.New(out.input)
		p := New(l)
		program := p.ParserProgram()
		if program == nil {
			t.Fatal("ParserProgram is nil.")
		}

		if len(program.Statements) != 1 {
			t.Fatalf("len(program.Statement) is not equal '%d'. got=%d", 1, len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not *ast.ExpressionStatement. got=%T", program.Statements[0])
		}

		infix, ok := stmt.Expression.(*ast.InfixExpression)
		if !ok {
			t.Fatalf("stmt.Expression is not *ast.InfixExpression. got=%T", stmt.Expression)
		}

		intLeft, ok := infix.Left.(*ast.Integer)
		if !ok {
			t.Fatalf("infix.Right is not *ast.Integer. got=%T", infix.Left)
		}

		if intLeft.Value != out.left {
			t.Fatalf("intRight is not '%d'. got=%d", out.left, intLeft.Value)
		}

		if infix.Operator != out.operator {
			t.Fatalf("infix.Operator is not '%s'. got=%s", out.operator, infix.Operator)
		}

		intRight, ok := infix.Right.(*ast.Integer)
		if !ok {
			t.Fatalf("infix.Right is not *ast.Integer. got=%T", infix.Right)
		}

		if intRight.Value != out.right {
			t.Fatalf("intRight is not '%d'. got=%d", out.right, intRight.Value)
		}
	}
}

func testInfixExpression(t *testing.T, expr ast.Expression, left interface{}, operator string, right interface{}) bool {

	opExpr, ok := expr.(*ast.InfixExpression)
	if !ok {
		t.Errorf("expr is not *ast.InfixExpression. got=%T(%s)", expr, expr)
		return false
	}

	if !testLiteralExpression(t, opExpr.Left, left) {
		return false
	}

	if opExpr.Operator != operator {
		t.Errorf("expr.Operator is not %s. got=%q", operator, opExpr.Operator)
		return false
	}

	if !testLiteralExpression(t, opExpr.Right, right) {
		return false
	}

	return true
}

func testLiteralExpression(t *testing.T, expr ast.Expression, expected interface{}) bool {

	switch v := expected.(type) {
	case int:
		return testIntegerLiteral(t, expr, int64(v))
	case int64:
		return testIntegerLiteral(t, expr, v)
	case string:
		return testIdentifier(t, expr, v)
	case bool:
		return testBooleanLiteral(t, expr, v)
	}

	t.Errorf("type of expr not handled. got=%T", expr)
	return false
}
func testBooleanLiteral(t *testing.T, expr ast.Expression, value bool) bool {
	bo, ok := expr.(*ast.Boolean)
	if !ok {
		t.Errorf("expr not *ast.Boolean. got=%T", expr)
		return false
	}

	if bo.Value != value {
		t.Errorf("bo.Value not %t. got=%t", value, bo.Value)
		return false
	}

	if bo.TokenLiteral() != fmt.Sprintf("%t", value) {
		t.Errorf("bo.TokenLiteral not %t. got=%s", value, bo.TokenLiteral())
		return false
	}

	return true
}

func testIdentifier(t *testing.T, expr ast.Expression, value string) bool {
	ident, ok := expr.(*ast.Identifier)
	if !ok {
		t.Errorf("expr no *ast.Identifier. got=%T", expr)
		return false
	}

	if ident.Name != value {
		t.Errorf("iden.Value is not %s. got=%s", value, ident.Name)
	}

	if ident.TokenLiteral() != value {
		t.Errorf("ident.TokenLiteral not %s. got=%s", value, ident.TokenLiteral())
		return false
	}

	return true
}

func testIntegerLiteral(t *testing.T, i ast.Expression, value int64) bool {
	interg, ok := i.(*ast.Integer)
	if !ok {
		t.Errorf("i not *ast.IntegerLiteral. got=%T", i)
		return false
	}

	if interg.Value != value {
		t.Errorf("i.Value is not %d. got=%d", value, interg.Value)
		return false
	}

	if interg.TokenLiteral() != fmt.Sprintf("%d", value) {
		t.Errorf("i.TokenLiteral is not %d. got=%s", value, interg.TokenLiteral())
		return false
	}

	return true
}

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.errors
	if len(errors) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(errors))
	for _, msg := range errors {
		t.Errorf("parser error: %q", msg)
	}
	t.FailNow()
}
