package lexer

import (
	"log"
	"testing"

	"github.com/kenshindeveloper/april/token"
)

func TestFunctionClosureExpression(t *testing.T) {
	input := `fn(x, y, z) {5+4;}`

	tests := []struct {
		tokenExpected token.TokenType
		expected      string
	}{
		{token.FNCLOSURE, "fn"},
		{token.LPAREN, "("},
		{token.IDENT, "x"},
		{token.COMMA, ","},
		{token.IDENT, "y"},
		{token.COMMA, ","},
		{token.IDENT, "z"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.INT, "5"},
		{token.PLUS, "+"},
		{token.INT, "4"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
	}

	l := New(input)

	for _, data := range tests {
		out := l.NextToken()

		if data.tokenExpected != out.Type {
			t.Fatalf("tokenExpected is no equal '%s'. got='%s'", out.Type, data.tokenExpected)
		}

		if data.expected != out.Literal {
			t.Fatalf("expected is no equal '%s'. got='%s'", out.Literal, data.expected)
		}

	}

}

func TestDoubleToken(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{"9.6", 9.6},
		{"945.68686", 945.68686},
		{"0.68686", 0.68686},
	}

	for _, data := range tests {
		l := New(data.input)

		tok := l.NextToken()
		if tok.Type != token.DOUBLE {
			t.Fatalf("tok.Type is not equal '%s'. got='%s'", token.DOUBLE, tok.Type)
		}
	}
}

func TestStringToken(t *testing.T) {
	input := `"hola"`

	l := New(input)

	tok := l.NextToken()

	if tok.Type != token.STRING {
		t.Fatalf("tok.Type is not 'STRING'. got='%T'", tok.Type)
	}

	if tok.Literal != "hola" {
		t.Fatalf("tok.Literal is not equal 'hola'. got='%s'", tok.Literal)
	}
}

func TestLexer(t *testing.T) {
	input := `
		// esto es una prueba
		var x:int = 5;
		// esto es otra prueba
		y := 6;
		[5];
		{"foo": "bar"};
		x++;
		y--;
		x += 1;
		y -= 1;
		z *= 1;
		w /= 1;
	`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.COMMENT, "//"},
		{token.VAR, "var"},
		{token.IDENT, "x"},
		{token.COLON, ":"},
		{token.OPEINT, "int"},
		{token.EQUAL, "="},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		{token.COMMENT, "//"},
		{token.IDENT, "y"},
		{token.DECLARATION, ":="},
		{token.INT, "6"},
		{token.SEMICOLON, ";"},
		{token.LBRACKET, "["},
		{token.INT, "5"},
		{token.RBRACKET, "]"},
		{token.SEMICOLON, ";"},
		{token.LBRACE, "{"},
		{token.STRING, "foo"},
		{token.COLON, ":"},
		{token.STRING, "bar"},
		{token.RBRACE, "}"},
		{token.SEMICOLON, ";"},
		{token.IDENT, "x"},
		{token.OPEPLUS, "++"},
		{token.SEMICOLON, ";"},
		{token.IDENT, "y"},
		{token.OPEMIN, "--"},
		{token.SEMICOLON, ";"},
		{token.IDENT, "x"},
		{token.ASIGPLUS, "+="},
		{token.INT, "1"},
		{token.SEMICOLON, ";"},
		{token.IDENT, "y"},
		{token.ASIGMIN, "-="},
		{token.INT, "1"},
		{token.SEMICOLON, ";"},
		{token.IDENT, "z"},
		{token.ASIGMUL, "*="},
		{token.INT, "1"},
		{token.SEMICOLON, ";"},
		{token.IDENT, "w"},
		{token.ASIGDIV, "/="},
		{token.INT, "1"},
		{token.SEMICOLON, ";"},
	}

	l := New(input)
	for _, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			log.Fatalf("el tipo del token experado es '%s', tiene='%s'", tt.expectedType, tok.Literal)
		}

		if tok.Literal != tt.expectedLiteral {
			log.Fatalf("el literal experado es '%s'. tiene='%s'", tt.expectedLiteral, tok.Literal)
		}
	}

}
