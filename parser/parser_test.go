package parser

import (
	"testing"

	"github.com/kenshindeveloper/april/ast"
	"github.com/kenshindeveloper/april/lexer"
)

func TestLetStatements(t *testing.T) {
	input := `
		var x = 5;
		var y = 10;
		var foobar = 838383;
	`
	l := lexer.New(input)
	p := New(l)

	program := p.ParserProgram()
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

		if !testLetStatement(t, stmt, tt.expectedIdentifier) {
			return
		}
	}
}

func testLetStatement(t *testing.T, stmt ast.Statement, name string) bool {

	if stmt.TokenLiteral() != "var" {
		t.Errorf("stmt.TokenLiteral not 'var'. got=%q", stmt.TokenLiteral())
		return false
	}

	letStmt, ok := stmt.(*ast.LetStatement) //no se que significa esta instruccion
	if !ok {
		t.Errorf("stmt not *ast.LetStatement. got=%T", stmt)
		return false
	}

	if letStmt.Name.Value != name {
		t.Errorf("letStmt.Name.Value not '%s'. got=%s", name, letStmt.Name.Value)
		return false
	}

	if letStmt.Name.TokenLiteral() != name {
		t.Errorf("letStmt.Name not '%s'. got=%s", name, letStmt.Name)
		return false
	}

	return true
}
