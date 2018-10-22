package evaluator

import (
	"testing"

	"github.com/kenshindeveloper/april/lexer"
	"github.com/kenshindeveloper/april/object"
	"github.com/kenshindeveloper/april/parser"
)

func TestFunctionApplication(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"var identity:func = fn(x:int) int { return x;}; identity(5);", 5},
		{"var identity:func = fn(x:int) int { return x; }; identity(5);", 5},
		{"var add:func = fn(x:int, y:int) int { return x + y;}; add(5, 5);", 10},
		{"var prueba:func = fn(x:int) int { return x * 2;}; prueba(5);", 10},
		{"var add:func = fn(x:int, y:int) int { return x + y;}; add(5 + 5, add(5, 5));", 20},
	}

	for _, tt := range tests {
		testIntegerObject(t, testEval(tt.input), tt.expected)
	}
}

func TestForExpression(t *testing.T) {

}

func TestListIndexExpresssion(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"[1, 2, 3][0]", 1},
		{"[1, 2, 3][1]", 2},
		{"[1, 2, 3][2]", 3},
		{"var i:int = 0; [1][i];", 1},
		{"[1, 2, 3][1 + 1];", 3},
		{"[1, 2, 3][1 + 1];", 3},
		{"[1, 2, 3][3];", nil},
		{"[1, 2, 3][-1];", nil},
	}

	for _, data := range tests {
		evaluated := testEval(data.input)
		integer, ok := data.expected.(int)
		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testNullObject(t, evaluated)
		}
	}
}

func TestListEvaluator(t *testing.T) {
	input := "[1, 2 * 2, 3 + 3]"

	evaluated := testEval(input)
	result, ok := evaluated.(*object.List)
	if !ok {
		t.Fatalf("evaluated is not '*object.List'. got='%T'", evaluated)
	}

	if len(result.Elements) != 3 {
		t.Fatalf("len(result.Elements) is not '%d'. got='%d'", 3, len(result.Elements))
	}

	testIntegerObject(t, result.Elements[0], 1)
	testIntegerObject(t, result.Elements[1], 4)
	testIntegerObject(t, result.Elements[2], 6)
}

func TestIfElseExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"if (true) { 10 }", 10},
		{"if (false) { 10 }", nil},
		{"if (1) { 10 }", 10},
		{"if (1 <= 2) { 10 }", 10},
		{"if (1 > 2) { 10 }", nil},
		{"if (1 > 2) { 10 } else { 20 }", 20},
		{"if (1 < 2) { 10 } else { 20 }", 10},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testNullObject(t, evaluated)
		}
	}
}

func testNullObject(t *testing.T, obj object.Object) bool {
	return obj == NIL
}

func testIntegerObject(t *testing.T, obj object.Object, expected int64) bool {
	result, ok := obj.(*object.Integer)
	if !ok {
		t.Errorf("object is not Integer. got=%T (%+v)", obj, obj)
		return false
	}

	if result.Value != expected {
		t.Errorf("object has wrong value. got=%d, want=%d", result.Value, expected)
		return false
	}

	return true
}

func TestEvalVarStatetement(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"var x:int = 15; x;", 15},
	}

	env := object.NewEnvironment()
	for i, response := range tests {
		l := lexer.New(response.input)
		p := parser.New(l)
		program := p.ParserProgram()
		if program == nil {
			t.Fatalf("program is null.")
		}

		if len(program.Statements) != 2 {
			t.Fatalf("len(program.Statement) is not '%d'. got='%d'.", 2, len(program.Statements))
		}
		evaluated := Eval(program, env)

		if i == 1 {
			integer, ok := evaluated.(*object.Integer)
			if !ok {
				t.Fatalf("evaluated is not '*object.Integer'. got='%T'", evaluated)
			}

			if integer.Value != response.expected {
				t.Fatalf("evaluated is not '%d'. got='%d'", response.expected, integer.Value)
			}
		}
	}
}

func TestNotExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"not true", false},
		{"not false", true},
		{"not not true", true},
		{"not not false", false},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func TestBooleanExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
		{"1 < 2", true},
		{"1 > 2", false},
		{"1 < 1", false},
		{"1 > 1", false},
		{"1 == 1", true},
		{"1 != 1", false},
		{"1 == 2", false},
		{"1 != 2", true},
		{"true == true and false == false", true},
		{"false == false", true},
		{"true == false or true == false", false},
		{"true == false or true == true", true},
		{"true != false", true},
		{"false != true", true},
		{"(1 < 2) == true", true},
		{"(1 <= 2) == false", false},
		{"(1 > 2) == true", false},
		{"(1 >= 2) == false", true},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func testEval(input string) object.Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParserProgram()
	env := object.NewEnvironment()

	return Eval(program, env)
}

func testBooleanObject(t *testing.T, obj object.Object, expected bool) bool {
	result, ok := obj.(*object.Boolean)
	if !ok {
		t.Errorf("object is not Boolean. got=%T (%+v)", obj, obj)
		return false
	}

	if result.Value != expected {
		t.Errorf("object has wrong value. got=%t, want=%t", result.Value, expected)
		return false
	}

	return true
}

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"5", 5},
		{"10", 10},
		{"5", 5},
		{"10", 10},
		{"-5", -5},
		{"-10", -10},
		{"5", 5},
		{"10", 10},
		{"-5", -5},
		{"-10", -10},
		{"8 + 6", 14},
		{"6 + 5 + 5 + 5 - 10", 11},
		{"2 * 2 * 2 * 2 * 2", 32},
		{"5 * 2 + 10", 20},
		{"5 + 2 * 10", 25},
		{"50 / 2 * 2 + 10", 60},
		{"2 * (5 + 10)", 30},
		{"3 * 3 * 3 + 10", 37},
		{"3 * (3 * 3) + 10", 37},
		{"5 + 5 + 5 + 5 - 10", 10},
		{"2 * 2 * 2 * 2 * 2", 32},
		{"-50 + 100 -50", 0},
		{"5 * 2 + 10", 20},
		{"5 + 2 * 10", 25},
		{"20 + 2 * -10", 0},
		{"50 / 2 * 2 + 10", 60},
		{"2 * (5 + 10)", 30},
		{"3 * 3 * 3 + 10", 37},
		{"3 * (3 * 3) + 10", 37},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
	}

	for _, response := range tests {
		l := lexer.New(response.input)
		p := parser.New(l)
		program := p.ParserProgram()
		if program == nil {
			t.Fatalf("ParserProgram is nil.")
		}

		if len(program.Statements) != 1 {
			t.Fatalf("len(program.Statement) is not 1. got=%d", len(program.Statements))
		}

		env := object.NewEnvironment()
		evaluated := Eval(program, env)

		integer, ok := evaluated.(*object.Integer)
		if !ok {
			t.Fatalf("object is not Integer. got=%T", evaluated)
		}

		if integer.Value != response.expected {
			t.Fatalf("object was wrong. got=%d, want=%d", integer.Value, response.expected)
		}
	}
}
