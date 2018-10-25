package parser

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/kenshindeveloper/april/ast"
	"github.com/kenshindeveloper/april/lexer"
	"github.com/kenshindeveloper/april/token"
)

const (
	_ = iota
	LESSVALUE
	EQUALS
	LESSGREATER
	SUM
	PRODUCT
	PREFIX
	CALL
	INDEX
)

var scope = map[string]int{
	"fn":  0,
	"if":  1,
	"for": 2,
}

var precedens = map[token.TokenType]int{
	token.PLUS:        SUM,
	token.MIN:         SUM,
	token.MUL:         PRODUCT,
	token.DIV:         PRODUCT,
	token.MOD:         PRODUCT,
	token.COMNE:       EQUALS,
	token.COMEQ:       EQUALS,
	token.COMLE:       LESSGREATER,
	token.COMGE:       LESSGREATER,
	token.COMLT:       LESSGREATER,
	token.COMGT:       LESSGREATER,
	token.AND:         LESSGREATER,
	token.OR:          LESSGREATER,
	token.EQUAL:       PREFIX,
	token.DECLARATION: PREFIX, //ojo con esta precedencia
	token.ASIGPLUS:    PREFIX,
	token.ASIGMIN:     PREFIX,
	token.ASIGMUL:     PREFIX,
	token.ASIGDIV:     PREFIX,
	token.ASIGMOD:     PREFIX,
	token.LPAREN:      CALL,
	token.LBRACKET:    INDEX,
	token.OPEPLUS:     INDEX,
	token.OPEMIN:      INDEX,
}

type (
	prefixFn  = func() ast.Expression
	infixFn   = func(ast.Expression) ast.Expression
	postfixFn = func() ast.Expression
)

type Parser struct {
	lexer     *lexer.Lexer
	curToken  token.Token
	peekToken token.Token

	prefixFns  map[token.TokenType]prefixFn
	infixFns   map[token.TokenType]infixFn
	postfixFns map[token.TokenType]postfixFn

	errors []string
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{lexer: l}
	p.nextToken()
	p.nextToken()

	//registro de funciones postfijas
	p.postfixFns = make(map[token.TokenType]postfixFn)
	p.registerPostfix(token.OPEPLUS, p.parsePostfixExpression)
	p.registerPostfix(token.OPEMIN, p.parsePostfixExpression)

	//registro de funciones prefijas
	p.prefixFns = make(map[token.TokenType]prefixFn)
	p.registerPrefix(token.INT, p.parseIntegerExpression)
	p.registerPrefix(token.IDENT, p.parseIdentifierExpression)
	p.registerPrefix(token.NIL, p.parseNilExpression)
	p.registerPrefix(token.LPAREN, p.parseParensExpression)
	p.registerPrefix(token.TRUE, p.parseBooleanExpression)
	p.registerPrefix(token.FALSE, p.parseBooleanExpression)
	p.registerPrefix(token.NOT, p.parsePrefixExpression)
	p.registerPrefix(token.MIN, p.parsePrefixExpression)
	p.registerPrefix(token.STRING, p.parseStringExpression)
	p.registerPrefix(token.DOUBLE, p.parseDoubleExpression)
	p.registerPrefix(token.IF, p.parseIfExpression)
	p.registerPrefix(token.FNCLOSURE, p.parseFnClosure)
	p.registerPrefix(token.LBRACKET, p.parseListExpression)
	p.registerPrefix(token.LBRACE, p.parseHashExpression)

	//registro de funciones infijas
	p.infixFns = make(map[token.TokenType]infixFn)
	p.registerInfix(token.PLUS, p.parseInfixExpression)
	p.registerInfix(token.MIN, p.parseInfixExpression)
	p.registerInfix(token.MUL, p.parseInfixExpression)
	p.registerInfix(token.DIV, p.parseInfixExpression)
	p.registerInfix(token.MOD, p.parseInfixExpression)
	p.registerInfix(token.COMNE, p.parseInfixExpression)
	p.registerInfix(token.COMEQ, p.parseInfixExpression)
	p.registerInfix(token.COMLE, p.parseInfixExpression)
	p.registerInfix(token.COMGE, p.parseInfixExpression)
	p.registerInfix(token.COMLT, p.parseInfixExpression)
	p.registerInfix(token.COMGT, p.parseInfixExpression)
	p.registerInfix(token.EQUAL, p.parseAssignExpression)
	p.registerInfix(token.ASIGPLUS, p.parseAssignOpeExpression)
	p.registerInfix(token.ASIGMIN, p.parseAssignOpeExpression)
	p.registerInfix(token.ASIGMUL, p.parseAssignOpeExpression)
	p.registerInfix(token.ASIGDIV, p.parseAssignOpeExpression)
	p.registerInfix(token.ASIGMOD, p.parseAssignOpeExpression)
	p.registerInfix(token.AND, p.parseInfixExpression)
	p.registerInfix(token.OR, p.parseInfixExpression)
	p.registerInfix(token.DECLARATION, p.parseImplicitExpression)
	p.registerInfix(token.LPAREN, p.parseCallExpression)
	p.registerInfix(token.LBRACKET, p.parseIndexExpression)

	return p
}

func (p *Parser) ParserProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for !p.curTokenIs(token.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}

	return program
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.VAR:
		return p.parseVarStatement()
	case token.GLOBAL:
		return p.parseGlobalStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	case token.BREAK:
		return p.parseBreakStatement()
	case token.FOR:
		return p.parseForStatement()
	case token.IMPORT:
		return p.parseImportStatement()
	case token.FNCLOSURE:
		return p.parseFunctionStatement()
	case token.COMMENT:
		return nil
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseFunctionStatement() ast.Statement {
	fn := &ast.Function{Token: p.curToken, Line: lexer.NUMBER_LINE}

	if !p.expectedTokenPeek(token.IDENT) {
		return nil
	}
	fn.Name = &ast.Identifier{Token: p.curToken, Name: p.curToken.Literal, Line: lexer.NUMBER_LINE}

	if !p.expectedTokenPeek(token.LPAREN) {
		return nil
	}

	fn.Parameters = p.parseFunctionParametersExpression()
	if !p.curTokenIs(token.RPAREN) {
		return nil
	}

	if p.isPeekBasicType() {
		p.nextToken()
		fn.Type = &ast.Identifier{Token: p.curToken, Name: p.curToken.Literal, Line: lexer.NUMBER_LINE}
		p.nextToken()

	} else if p.peekTokenIs(token.IDENT) {
		msg := fmt.Sprintf("Line: %d - function expression is incorrect.", lexer.NUMBER_LINE)
		p.errors = append(p.errors, msg)
		return nil
	} else if p.peekTokenIs(token.LBRACE) {
		fn.Type = nil
		p.nextToken()
	}

	//ultima validacion para estar seguro del token actual.
	if !p.curTokenIs(token.LBRACE) {
		msg := fmt.Sprintf("Line: %d - function expression is incorrect. : %s", lexer.NUMBER_LINE, p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}

	fn.Body = p.parseBlockStatement()
	if fn.Body == nil {
		return nil
	}
	fn.Body.Type = scope["fn"]
	return fn
}

func (p *Parser) parseFunctionParametersExpression() []*ast.FunctionParameters {
	paramenters := []*ast.FunctionParameters{}

	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		return paramenters
	}

	p.nextToken()
	param := &ast.FunctionParameters{Token: p.curToken, Line: lexer.NUMBER_LINE}
	param.Name = &ast.Identifier{Token: p.curToken, Name: p.curToken.Literal, Line: lexer.NUMBER_LINE}

	if !p.expectedTokenPeek(token.COLON) {
		return nil
	}
	if !p.isPeekBasicType() {
		msg := fmt.Sprintf("Line: %d - function paramenters incorrect. : %s", lexer.NUMBER_LINE, p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}
	p.nextToken()
	param.Type = &ast.Identifier{Token: p.curToken, Name: p.curToken.Literal, Line: lexer.NUMBER_LINE}
	paramenters = append(paramenters, param)

	for p.peekTokenIs(token.COMMA) {
		p.nextToken() //coma
		p.nextToken() //valor
		param := &ast.FunctionParameters{Token: p.curToken, Line: lexer.NUMBER_LINE}
		param.Name = &ast.Identifier{Token: p.curToken, Name: p.curToken.Literal, Line: lexer.NUMBER_LINE}

		if !p.expectedTokenPeek(token.COLON) {
			return nil
		}
		if !p.isPeekBasicType() {
			msg := fmt.Sprintf("Line: %d - function paramenters incorrect. : %s", lexer.NUMBER_LINE, p.curToken.Literal)
			p.errors = append(p.errors, msg)
			return nil
		}
		p.nextToken()
		param.Type = &ast.Identifier{Token: p.curToken, Name: p.curToken.Literal, Line: lexer.NUMBER_LINE}
		paramenters = append(paramenters, param)
	}
	if !p.peekTokenIs(token.RPAREN) {
		msg := fmt.Sprintf("Line: %d - function paramenters incorrect.", lexer.NUMBER_LINE)
		p.errors = append(p.errors, msg)
		return nil
	}
	p.nextToken()

	return paramenters
}

func (p *Parser) parseImportStatement() ast.Statement {
	if !p.peekTokenIs(token.STRING) {
		msg := fmt.Sprintf("Line: %d - import expression wrong.", lexer.NUMBER_LINE)
		p.errors = append(p.errors, msg)
		return nil
	}

	if !p.lexer.ReadFile(p.peekToken.Literal) {
		msg := fmt.Sprintf("Line: %d - incorrect path file: '%s'.", lexer.NUMBER_LINE, p.peekToken.Literal)
		p.errors = append(p.errors, msg)
	}
	return nil
}

func (p *Parser) parseForStatement() *ast.ForStatement {
	fs := &ast.ForStatement{Token: p.curToken, Line: lexer.NUMBER_LINE}
	flag := false

	if p.peekTokenIs(token.LPAREN) {
		flag = true
		p.nextToken()
	}
	p.nextToken()

	//-------------------------------------------------

	p.lexer.Save()
	auxCurToken := p.curToken
	auxPeekToken := p.peekToken
	countSemicolon := 0
	for !p.curTokenIs(token.LBRACE) && !p.curTokenIs(token.EOF) {
		if p.curTokenIs(token.SEMICOLON) {
			countSemicolon++
		}
		p.nextToken()
	}

	p.lexer.Continue()
	p.curToken = auxCurToken
	p.peekToken = auxPeekToken
	//-------------------------------------------------
	switch countSemicolon {
	case 0:
		fs.Condition = p.parseExpression(LESSVALUE) //aqui
		_, ok := fs.Condition.(*ast.ImplicitDeclarationExpression)
		if !ok && flag && !p.expectedTokenPeek(token.RPAREN) {
			msg := fmt.Sprintf("Line: %d - for expression is incorrect, expected token ')'", lexer.NUMBER_LINE)
			p.errors = append(p.errors, msg)
			return nil
		}

		if !ok && !p.peekTokenIs(token.LBRACE) {
			msg := fmt.Sprintf("Line: %d - for expression is incorrect, expected token '{'", lexer.NUMBER_LINE)
			p.errors = append(p.errors, msg)
			return nil
		}

		if !ok || flag {
			p.nextToken()
		}
		fs.Body = p.parseBlockStatement()
		if fs.Body == nil {
			return nil
		}
		fs.Body.Type = scope["for"]
		return fs
	case 2:
		fs.Declaration = p.parseExpression(LESSVALUE) //aqui

		p.nextToken()
		fs.Condition = p.parseExpression(LESSVALUE) //aqui
		if !p.expectedTokenPeek(token.SEMICOLON) {
			msg := fmt.Sprintf("Line: %d - for expression is incorrect, expected token ';'", lexer.NUMBER_LINE)
			p.errors = append(p.errors, msg)
			return nil
		}
		p.nextToken()
		fs.Operation = p.parseExpression(LESSVALUE)
		if flag && !p.expectedTokenPeek(token.LBRACE) {
			msg := fmt.Sprintf("Line: %d - for expression is incorrect, expected token '{'", lexer.NUMBER_LINE)
			p.errors = append(p.errors, msg)
			return nil
		}

		fs.Body = p.parseBlockStatement()
		if fs.Body == nil {
			return nil
		}
		fs.Body.Type = scope["for"]
		return fs
	default:
		msg := fmt.Sprintf("Line: %d - count ';' incorrect.", lexer.NUMBER_LINE)
		p.errors = append(p.errors, msg)
		return nil
	}
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	rs := &ast.ReturnStatement{Token: p.curToken, Line: lexer.NUMBER_LINE}

	p.nextToken()
	rs.Expression = p.parseExpression(LESSVALUE)

	if !p.curTokenIs(token.SEMICOLON) && !p.peekTokenIs(token.RBRACE) {
		p.nextToken()
	}
	return rs
}

func (p *Parser) parseBreakStatement() *ast.BreakStatement {
	b := &ast.BreakStatement{Token: p.curToken, Line: lexer.NUMBER_LINE}

	p.nextToken()
	if !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return b
}

func (p *Parser) parseVarStatement() *ast.VarStatement {
	vs := &ast.VarStatement{Token: p.curToken, Line: lexer.NUMBER_LINE}

	if !p.expectedTokenPeek(token.IDENT) {
		msg := fmt.Sprintf("Line: %d - var is not declated.", lexer.NUMBER_LINE)
		p.errors = append(p.errors, msg)
		return nil
	}
	vs.Name = &ast.Identifier{Token: p.curToken, Name: p.curToken.Literal, Line: lexer.NUMBER_LINE}

	if !p.expectedTokenPeek(token.COLON) {
		msg := fmt.Sprintf("Line: %d - colon ':' is not declated.", lexer.NUMBER_LINE)
		p.errors = append(p.errors, msg)
		return nil
	}

	if !p.isPeekBasicType() {
		msg := fmt.Sprintf("Line: %d - type '%s' is not declated.", lexer.NUMBER_LINE, p.peekToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}
	p.nextToken()

	typeName := p.curToken.Literal
	vs.Type = &ast.Identifier{Token: p.curToken, Name: p.curToken.Literal, Line: lexer.NUMBER_LINE}

	p.nextToken()
	//!p.curTokenIs(token.EQUAL)
	if p.curTokenIs(token.EOF) || p.curTokenIs(token.SEMICOLON) {
		if typeName == "int" {
			tok := token.Token{Type: token.INT, Literal: "0"}
			vs.Value = &ast.Integer{Token: tok, Value: 0, Line: lexer.NUMBER_LINE}
		} else if typeName == "bool" {
			tok := token.Token{Type: token.FALSE, Literal: "false"}
			vs.Value = &ast.Boolean{Token: tok, Value: false, Line: lexer.NUMBER_LINE}
		} else if typeName == "string" {
			tok := token.Token{Type: token.STRING, Literal: ""}
			vs.Value = &ast.String{Token: tok, Value: "", Line: lexer.NUMBER_LINE}
		} else if typeName == "double" {
			tok := token.Token{Type: token.STRING, Literal: "0.0"}
			vs.Value = &ast.Double{Token: tok, Value: 0.0, Line: lexer.NUMBER_LINE}
		} else if typeName == "list" {
			tok := token.Token{Type: token.STRING, Literal: "[]"}
			vs.Value = &ast.List{Token: tok, Line: lexer.NUMBER_LINE}
		} else if typeName == "map" {
			tok := token.Token{Type: token.STRING, Literal: "{}"}
			vs.Value = &ast.Hash{Token: tok, Line: lexer.NUMBER_LINE}
		} else if typeName == "func" {
			msg := fmt.Sprintf("Line: %d - declaration func error: var %s '%s' must be a expression.", lexer.NUMBER_LINE, vs.Name, strings.ToUpper(typeName))
			p.errors = append(p.errors, msg)
			return nil
		} else if typeName == "stream" {
			msg := fmt.Sprintf("Line: %d - declaration stream error: var %s '%s' must be a expression.", lexer.NUMBER_LINE, vs.Name, strings.ToUpper(typeName))
			p.errors = append(p.errors, msg)
			return nil
		} else {
			msg := fmt.Sprintf("Line: %d - declaration error: var %s '%s' ", lexer.NUMBER_LINE, vs.Name, strings.ToUpper(typeName))
			p.errors = append(p.errors, msg)
			return nil
		}

		return vs
	} else if p.curTokenIs(token.EQUAL) {
		p.nextToken()
		vs.Value = p.parseExpression(LESSVALUE)
		if vs.Value != nil {
			switch vs.Value.(type) {
			case *ast.Hash:
				if typeName != "map" {
					msg := fmt.Sprintf("Line: %d - declaration error: var %s '%s' = MAP", lexer.NUMBER_LINE, vs.Name, strings.ToUpper(typeName))
					p.errors = append(p.errors, msg)
					return nil
				}
			case *ast.List:
				if typeName != "list" {
					msg := fmt.Sprintf("Line: %d - declaration error: var %s '%s' = LIST", lexer.NUMBER_LINE, vs.Name, strings.ToUpper(typeName))
					p.errors = append(p.errors, msg)
					return nil
				}
			case *ast.Integer:
				if typeName == "double" {
					v := vs.Value.(*ast.Integer)
					vs.Value = &ast.Double{Token: v.Token, Value: float64(v.Value), Line: lexer.NUMBER_LINE}
				} else if typeName != "int" {
					msg := fmt.Sprintf("Line: %d - declaration error: var %s '%s' = INTEGER", lexer.NUMBER_LINE, vs.Name, strings.ToUpper(typeName))
					p.errors = append(p.errors, msg)
					return nil
				}
			case *ast.Boolean:
				if typeName != "bool" {
					msg := fmt.Sprintf("Line: %d - declaration error: var %s '%s' = BOOLEAN", lexer.NUMBER_LINE, vs.Name, strings.ToUpper(typeName))
					p.errors = append(p.errors, msg)
					return nil
				}
			case *ast.String:
				if typeName != "string" {
					msg := fmt.Sprintf("Line: %d - declaration error: var %s '%s' = STRING", lexer.NUMBER_LINE, vs.Name, strings.ToUpper(typeName))
					p.errors = append(p.errors, msg)
					return nil
				}
			case *ast.Double:
				if typeName != "double" {
					msg := fmt.Sprintf("Line: %d - declaration error: var %s '%s' = DOUBLE", lexer.NUMBER_LINE, vs.Name, strings.ToUpper(typeName))
					p.errors = append(p.errors, msg)
					return nil
				}
			case *ast.FunctionClosure:
				if typeName != "func" {
					msg := fmt.Sprintf("Line: %d - declaration error: var %s '%s' = FUNC", lexer.NUMBER_LINE, vs.Name, strings.ToUpper(typeName))
					p.errors = append(p.errors, msg)
					return nil
				}
			case *ast.Stream:
				if typeName != "func" {
					msg := fmt.Sprintf("Line: %d - declaration error: var %s '%s' = STREAM", lexer.NUMBER_LINE, vs.Name, strings.ToUpper(typeName))
					p.errors = append(p.errors, msg)
					return nil
				}
			}
		}

		if p.peekTokenIs(token.SEMICOLON) {
			p.nextToken()
		}

		return vs
	} else {
		msg := fmt.Sprintf("Line: %d - token incorrect '%s'.", lexer.NUMBER_LINE, p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}
}

func (p *Parser) parseGlobalStatement() *ast.GlobalStatement {
	gs := &ast.GlobalStatement{Token: p.curToken, Line: lexer.NUMBER_LINE}

	if !p.expectedTokenPeek(token.IDENT) {
		msg := fmt.Sprintf("Line: %d - var is not declated.", lexer.NUMBER_LINE)
		p.errors = append(p.errors, msg)
		return nil
	}
	gs.Name = &ast.Identifier{Token: p.curToken, Name: p.curToken.Literal, Line: lexer.NUMBER_LINE}

	if !p.expectedTokenPeek(token.COLON) {
		msg := fmt.Sprintf("Line: %d - colon ':' is not declated.", lexer.NUMBER_LINE)
		p.errors = append(p.errors, msg)
		return nil
	}

	if !p.isPeekBasicType() {
		msg := fmt.Sprintf("Line: %d - type '%s' is not declated.", lexer.NUMBER_LINE, p.peekToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}
	p.nextToken()

	typeName := p.curToken.Literal
	gs.Type = &ast.Identifier{Token: p.curToken, Name: p.curToken.Literal, Line: lexer.NUMBER_LINE}

	p.nextToken()
	//!p.curTokenIs(token.EQUAL)
	if p.curTokenIs(token.EOF) || p.curTokenIs(token.SEMICOLON) {
		if typeName == "int" {
			tok := token.Token{Type: token.INT, Literal: "0"}
			gs.Value = &ast.Integer{Token: tok, Value: 0, Line: lexer.NUMBER_LINE}
		} else if typeName == "bool" {
			tok := token.Token{Type: token.FALSE, Literal: "false"}
			gs.Value = &ast.Boolean{Token: tok, Value: false, Line: lexer.NUMBER_LINE}
		} else if typeName == "string" {
			tok := token.Token{Type: token.STRING, Literal: ""}
			gs.Value = &ast.String{Token: tok, Value: "", Line: lexer.NUMBER_LINE}
		} else if typeName == "double" {
			tok := token.Token{Type: token.STRING, Literal: "0.0"}
			gs.Value = &ast.Double{Token: tok, Value: 0.0, Line: lexer.NUMBER_LINE}
		} else if typeName == "list" {
			tok := token.Token{Type: token.STRING, Literal: "[]"}
			gs.Value = &ast.List{Token: tok, Line: lexer.NUMBER_LINE}
		} else if typeName == "map" {
			tok := token.Token{Type: token.STRING, Literal: "{}"}
			gs.Value = &ast.Hash{Token: tok, Line: lexer.NUMBER_LINE}
		} else if typeName == "func" {
			msg := fmt.Sprintf("Line: %d - declaration func error: var %s '%s' must be a expression.", lexer.NUMBER_LINE, gs.Name, strings.ToUpper(typeName))
			p.errors = append(p.errors, msg)
			return nil
		} else if typeName == "stream" {
			msg := fmt.Sprintf("Line: %d - declaration stream error: var %s '%s' must be a expression.", lexer.NUMBER_LINE, gs.Name, strings.ToUpper(typeName))
			p.errors = append(p.errors, msg)
			return nil
		} else {
			msg := fmt.Sprintf("Line: %d - declaration error: var %s '%s' ", lexer.NUMBER_LINE, gs.Name, strings.ToUpper(typeName))
			p.errors = append(p.errors, msg)
			return nil
		}

		return gs
	} else if p.curTokenIs(token.EQUAL) {
		p.nextToken()
		gs.Value = p.parseExpression(LESSVALUE)
		if gs.Value != nil {
			switch gs.Value.(type) {
			case *ast.Hash:
				if typeName != "map" {
					msg := fmt.Sprintf("Line: %d - declaration error: var %s '%s' = MAP", lexer.NUMBER_LINE, gs.Name, strings.ToUpper(typeName))
					p.errors = append(p.errors, msg)
					return nil
				}
			case *ast.List:
				if typeName != "list" {
					msg := fmt.Sprintf("Line: %d - declaration error: var %s '%s' = LIST", lexer.NUMBER_LINE, gs.Name, strings.ToUpper(typeName))
					p.errors = append(p.errors, msg)
					return nil
				}
			case *ast.Integer:
				if typeName == "double" {
					v := gs.Value.(*ast.Integer)
					gs.Value = &ast.Double{Token: v.Token, Value: float64(v.Value), Line: lexer.NUMBER_LINE}
				} else if typeName != "int" {
					msg := fmt.Sprintf("Line: %d - declaration error: var %s '%s' = INTEGER", lexer.NUMBER_LINE, gs.Name, strings.ToUpper(typeName))
					p.errors = append(p.errors, msg)
					return nil
				}
			case *ast.Boolean:
				if typeName != "bool" {
					msg := fmt.Sprintf("Line: %d - declaration error: var %s '%s' = BOOLEAN", lexer.NUMBER_LINE, gs.Name, strings.ToUpper(typeName))
					p.errors = append(p.errors, msg)
					return nil
				}
			case *ast.String:
				if typeName != "string" {
					msg := fmt.Sprintf("Line: %d - declaration error: var %s '%s' = STRING", lexer.NUMBER_LINE, gs.Name, strings.ToUpper(typeName))
					p.errors = append(p.errors, msg)
					return nil
				}
			case *ast.Double:
				if typeName != "double" {
					msg := fmt.Sprintf("Line: %d - declaration error: var %s '%s' = DOUBLE", lexer.NUMBER_LINE, gs.Name, strings.ToUpper(typeName))
					p.errors = append(p.errors, msg)
					return nil
				}
			case *ast.FunctionClosure:
				if typeName != "func" {
					msg := fmt.Sprintf("Line: %d - declaration error: var %s '%s' = FUNC", lexer.NUMBER_LINE, gs.Name, strings.ToUpper(typeName))
					p.errors = append(p.errors, msg)
					return nil
				}
			case *ast.Stream:
				if typeName != "func" {
					msg := fmt.Sprintf("Line: %d - declaration error: var %s '%s' = STREAM", lexer.NUMBER_LINE, gs.Name, strings.ToUpper(typeName))
					p.errors = append(p.errors, msg)
					return nil
				}
			}
		}

		if p.peekTokenIs(token.SEMICOLON) {
			p.nextToken()
		}

		return gs
	} else {
		msg := fmt.Sprintf("Line: %d - token incorrect '%s'.", lexer.NUMBER_LINE, p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}
}

func (p *Parser) isPeekBasicType() bool {
	switch {
	case p.peekTokenIs(token.OPEINT):
		return true
	case p.peekTokenIs(token.OPEBOOL):
		return true
	case p.peekTokenIs(token.OPEDOUBLE):
		return true
	case p.peekTokenIs(token.OPESTR):
		return true
	case p.peekTokenIs(token.LIST):
		return true
	case p.peekTokenIs(token.MAP):
		return true
	case p.peekTokenIs(token.FUNC):
		return true
	case p.peekTokenIs(token.STREAM):
		return true
	default:
		return false
	}
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	es := &ast.ExpressionStatement{Token: p.curToken, Line: lexer.NUMBER_LINE}
	es.Expression = p.parseExpression(LESSVALUE)
	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return es
}

// 5 + 5 * 6
func (p *Parser) parseExpression(preceden int) ast.Expression {

	if postfix, ok := p.postfixFns[p.peekToken.Type]; ok {
		return postfix()
	}
	// --------------------------------------
	var leftExpression ast.Expression

	prefix := p.prefixFns[p.curToken.Type]
	if prefix == nil {
		if p.curToken.Literal == "int" || p.curToken.Literal == "double" {
			tok := p.curToken
			p.nextToken()
			leftExpression = p.parseCallExpression(&ast.Identifier{Token: tok, Name: tok.Literal, Line: lexer.NUMBER_LINE})
		} else {
			msg := fmt.Sprintf("Line: %d - no prefix parse function for '%s' found, literal: '%s'", lexer.NUMBER_LINE, p.curToken.Type, p.curToken.Literal)
			p.errors = append(p.errors, msg)
			return nil
		}
	}

	if leftExpression == nil {
		leftExpression = prefix()
	}

	for !p.peekTokenIs(token.SEMICOLON) && preceden < p.peekPrecedence() {
		infix := p.infixFns[p.peekToken.Type]
		if infix == nil {
			return leftExpression
		}
		p.nextToken()
		leftExpression = infix(leftExpression)
	}

	return leftExpression
}

//***************************************************************************************
//***************************************************************************************
//***************************************************************************************

func (p *Parser) parsePostfixExpression() ast.Expression {
	postfix := &ast.PostfixExpression{
		Token:    p.peekToken,
		Operator: p.peekToken.Literal,
		Line:     lexer.NUMBER_LINE,
	}

	prefix := p.prefixFns[p.curToken.Type]
	if prefix == nil {
		msg := fmt.Sprintf("Line: %d - no prefix parse function for '%s' found, literal: '%s'", lexer.NUMBER_LINE, p.curToken.Type, p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}

	ident, ok := prefix().(*ast.Identifier)
	if !ok {
		msg := fmt.Sprintf("Line: %d - error expression is not type 'identifier'. got'%T'", lexer.NUMBER_LINE, prefix())
		p.errors = append(p.errors, msg)
		return nil
	}

	postfix.Left = ident
	p.nextToken()
	p.nextToken()
	return postfix
}

func (p *Parser) parseHashExpression() ast.Expression {
	hash := &ast.Hash{Token: p.curToken, Line: lexer.NUMBER_LINE}
	hash.Pairs = make(map[ast.Expression]ast.Expression)

	for !p.peekTokenIs(token.RBRACE) {
		p.nextToken()
		key := p.parseExpression(LESSVALUE)
		if !p.expectedTokenPeek(token.COLON) {
			return nil
		}
		p.nextToken()
		value := p.parseExpression(LESSVALUE)

		hash.Pairs[key] = value
		if !p.peekTokenIs(token.RBRACE) && !p.expectedTokenPeek(token.COMMA) {
			return nil
		}
	}

	if !p.expectedTokenPeek(token.RBRACE) {
		return nil
	}
	return hash
}

func (p *Parser) parseIndexExpression(left ast.Expression) ast.Expression {
	expr := &ast.IndexExpression{Token: p.curToken, Left: left, Line: lexer.NUMBER_LINE}

	p.nextToken()
	expr.Index = p.parseExpression(LESSVALUE)

	if !p.expectedTokenPeek(token.RBRACKET) {
		return nil
	}

	return expr
}

func (p *Parser) parseListExpressions(end token.TokenType) []ast.Expression {
	l := []ast.Expression{}

	if p.peekTokenIs(end) {
		p.nextToken()
		return l
	}

	p.nextToken()
	l = append(l, p.parseExpression(LESSVALUE))

	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		l = append(l, p.parseExpression(LESSVALUE))
	}

	if !p.expectedTokenPeek(end) {
		return nil
	}

	return l
}

func (p *Parser) parseListExpression() ast.Expression {
	l := &ast.List{Token: p.curToken, Line: lexer.NUMBER_LINE}

	l.Elements = p.parseListExpressions(token.RBRACKET)

	return l
}

func (p *Parser) parseCallExpression(function ast.Expression) ast.Expression {
	expr := &ast.CallExpression{Token: p.curToken, Function: function, Line: lexer.NUMBER_LINE}
	expr.Arguments = p.parseListExpressions(token.RPAREN)
	return expr
}

func (p *Parser) parseCallArguments() []ast.Expression {
	args := []ast.Expression{}

	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		return args
	}

	p.nextToken()
	args = append(args, p.parseExpression(LESSVALUE))

	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		args = append(args, p.parseExpression(LESSVALUE))
	}

	if !p.expectedTokenPeek(token.RPAREN) {
		return nil
	}

	return args
}

func (p *Parser) parseFnClosure() ast.Expression {
	fn := &ast.FunctionClosure{Token: p.curToken, Line: lexer.NUMBER_LINE}

	p.nextToken()
	if !p.curTokenIs(token.LPAREN) {
		return nil
	}

	fn.Parameters = p.parseFunctionParametersExpression()
	if !p.curTokenIs(token.RPAREN) {
		return nil
	}

	if p.isPeekBasicType() {
		p.nextToken()
		fn.Type = &ast.Identifier{Token: p.curToken, Name: p.curToken.Literal, Line: lexer.NUMBER_LINE}
		p.nextToken()

	} else if p.peekTokenIs(token.IDENT) {
		msg := fmt.Sprintf("Line: %d - function closure expression is incorrect.", lexer.NUMBER_LINE)
		p.errors = append(p.errors, msg)
		return nil
	} else if p.peekTokenIs(token.LBRACE) {
		fn.Type = nil
		p.nextToken()
	}

	// p.nextToken()
	fn.Body = p.parseBlockStatement()
	if fn.Body == nil {
		return nil
	}
	fn.Body.Type = scope["fn"]
	return fn
}

func (p *Parser) parseFunctionParameters() []*ast.Identifier {
	identifiers := []*ast.Identifier{}

	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		return identifiers
	}

	p.nextToken()
	ident := &ast.Identifier{Token: p.curToken, Name: p.curToken.Literal, Line: lexer.NUMBER_LINE}
	identifiers = append(identifiers, ident)

	for p.peekTokenIs(token.COMMA) {
		p.nextToken() //coma
		p.nextToken() //valor
		ident := &ast.Identifier{Token: p.curToken, Name: p.curToken.Literal, Line: lexer.NUMBER_LINE}
		identifiers = append(identifiers, ident)
	}
	if !p.peekTokenIs(token.RPAREN) {
		return nil
	}
	p.nextToken()

	return identifiers
}

func (p *Parser) parseIfExpression() ast.Expression {
	ie := &ast.IfExpression{Token: p.curToken, Line: lexer.NUMBER_LINE}

	p.nextToken()
	ie.Codition = p.parseExpression(LESSVALUE)
	p.nextToken()

	if !p.curTokenIs(token.LBRACE) {
		msg := fmt.Sprintf("Line: %d - if expression is incorrect, expected token '{'", lexer.NUMBER_LINE)
		p.errors = append(p.errors, msg)
		return nil
	}
	ie.Consequence = p.parseBlockStatement()
	if ie.Consequence == nil {
		return nil
	}
	ie.Consequence.Type = scope["if"]
	if p.peekTokenIs(token.ELSE) {
		p.nextToken()
		if !p.peekTokenIs(token.LBRACE) {
			msg := fmt.Sprintf("Line: %d - else expression is incorrect, expected token '{'", lexer.NUMBER_LINE)
			p.errors = append(p.errors, msg)
			return nil
		}
		p.nextToken()

		ie.Alternative = p.parseBlockStatement()
		if ie.Alternative == nil {
			return nil
		}
		ie.Alternative.Type = scope["if"]
	}

	return ie
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	if !p.curTokenIs(token.LBRACE) {
		msg := fmt.Sprintf("Line: %d - block definition is incorrect, expected token '{'", lexer.NUMBER_LINE)
		p.errors = append(p.errors, msg)
		return nil
	}

	block := &ast.BlockStatement{Token: p.curToken, Line: lexer.NUMBER_LINE}
	block.Statements = []ast.Statement{}
	p.nextToken()

	for !p.curTokenIs(token.RBRACE) {
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.nextToken()
	}

	return block
}

func (p *Parser) parseImplicitExpression(left ast.Expression) ast.Expression {
	ident, ok := left.(*ast.Identifier)
	if !ok {
		msg := fmt.Sprintf("Line: %d - declaration is not possible,  %T type is not IDENTIFIER", lexer.NUMBER_LINE, left)
		p.errors = append(p.errors, msg)
		return nil
	}

	ide := &ast.ImplicitDeclarationExpression{
		Token: p.curToken,
		Left:  ident,
		Line:  lexer.NUMBER_LINE,
	}

	p.nextToken()
	ide.Right = p.parseExpression(LESSVALUE)

	p.nextToken()
	if !p.curTokenIs(token.EOF) && !p.curTokenIs(token.SEMICOLON) && !p.curTokenIs(token.RPAREN) && !p.curTokenIs(token.LBRACE) {
		msg := fmt.Sprintf("Line: %d - implicit assigment incorrect, expected token ';'.\n", lexer.NUMBER_LINE)
		p.errors = append(p.errors, msg)
		return nil
	}

	return ide
}

func (p *Parser) parseAssignOpeExpression(left ast.Expression) ast.Expression {
	ident, ok := left.(*ast.Identifier)
	if !ok {
		msg := fmt.Sprintf("Line: %d - operation is not possible,  %T type is not IDENTIFIER", lexer.NUMBER_LINE, left)
		p.errors = append(p.errors, msg)
		return nil
	}

	aoe := &ast.AssignOperationExpression{
		Token:    p.curToken,
		Left:     ident,
		Operator: p.curToken.Literal,
		Line:     lexer.NUMBER_LINE,
	}

	p.nextToken()
	aoe.Right = p.parseExpression(LESSVALUE)

	p.nextToken()
	if !p.curTokenIs(token.EOF) && !p.curTokenIs(token.SEMICOLON) && !p.curTokenIs(token.RPAREN) && !p.curTokenIs(token.LBRACE) {
		msg := fmt.Sprintf("Line: %d - expression incorrect, expected token ';'.\n", lexer.NUMBER_LINE)
		p.errors = append(p.errors, msg)
		return nil
	}
	return aoe
}

func (p *Parser) parseAssignExpression(left ast.Expression) ast.Expression {
	_, okIdent := left.(*ast.Identifier)
	_, okIndex := left.(*ast.IndexExpression)

	if !okIdent && !okIndex {
		msg := fmt.Sprintf("Line: %d - declaration is not possible\n", lexer.NUMBER_LINE)
		p.errors = append(p.errors, msg)
		return nil
	}

	ae := &ast.AssignExpression{
		Token: p.curToken,
		Left:  left,
		Line:  lexer.NUMBER_LINE,
	}

	p.nextToken()
	ae.Right = p.parseExpression(LESSVALUE)

	p.nextToken()
	if !p.curTokenIs(token.EOF) && !p.curTokenIs(token.SEMICOLON) {
		msg := fmt.Sprintf("Line: %d - expression incorrect, expected token ';'.\n", lexer.NUMBER_LINE)
		p.errors = append(p.errors, msg)
		return nil
	}

	return ae
}

func (p *Parser) parseBooleanExpression() ast.Expression {
	return &ast.Boolean{Token: p.curToken, Value: p.curTokenIs(token.TRUE), Line: lexer.NUMBER_LINE}
}

func (p *Parser) parseParensExpression() ast.Expression {
	p.nextToken()
	expression := p.parseExpression(LESSVALUE)

	if !p.peekTokenIs(token.RPAREN) {
		return nil
	}
	p.nextToken()
	return expression
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	prefix := &ast.PrefixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Line:     lexer.NUMBER_LINE,
	}

	p.nextToken()
	prefix.Right = p.parseExpression(PREFIX)

	return prefix
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	ie := &ast.InfixExpression{
		Token:    p.curToken,
		Left:     left,
		Operator: p.curToken.Literal,
		Line:     lexer.NUMBER_LINE,
	}
	preceden := p.curPrecedence()
	p.nextToken()
	ie.Right = p.parseExpression(preceden)

	return ie
}

func (p *Parser) parseIntegerExpression() ast.Expression {
	i := &ast.Integer{Token: p.curToken, Line: lexer.NUMBER_LINE}
	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("Line: %d - could not parse %q as integer", lexer.NUMBER_LINE, p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}
	i.Value = value
	return i
}

func (p *Parser) parseStringExpression() ast.Expression {
	return &ast.String{Token: p.curToken, Value: p.curToken.Literal, Line: lexer.NUMBER_LINE}
}

func (p *Parser) parseDoubleExpression() ast.Expression {
	d := &ast.Double{Token: p.curToken, Line: lexer.NUMBER_LINE}

	value, err := strconv.ParseFloat(p.curToken.Literal, 64)
	if err != nil {
		msg := fmt.Sprintf("Line: %d - could not parse %q as double", lexer.NUMBER_LINE, p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}
	d.Value = value
	return d
}

func (p *Parser) parseIdentifierExpression() ast.Expression {
	return &ast.Identifier{Token: p.curToken, Name: p.curToken.Literal, Line: lexer.NUMBER_LINE}
}

func (p *Parser) parseNilExpression() ast.Expression {
	return &ast.Nil{Token: p.curToken, Name: p.curToken.Literal, Line: lexer.NUMBER_LINE}
}

//***************************************************************************************
//***************************************************************************************
//***************************************************************************************

func (p *Parser) expectedTokenPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true

	} else {
		p.peekError(t)
		return false
	}
}

func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("Line: %d - expected next token to be '%s', got='%s' instead", lexer.NUMBER_LINE, t, p.peekToken.Literal)
	p.errors = append(p.errors, msg)
}

func (p *Parser) Error() []string {
	return p.errors
}

func (p *Parser) curPrecedence() int {
	value, ok := precedens[p.curToken.Type]
	if !ok {
		return LESSVALUE
	}
	return value
}

func (p *Parser) peekPrecedence() int {
	value, ok := precedens[p.peekToken.Type]
	if !ok {
		return LESSVALUE
	}
	return value
}

func (p *Parser) registerPrefix(tokType token.TokenType, fn prefixFn) {
	p.prefixFns[tokType] = fn
}

func (p *Parser) registerInfix(tokType token.TokenType, fn infixFn) {
	p.infixFns[tokType] = fn
}

func (p *Parser) registerPostfix(tokType token.TokenType, fn postfixFn) {
	p.postfixFns[tokType] = fn
}

func (p *Parser) curTokenIs(tokType token.TokenType) bool {
	return p.curToken.Type == tokType
}

func (p *Parser) peekTokenIs(tokType token.TokenType) bool {
	return p.peekToken.Type == tokType
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.lexer.NextToken()
}
