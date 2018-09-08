package ast

import (
	"testing"

	"github.com/kenshindeveloper/april/token"
)

func TestString(t *testing.T) {
	program := &Program{
		Statements: []Statement{
			&VarStatement{
				Token: token.Token{Type: token.VAR, Literal: "var"},
				Name: &Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "myvar"},
					Value: "myvar",
				},
				Value: &Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "anotherVar"},
					Value: "anotherVar",
				},
			},
		},
	}

	if program.String() != "var myvar = anotherVar;" {
		t.Errorf("program.String() wrong. got=%q", program.String())
	}
}
