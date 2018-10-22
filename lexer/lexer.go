package lexer

import (
	"io/ioutil"

	"github.com/kenshindeveloper/april/token"
)

var (
	currentPosition int
	NUMBER_LINE     = 1
)

type fileInput struct {
	input    string
	position int
	char     byte
	prev     *fileInput
}

type Lexer struct {
	top    *fileInput
	length int
}

func New(input string) *Lexer {
	fi := &fileInput{input: input, prev: nil}
	l := &Lexer{top: fi, length: 1}
	l.readToken()
	return l
}

//ReadFile push input into stack
func (l *Lexer) ReadFile(path string) bool {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return false
	}

	if len(path) <= len(".april") || path[len(path)-6:len(path)] != ".april" {
		return false
	}

	fi := &fileInput{input: string(data), prev: l.top}
	l.top = fi
	l.readToken()
	return true
}

func (l *Lexer) readString() string {
	str := ""
	for l.top.char != '"' {
		str += string(l.top.char)
		l.readToken()
	}
	l.readToken()
	return str
}

func (l *Lexer) readToken() {
	if l.top.position >= len(l.top.input) {
		if !l.expectedFilePrev() {
			l.top.char = 0
		}
	} else {
		l.top.char = l.top.input[l.top.position]
		l.top.position++
	}
}

func (l *Lexer) expectedFilePrev() bool {
	if l.top.prev == nil {
		return false
	}
	fi := l.top
	l.top = fi.prev
	return true
}

func (l *Lexer) comments() {
	if l.top.char == '/' && l.top.char == l.top.input[l.top.position] {
		for l.top.char != '\n' && l.top.char != 0 && l.top.char != '\r' {
			l.readToken()
		}
	}
}

func (l *Lexer) NextToken() token.Token {
	l.comments()
	l.skypeSpace()

	switch l.top.char {
	case ']':
		l.readToken()
		return token.Token{Type: token.RBRACKET, Literal: "]"}
	case '[':
		l.readToken()
		return token.Token{Type: token.LBRACKET, Literal: "["}
	case ',':
		l.readToken()
		return token.Token{Type: token.COMMA, Literal: ","}
	case '"':
		l.readToken()
		tok := token.Token{Type: token.STRING}
		tok.Literal = l.readString()

		return tok
	case ':':
		if l.top.input[l.top.position] == '=' {
			l.readToken()
			l.readToken()
			return token.Token{Type: token.DECLARATION, Literal: ":="}
		}
		l.readToken()
		return token.Token{Type: token.COLON, Literal: ":"}
	case '<':
		if l.top.input[l.top.position] == '=' {
			l.readToken()
			l.readToken()
			return token.Token{Type: token.COMLE, Literal: "<="}
		}
		l.readToken()
		return token.Token{Type: token.COMLT, Literal: "<"}

	case '>':
		if l.top.input[l.top.position] == '=' {
			l.readToken()
			l.readToken()
			return token.Token{Type: token.COMGE, Literal: ">="}
		}
		l.readToken()
		return token.Token{Type: token.COMGT, Literal: ">"}

	case '!':
		if l.top.input[l.top.position] == '=' {
			l.readToken()
			l.readToken()
			return token.Token{Type: token.COMNE, Literal: "!="}
		}
		l.readToken()
		return token.Token{Type: token.BANG, Literal: "!"}

	case '=':
		if l.top.input[l.top.position] == '=' {
			l.readToken()
			l.readToken()
			return token.Token{Type: token.COMEQ, Literal: "=="}
		}
		l.readToken()
		return token.Token{Type: token.EQUAL, Literal: "="}
	case ';':
		l.readToken()
		return token.Token{Type: token.SEMICOLON, Literal: ";"}
	case '%':
		if l.top.input[l.top.position] == '=' {
			l.readToken()
			l.readToken()
			return token.Token{Type: token.ASIGMOD, Literal: "%="}
		}
		l.readToken()
		return token.Token{Type: token.MOD, Literal: "%"}
	case '+':
		if l.top.input[l.top.position] == '+' {
			l.readToken()
			l.readToken()
			return token.Token{Type: token.OPEPLUS, Literal: "++"}
		} else if l.top.input[l.top.position] == '=' {
			l.readToken()
			l.readToken()
			return token.Token{Type: token.ASIGPLUS, Literal: "+="}
		}
		l.readToken()
		return token.Token{Type: token.PLUS, Literal: "+"}
	case '-':
		if l.top.input[l.top.position] == '-' {
			l.readToken()
			l.readToken()
			return token.Token{Type: token.OPEMIN, Literal: "--"}
		} else if l.top.input[l.top.position] == '=' {
			l.readToken()
			l.readToken()
			return token.Token{Type: token.ASIGMIN, Literal: "-="}
		}
		l.readToken()
		return token.Token{Type: token.MIN, Literal: "-"}
	case '*':
		if l.top.input[l.top.position] == '=' {
			l.readToken()
			l.readToken()
			return token.Token{Type: token.ASIGMUL, Literal: "*="}
		}
		l.readToken()
		return token.Token{Type: token.MUL, Literal: "*"}
	case '/':
		if (l.top.position) < len(l.top.input) && l.top.input[l.top.position] == '/' {
			l.comments()
			return token.Token{Type: token.COMMENT, Literal: "//"}
		} else if l.top.input[l.top.position] == '=' {
			l.readToken()
			l.readToken()
			return token.Token{Type: token.ASIGDIV, Literal: "/="}
		}
		l.readToken()
		return token.Token{Type: token.DIV, Literal: "/"}
	case '(':
		l.readToken()
		return token.Token{Type: token.LPAREN, Literal: "("}
	case ')':
		l.readToken()
		return token.Token{Type: token.RPAREN, Literal: ")"}
	case '{':
		l.readToken()
		return token.Token{Type: token.LBRACE, Literal: "{"}
	case '}':
		l.readToken()
		return token.Token{Type: token.RBRACE, Literal: "}"}
	case 0:
		return token.Token{Type: token.EOF, Literal: "EOF"}
	default:
		if l.isChar() {
			return l.tokenEvaluation()
		}
		char := string(l.top.char)
		l.readToken()
		return token.Token{Type: token.ILLEGAL, Literal: string(char)}
	}
}

func (l *Lexer) tokenEvaluation() token.Token {
	ident := ""

	for l.isChar() {
		ident += string(l.top.char)
		l.readToken()
	}

	if isNumeric(ident) {
		return token.Token{Type: token.INT, Literal: ident}
	}

	if isDouble(ident) {
		return token.Token{Type: token.DOUBLE, Literal: ident}
	}

	return token.Token{Type: token.LookKeyword(ident), Literal: ident}
}

func isNumeric(value string) bool {
	for _, data := range value {
		if data < '0' || data > '9' {
			return false
		}
	}
	return true
}

func isDouble(value string) bool {
	for _, data := range value {
		if !(data >= '0' && data <= '9' || data == '.') {
			return false
		}
	}
	return true
}

func (l *Lexer) Save() {
	currentPosition = l.top.position
}

func (l *Lexer) Continue() {
	l.top.position = currentPosition
	l.top.char = l.top.input[l.top.position-1]
}

func (l *Lexer) skypeSpace() {
	for l.top.char == ' ' || l.top.char == '\n' || l.top.char == '\t' || l.top.char == '\a' || l.top.char == '\r' {
		if l.top.char == '\n' {
			NUMBER_LINE += 1
		}

		l.readToken()
	}
}

func (l *Lexer) isChar() bool {
	if l.top.char >= 'a' && l.top.char <= 'z' || l.top.char >= 'A' && l.top.char <= 'Z' || l.top.char == '_' || l.isDigit() || l.top.char == '.' {
		return true
	}
	return false
}

func (l *Lexer) isDigit() bool {
	if l.top.char >= '0' && l.top.char <= '9' {
		return true
	}
	return false
}
