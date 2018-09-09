package parser

import (
	"testing"

	"github.com/kenshindeveloper/april/ast"
	"github.com/kenshindeveloper/april/lexer"
)

func TestIntegerLiteralExpression(t *testing.T) {
	input := `5;`

	l := lexer.New(input)
	p := New(l)
	program := p.ParserProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statemets is not enough statements. got=%d\n", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	literal, ok := stmt.Expression.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.IntegerLiteral. got=%T", stmt.Expression)
	}

	if literal.Value != 5 {
		t.Fatalf("literal is not %d. got=%q", 5, literal.Value)
	}

	if literal.TokenLiteral() != "5" {
		t.Fatalf("literal is not %s. got=%s", "5", literal.TokenLiteral())
	}

}

func TestIdentifierExpression(t *testing.T) {
	input := `foobar;`

	l := lexer.New(input)
	p := New(l)
	program := p.ParserProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d\n", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statemets[0] is not *ast.ExpressionStatement. got=%T\n", program.Statements[0])
	}

	ident, ok := stmt.Expression.(*ast.Identifier)
	if !ok {
		t.Fatalf("stmt.Expression is not *ast.Identifier. got=%T\n", stmt.Expression)
	}

	if ident.Value != "foobar" {
		t.Fatalf("ident is not %s. got=%q", "foobar", ident.Value)
	}

	if ident.TokenLiteral() != "foobar" {
		t.Fatalf("ident.TokenLiteral is not %s. got=%s", "foobar", ident.TokenLiteral())
	}
}

func TestReturnStatements(t *testing.T) {
	input := `
		return 5;
		return 10;
		return 993322;
	`
	l := lexer.New(input) //Modulo lexico
	p := New(l)           //Modulo Sintactico

	program := p.ParserProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not containt 3 statements. got=%d", len(program.Statements))
	}

	for _, stmt := range program.Statements {
		returnStmt, ok := stmt.(*ast.ReturnStatement)
		if !ok {
			t.Fatalf("stmt not *ast.returnStatement. got=%T", stmt)
			continue
		}
		if returnStmt.TokenLiteral() != "return" {
			t.Errorf("returnStmt.TokenLiteral not 'return', got %q", returnStmt.TokenLiteral())
		}
	}
}

func TestVarStatements(t *testing.T) {
	input := `
		var x int = 5;
		var y int = 10;
		var foobar int = 838383;
	`
	l := lexer.New(input)
	p := New(l)

	program := p.ParserProgram()
	checkParserErrors(t, p)
	if program == nil {
		t.Fatalf("ParserProgram() return nil")
	}

	if len(program.Statements) != 3 {
		t.Fatalf("program.Statement does not containt 3 statement. got=%d", len(program.Statements))
	}

	tests := []struct {
		expectedIdentifier string
	}{
		{"x"},
		{"y"},
		{"foobar"},
	}

	for i, tt := range tests {

		stmt := program.Statements[i]

		if !testVarStatement(t, stmt, tt.expectedIdentifier) {
			return
		}
	}
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

func testVarStatement(t *testing.T, stmt ast.Statement, name string) bool {

	if stmt.TokenLiteral() != "var" {
		t.Errorf("stmt.TokenLiteral not 'var'. got=%q", stmt.TokenLiteral())
		return false
	}

	varStmt, ok := stmt.(*ast.VarStatement) //no se que significa esta instruccion
	if !ok {
		t.Errorf("stmt not *ast.LetStatement. got=%T", stmt)
		return false
	}

	if varStmt.Name.Value != name {
		t.Errorf("varStmt.Name.Value not '%s'. got=%s", name, varStmt.Name.Value)
		return false
	}

	if varStmt.Name.TokenLiteral() != name {
		t.Errorf("varStmt.Name not '%s'. got=%s", name, varStmt.Name)
		return false
	}

	return true
}
