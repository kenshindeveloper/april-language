package token

//***************************************************************************************
//***************************************************************************************
//***************************************************************************************
//Constantes que representan los lexemas dentro del lenguaje April.
const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"
	IDENT   = "IDENT"
	RETURN  = "RETURN"
	IMPORT  = "IMPORT"

	OPEINT    = "OPEINT"
	OPEBOOL   = "OPEBOOL"
	OPESTR    = "OPESTR"
	OPEDOUBLE = "OPEDOUBLE"
	LIST      = "LIST"
	MAP       = "MAP"
	FUNC      = "FUNC"
	STREAM    = "STREAM"
	NIL       = "NIL"

	INT    = "INT"
	DOUBLE = "DOUBLE"
	STRING = "STRING"
	TRUE   = "TRUE"
	FALSE  = "FALSE"

	COMMENT = "//"
	JUMP    = "\n"

	NOT         = "NOT"
	EQUAL       = "="
	SEMICOLON   = ";"
	COLON       = ":"
	DECLARATION = ":="
	DOT         = "."
	COMMA       = ","
	BANG        = "!"

	AND = "AND"
	OR  = "OR"

	COMNE = "!="
	COMEQ = "=="
	COMLE = "<="
	COMGE = ">="
	COMLT = "<"
	COMGT = ">"

	OPEPLUS  = "++"
	OPEMIN   = "--"
	ASIGPLUS = "+="
	ASIGMIN  = "-="
	ASIGMUL  = "*="
	ASIGDIV  = "/="
	ASIGMOD  = "%="

	LPAREN   = "("
	RPAREN   = ")"
	LBRACE   = "{"
	RBRACE   = "}"
	LBRACKET = "["
	RBRACKET = "]"

	PLUS = "+"
	MIN  = "-"
	MUL  = "*"
	DIV  = "/"
	MOD  = "%"

	VAR       = "VAR"
	GLOBAL    = "GLOBAL"
	IF        = "IF"
	ELSE      = "ELSE"
	FNCLOSURE = "FNCLOSURE"
	FOR       = "FOR"
	BREAK     = "BREAK"
)

//TokenType es un tipo de dato que representa el 'tipo' de los lexemas definidos.
type TokenType string

//Token es un tipo de dato que representa un lexema.
type Token struct {
	Type    TokenType
	Literal string
}

var keywords = map[string]TokenType{
	"var":    VAR,
	"global": GLOBAL,
	"true":   TRUE,
	"false":  FALSE,
	"not":    NOT,
	"int":    OPEINT,
	"bool":   OPEBOOL,
	"string": OPESTR,
	"double": OPEDOUBLE,
	"func":   FUNC,
	"stream": STREAM,
	"and":    AND,
	"or":     OR,
	"return": RETURN,
	"if":     IF,
	"else":   ELSE,
	"fn":     FNCLOSURE,
	"for":    FOR,
	"list":   LIST,
	"map":    MAP,
	"import": IMPORT,
	"break":  BREAK,
	"nil":    NIL,
}

//LookKeyword comprueba si la variable 'name' es una palabra clave dentro del map de keywords.
func LookKeyword(name string) TokenType {
	if tok, ok := keywords[name]; ok {
		return tok
	}
	return IDENT
}
