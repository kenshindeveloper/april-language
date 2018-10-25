package evaluator

import (
	"fmt"
	"strings"

	"github.com/kenshindeveloper/april/ast"
	"github.com/kenshindeveloper/april/object"
)

var (
	NIL   = &object.Nil{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {
	//Statements
	case *ast.Program:
		return evalProgram(node, env)

	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)

	case *ast.VarStatement:
		return evalVarStatement(node, env)

	case *ast.GlobalStatement:
		return evalGlobalStatement(node, env)

	case *ast.ReturnStatement:
		return evalReturnStatement(node, env)

	case *ast.BreakStatement:
		return evalBreakStatement(node, env)

	case *ast.BlockStatement:
		return evalBlockStatement(node, env)

	case *ast.ForStatement:
		return evalForStatement(node, env)

	case *ast.Function:
		return evalFunctionStatement(node, env)

	//Expressions
	case *ast.Nil:
		return &object.Nil{}

	case *ast.Hash:
		return evalHashExpression(node, env)

	case *ast.IndexExpression:
		return evalIndexExpressions(node, env)

	case *ast.List:
		return evalListExpression(node, env)

	case *ast.CallExpression:
		return evalCallExpression(node, env)

	case *ast.FunctionClosure:
		return evalFunctionClosureExpression(node, env)

	case *ast.IfExpression:
		return evalIfExpression(node, env)

	case *ast.Identifier:
		return evalIdentifier(node, env)

	case *ast.Integer:
		return &object.Integer{Value: node.Value}

	case *ast.Boolean:
		return boolToBooleanObject(node.Value)

	case *ast.String:
		return &object.String{Value: node.Value}

	case *ast.Double:
		return &object.Double{Value: node.Value}

	case *ast.ImplicitDeclarationExpression:
		return evalImplicitDeclarationExpression(node, env)

	case *ast.AssignExpression:
		return evalAssignExpression(node, env)

	case *ast.AssignOperationExpression:
		return evalAssignOperationExpression(node, env)

	case *ast.PrefixExpression:
		return evalPrefixExpressions(node, env)

	case *ast.InfixExpression:
		return evalInfixExpressions(node, env)

	case *ast.PostfixExpression:
		return evalPostfixExpressions(node, env)
	}

	return nil
}

//***************************************************************************************
//***************************************STATEMENTS**************************************
//***************************************************************************************

func evalFunctionStatement(node *ast.Function, env *object.Environment) object.Object {

	if !env.Scope {
		return newError("Line: %d - name function '%s' cannot declared in this scope. ", node.Line, node.Name.Name)
	}

	_, inEnv := env.Get(node.Name.Name)
	_, inBuilt := builtins[node.Name.Name]
	if inEnv || inBuilt {
		return newError("Line: %d - name function '%s' already exist. ", node.Line, node.Name.Name)
	}

	for _, param := range node.Parameters {
		_, inEnv = env.Get(param.Name.Name)
		_, inBuilt = builtins[param.Name.Name]
		if inEnv || inBuilt {
			return newError("Line: %d - name variable '%s' already exist. ", node.Line, param.Name.Name)
		}
	}

	fn := &object.Function{Name: node.Name, Parameters: node.Parameters, Return: node.Type, Body: node.Body, Env: env}
	env.SaveGlobal(node.Name.Name, fn)

	return NIL
}

func evalVarStatement(node *ast.VarStatement, env *object.Environment) object.Object {
	val := Eval(node.Value, env)
	if isError(val) {
		return val
	}

	_, inEnv := env.Get(node.Name.Name)
	_, inBuilt := builtins[node.Name.Name]

	if !inEnv && !inBuilt {
		switch val.(type) {
		case *object.Integer:
			if node.Type.Name != "int" {
				return newError("Line: %d - declaration error: var %s:%s = INTEGER", node.Line, node.Name.Name, strings.ToUpper(node.Type.Name))
			}
		case *object.Double:
			if node.Type.Name != "double" {
				return newError("Line: %d - declaration error: var %s:%s = DOUBLE", node.Line, node.Name.Name, strings.ToUpper(node.Type.Name))
			}
		case *object.Boolean:
			if node.Type.Name != "bool" {
				return newError("Line: %d - declaration error: var %s:%s = BOOL", node.Line, node.Name.Name, strings.ToUpper(node.Type.Name))
			}
		case *object.String:
			if node.Type.Name != "string" {
				return newError("Line: %d - declaration error: var %s:%s = STRING", node.Line, node.Name.Name, strings.ToUpper(node.Type.Name))
			}
		case *object.List:
			if node.Type.Name != "list" {
				return newError("Line: %d - declaration error: var %s:%s = LIST", node.Line, node.Name.Name, strings.ToUpper(node.Type.Name))
			}
		case *object.Hash:
			if node.Type.Name != "map" {
				return newError("Line: %d - declaration error: var %s:%s = MAP", node.Line, node.Name.Name, strings.ToUpper(node.Type.Name))
			}
		case *object.Stream:
			if node.Type.Name != "stream" {
				return newError("Line: %d - declaration error: var %s:%s = STREAM", node.Line, node.Name.Name, strings.ToUpper(node.Type.Name))
			}
		case *object.FunctionClosure:
			if node.Type.Name != "func" {
				return newError("Line: %d - declaration error: var %s:%s = FUNC", node.Line, node.Name.Name, strings.ToUpper(node.Type.Name))
			}
		}

		env.Save(node.Name.Name, val)
	} else {
		return newError("Line: %d - variable '%s' already exist. ", node.Line, node.Name.Name)
	}
	return NIL
}

func evalGlobalStatement(node *ast.GlobalStatement, env *object.Environment) object.Object {
	val := Eval(node.Value, env)
	if isError(val) {
		return val
	}

	_, inEnv := env.Get(node.Name.Name)
	_, inBuilt := builtins[node.Name.Name]

	if !inEnv && !inBuilt {
		switch val.(type) {
		case *object.Integer:
			if node.Type.Name != "int" {
				return newError("Line: %d - declaration error: var %s:%s = INTEGER", node.Line, node.Name.Name, strings.ToUpper(node.Type.Name))
			}
		case *object.Double:
			if node.Type.Name != "double" {
				return newError("Line: %d - declaration error: var %s:%s = DOUBLE", node.Line, node.Name.Name, strings.ToUpper(node.Type.Name))
			}
		case *object.Boolean:
			if node.Type.Name != "bool" {
				return newError("Line: %d - declaration error: var %s:%s = BOOL", node.Line, node.Name.Name, strings.ToUpper(node.Type.Name))
			}
		case *object.String:
			if node.Type.Name != "string" {
				return newError("Line: %d - declaration error: var %s:%s = STRING", node.Line, node.Name.Name, strings.ToUpper(node.Type.Name))
			}
		case *object.List:
			if node.Type.Name != "list" {
				return newError("Line: %d - declaration error: var %s:%s = LIST", node.Line, node.Name.Name, strings.ToUpper(node.Type.Name))
			}
		case *object.Hash:
			if node.Type.Name != "map" {
				return newError("Line: %d - declaration error: var %s:%s = MAP", node.Line, node.Name.Name, strings.ToUpper(node.Type.Name))
			}
		case *object.Stream:
			if node.Type.Name != "stream" {
				return newError("Line: %d - declaration error: var %s:%s = STREAM", node.Line, node.Name.Name, strings.ToUpper(node.Type.Name))
			}
		case *object.FunctionClosure:
			if node.Type.Name != "func" {
				return newError("Line: %d - declaration error: var %s:%s = FUNC", node.Line, node.Name.Name, strings.ToUpper(node.Type.Name))
			}
		}

		env.SaveGlobal(node.Name.Name, val)
	} else {
		return newError("Line: %d - variable '%s' already exist. ", node.Line, node.Name.Name)
	}
	return NIL
}

func evalReturnStatement(node *ast.ReturnStatement, env *object.Environment) object.Object {
	expr := Eval(node.Expression, env)
	if isError(expr) {
		return expr
	}

	return &object.ReturnStatement{Value: expr}
}

func evalBreakStatement(node *ast.BreakStatement, env *object.Environment) object.Object {
	return &object.BreakStatement{}
}

func evalForStatement(node *ast.ForStatement, env *object.Environment) object.Object {
	extendEnv := object.NewEncloseEnvironment(env)
	var condition object.Object
	//---------------------------------
	if impl, ok := node.Condition.(*ast.ImplicitDeclarationExpression); ok {

		if node.Declaration != nil || node.Operation != nil {
			return newError("Line: %d - for declaration is incorrect.", node.Line)
		}

		var listExpr object.Object
		switch impl.Right.(type) {
		case *ast.CallExpression:
			listExpr = Eval(impl.Right, extendEnv)
		case *ast.List:
			listExpr = Eval(impl.Right, extendEnv)
		case *ast.Identifier:
			obj, ok := env.Get(impl.Right.(*ast.Identifier).Name)
			if ok {
				_, ok = obj.(*object.List)
				if !ok {
					return newError("Line: %d - expression incompatible with for, expression must be type list. %T", node.Line, impl.Right)
				}
				listExpr = Eval(impl.Right, extendEnv)
			}
		default:
			return newError("Line: %d - expression incompatible with for, expression must be type list. %T", node.Line, impl.Right)
		}

		if isError(listExpr) {
			return listExpr
		}

		if list, ok := listExpr.(*object.List); ok {
			extendEnv.Save(impl.Left.Name, &object.Integer{})
			for _, obj := range list.Elements {
				newExtendEnv := object.NewEncloseEnvironment(extendEnv)
				extendEnv.Set(impl.Left.Name, obj)
				body := Eval(node.Body, newExtendEnv)
				if isError(body) {
					return body
				} else if body == nil {
					break
				} else if body != nil && body.Type() == object.RETURN_OBJ || body.Type() == object.BREAK_OBJ {
					return body
				}
			}
			return NIL
		}
		return newError("Line: %d - expression is not type list.", node.Line)
	}
	//---------------------------------
	var declaration object.Object
	if node.Declaration != nil {
		declaration = Eval(node.Declaration, extendEnv)
		if isError(declaration) {
			return declaration
		}
	}

	//----------------------------------
	condition = Eval(node.Condition, extendEnv)
	if isError(condition) {
		return condition
	}

	value, ok := condition.(*object.Boolean)
	if !ok {
		return newError("Line: %d - the expression is not type boolean.", node.Line)
	}
	run := value.Value
	for run {
		newExtendEnv := object.NewEncloseEnvironment(extendEnv)
		body := Eval(node.Body, newExtendEnv)
		if isError(body) {
			return body
		} else if body == nil {
			break
		} else if body != nil && body.Type() == object.RETURN_OBJ || body.Type() == object.BREAK_OBJ {
			return body
		}

		if node.Operation != nil {
			operation := Eval(node.Operation, extendEnv)
			if isError(operation) {
				return operation
			}
		}
		run = isTruthy(Eval(node.Condition, extendEnv))
	}

	return NIL
}

//***************************************************************************************
//************************************EXPRESSIONS****************************************
//***************************************************************************************

func evalIndexExpressions(node *ast.IndexExpression, env *object.Environment) object.Object {
	left := Eval(node.Left, env)
	if isError(left) {
		return left
	}

	index := Eval(node.Index, env)
	if isError(index) {
		return index
	}

	return evalIndexExpression(left, index)
}

func evalListExpression(node *ast.List, env *object.Environment) object.Object {
	elements := evalExpressions(node.Elements, env)
	if len(elements) == 1 && isError(elements[0]) {
		return elements[0]
	}

	return &object.List{Elements: elements}
}

func evalCallExpression(node *ast.CallExpression, env *object.Environment) object.Object {
	function := Eval(node.Function, env)
	if isError(function) {
		return function
	}
	args := evalExpressions(node.Arguments, env)
	if len(args) == 1 && isError(args[0]) {
		return args[0]
	}

	switch function.(type) {
	case *object.Function:
		if len(function.(*object.Function).Parameters) != len(args) {
			return newError("Line: %d - count parameters not match.", node.Line)
		}
	case *object.FunctionClosure:
		if len(function.(*object.FunctionClosure).Parameters) != len(args) {
			return newError("Line: %d - count parameters not match.", node.Line)
		}
		// default:
		// 	return newError("data type error in call function.")
	}

	return applyFunction(function, args)
}

func evalFunctionClosureExpression(node *ast.FunctionClosure, env *object.Environment) object.Object {
	params := node.Parameters
	body := node.Body
	return &object.FunctionClosure{Parameters: params, Return: node.Type, Env: env, Body: body}
}

func evalImplicitDeclarationExpression(node *ast.ImplicitDeclarationExpression, env *object.Environment) object.Object {
	right := Eval(node.Right, env)
	if isError(right) {
		return right
	}

	_, inEnv := env.Get(node.Left.Name)
	_, inBuilt := builtins[node.Left.Name]

	if !inEnv && !inBuilt {

		if !isBasicDataType(right) {
			return newError("Line: %d - declaration not compatible '%s' := '%s'. ", node.Line, node.Left.Name, right.Type())
		}

		env.Save(node.Left.Name, right)
	} else {
		return newError("Line: %d - variable '%s' already exist. ", node.Line, node.Left.Name)
	}

	return NIL
}

func isBasicDataType(obj object.Object) bool {
	switch obj.(type) {
	case *object.Integer:
		return true
	case *object.Boolean:
		return true
	case *object.String:
		return true
	case *object.Double:
		return true
	case *object.List:
		return true
	case *object.Hash:
		return true
	case *object.FunctionClosure:
		return true
	case *object.Stream:
		return true
	default:
		return false
	}
}

func evalAssignExpression(node *ast.AssignExpression, env *object.Environment) object.Object {
	right := Eval(node.Right, env)
	if isError(right) {
		return right
	}

	switch node.Left.(type) {
	case *ast.Identifier:
		ident := node.Left.(*ast.Identifier)
		if value, ok := env.Get(ident.Name); ok {
			if right.Type() != value.Type() {
				switch {
				case value.Type() == object.DOUBLE_OBJ && right.Type() == object.INTEGER_OBJ:
					env.Set(ident.Name, &object.Double{Value: float64(right.(*object.Integer).Value)})
				default:
					return newError("Line: %d - assign not compatible '%s' : '%s'. ", node.Line, right.Type(), value.Type())
				}
			} else {
				env.Set(ident.Name, right)
			}
		} else {
			return newError("Line: %d - variable '%s' not exist. ", node.Line, ident.Name)
		}

		return NIL
	case *ast.IndexExpression:
		indexExpr := node.Left.(*ast.IndexExpression)
		left := Eval(indexExpr.Left, env)
		if isError(left) {
			return left
		}

		index := Eval(indexExpr.Index, env)
		if isError(index) {
			return index
		}

		value := Eval(node.Right, env)
		if isError(value) {
			return value
		}

		return evalSetIndexExpression(left, index, value)

	default:
		return newError("Line: %d - expression assignment is not possible. ", node.Line)
	}
}

func evalAssignOperationExpression(node *ast.AssignOperationExpression, env *object.Environment) object.Object {
	right := Eval(node.Right, env)
	if isError(right) {
		return right
	}

	if value, ok := env.Get(node.Left.Name); ok {
		v := evalInfixExpression(string(node.Operator[0]), value, right)
		if isError(v) {
			return v
		}
		env.Set(node.Left.Name, v)
		return v
	}

	return newError("Line: %d - variable '%s' not exist. ", node.Line, node.Left.Name)
}

func evalPrefixExpressions(node *ast.PrefixExpression, env *object.Environment) object.Object {
	right := Eval(node.Right, env)
	if isError(right) {
		return right
	}
	return evalPrefixExpression(node.Operator, right)
}

func evalInfixExpressions(node *ast.InfixExpression, env *object.Environment) object.Object {
	left := Eval(node.Left, env)
	if isError(left) {
		return left
	}

	//--------------------------------------
	//--------------------------------------
	if left.Type() == object.BOOLEAN_OBJ && left.(*object.Boolean).Value == false && node.Operator == "and" {
		_, ok := node.Right.(*ast.Boolean)
		if ok {
			return boolToBooleanObject(false)
		}

		_, ok = node.Right.(*ast.InfixExpression)
		if ok {
			return boolToBooleanObject(false)
		}
	}
	//--------------------------------------
	//--------------------------------------

	right := Eval(node.Right, env)
	if isError(right) {
		return right
	}

	return evalInfixExpression(node.Operator, left, right)
}

func evalPostfixExpressions(node *ast.PostfixExpression, env *object.Environment) object.Object {
	value, ok := env.Get(node.Left.Name)
	if ok {
		_, isInteger := value.(*object.Integer)
		if !isInteger {
			return newError("Line: %d - operator '%s' must be integer. got='%T'", node.Line, node.Operator, value)
		}

		v := evalInfixExpression(string(node.Operator[0]), value, &object.Integer{Value: 1})
		if isError(v) {
			return v
		}
		env.Set(node.Left.Name, v)
		return v
	}
	return newError("Line: %d - variable '%s' not exist. ", node.Line, node.Left.Name)
}

//***************************************************************************************
//***************************************************************************************
//***************************************************************************************

func evalBlockStatement(block *ast.BlockStatement, env *object.Environment) object.Object {
	var result object.Object

	for _, stmt := range block.Statements {
		result = Eval(stmt, env)
		if result != nil {
			if result.Type() == object.RETURN_OBJ || result.Type() == object.ERROR_OBJ || block.Type > 0 && result.Type() == object.BREAK_OBJ {
				return result
			}
		}
	}

	return result
}

func evalIfExpression(ie *ast.IfExpression, env *object.Environment) object.Object {
	condition := Eval(ie.Codition, env)
	if isError(condition) {
		return condition
	}
	flag := env.Scope
	env.Scope = false
	if isTruthy(condition) {
		eval := Eval(ie.Consequence, env)
		if isError(eval) {
			return eval
		}
		env.Scope = flag
		return eval
	} else if ie.Alternative != nil {
		eval := Eval(ie.Alternative, env)
		if isError(eval) {
			return eval
		}
		env.Scope = flag
		return eval
	} else {
		return NIL
	}
}

func isTruthy(obj object.Object) bool {
	switch obj {
	case NIL:
		return false
	case TRUE:
		return true
	case FALSE:
		return false
	default:
		return true
	}
}

func evalProgram(program *ast.Program, env *object.Environment) object.Object {
	var result object.Object

	for _, statement := range program.Statements {
		result = Eval(statement, env)

		switch result := result.(type) {
		case *object.Error:
			return result
		case *object.ReturnStatement:
			return result.Value
		}
	}

	return result
}

func evalIdentifier(node *ast.Identifier, env *object.Environment) object.Object {
	if val, ok := env.Get(node.Name); ok {
		return val
	}

	if builtin, ok := builtins[node.Name]; ok {
		return builtin
	}

	return newError("Line: %d - identifier not found: %s", node.Line, node.Name)
}

func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "not":
		return evalNotOperatorExpression(right)
	case "-":
		return evalMinOperatorExpression(right)
	default:
		return newError("unknown operator: %s%s", operator, right.Type())
	}
}

func evalMinOperatorExpression(right object.Object) object.Object {
	switch right.Type() {
	case object.INTEGER_OBJ:
		return &object.Integer{Value: -right.(*object.Integer).Value}
	case object.DOUBLE_OBJ:
		return &object.Double{Value: -right.(*object.Double).Value}
	default:
		return newError("unkown operator: -%s", right.Type())
	}
}

func evalInfixExpression(operator string, left object.Object, right object.Object) object.Object {
	switch {
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalIntegerInfixExpression(operator, left, right)
	case left.Type() == object.STRING_OBJ && right.Type() == object.STRING_OBJ:
		return evalStringInfixExpression(operator, left, right)
	case left.Type() == object.BOOLEAN_OBJ && right.Type() == object.BOOLEAN_OBJ:
		return evalBooleanInfixExpression(operator, left, right)
	case left.Type() == object.DOUBLE_OBJ && right.Type() == object.DOUBLE_OBJ:
		return evalDoubleInfixExpression(operator, left.(*object.Double).Value, right.(*object.Double).Value)
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.DOUBLE_OBJ:
		leftVar := left.(*object.Integer).Value
		return evalDoubleInfixExpression(operator, float64(leftVar), right.(*object.Double).Value)
	case left.Type() == object.DOUBLE_OBJ && right.Type() == object.INTEGER_OBJ:
		rightVar := right.(*object.Integer).Value
		return evalDoubleInfixExpression(operator, left.(*object.Double).Value, float64(rightVar))
	case left.Type() == object.NIL_OBJ || right.Type() == object.NIL_OBJ:
		return evalNilInfixExpression(operator, left, right)
	case left.Type() != right.Type():
		return newError("type mismatch: %s %s %s", left.Type(), operator, right.Type())
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalNilInfixExpression(operator string, left object.Object, right object.Object) object.Object {
	switch operator {
	case "!=":
		return boolToBooleanObject(left.Type() != right.Type())
	case "==":
		return boolToBooleanObject(left.Type() == right.Type())
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalDoubleInfixExpression(operator string, leftVar float64, rightVar float64) object.Object {
	switch operator {
	case "+":
		return &object.Double{Value: (leftVar + rightVar)}
	case "-":
		return &object.Double{Value: (leftVar - rightVar)}
	case "*":
		return &object.Double{Value: (leftVar * rightVar)}
	case "/":
		if rightVar == 0 {
			return newError("division by zero")
		}
		return &object.Double{Value: (leftVar / rightVar)}
	case "<":
		return boolToBooleanObject(leftVar < rightVar)
	case ">":
		return boolToBooleanObject(leftVar > rightVar)
	case "!=":
		return boolToBooleanObject(leftVar != rightVar)
	case "==":
		return boolToBooleanObject(leftVar == rightVar)
	case "<=":
		return boolToBooleanObject(leftVar <= rightVar)
	case ">=":
		return boolToBooleanObject(leftVar >= rightVar)
	default:
		return newError("unknown operator: DOUBLE %s DOUBLE", operator)
	}
}

func evalBooleanInfixExpression(operator string, left object.Object, right object.Object) object.Object {
	leftVar := left.(*object.Boolean).Value
	rightVar := right.(*object.Boolean).Value

	switch operator {
	case "==":
		return boolToBooleanObject(leftVar == rightVar)
	case "!=":
		return boolToBooleanObject(leftVar != rightVar)
	case "and":
		return boolToBooleanObject(leftVar && rightVar)
	case "or":
		return boolToBooleanObject(leftVar || rightVar)
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalStringInfixExpression(operator string, left object.Object, right object.Object) object.Object {
	leftVar := left.(*object.String).Value
	rightVar := right.(*object.String).Value

	switch operator {
	case "+":
		return &object.String{Value: (leftVar + rightVar)}
	case "==":
		return boolToBooleanObject(leftVar == rightVar)
	case "!=":
		return boolToBooleanObject(leftVar != rightVar)
	case "<=":
		return boolToBooleanObject(leftVar <= rightVar)
	case ">=":
		return boolToBooleanObject(leftVar >= rightVar)
	case "<":
		return boolToBooleanObject(leftVar < rightVar)
	case ">":
		return boolToBooleanObject(leftVar > rightVar)
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalIntegerInfixExpression(operator string, left object.Object, right object.Object) object.Object {
	leftVar := left.(*object.Integer).Value
	rightVar := right.(*object.Integer).Value

	switch operator {
	case "+":
		return &object.Integer{Value: (leftVar + rightVar)}
	case "-":
		return &object.Integer{Value: (leftVar - rightVar)}
	case "*":
		return &object.Integer{Value: (leftVar * rightVar)}
	case "/":
		if rightVar == 0 {
			return newError("division by zero")
		}
		return &object.Integer{Value: (leftVar / rightVar)}
	case "%":
		return &object.Integer{Value: (leftVar % rightVar)}
	case "<":
		return boolToBooleanObject(leftVar < rightVar)
	case ">":
		return boolToBooleanObject(leftVar > rightVar)
	case "!=":
		return boolToBooleanObject(leftVar != rightVar)
	case "==":
		return boolToBooleanObject(leftVar == rightVar)
	case "<=":
		return boolToBooleanObject(leftVar <= rightVar)
	case ">=":
		return boolToBooleanObject(leftVar >= rightVar)
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func boolToBooleanObject(value bool) object.Object {
	if value {
		return TRUE
	}
	return FALSE
}

func evalNotOperatorExpression(right object.Object) object.Object {
	switch right {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	case NIL:
		return TRUE
	default:
		return FALSE
	}
}

func evalExpressions(exprs []ast.Expression, env *object.Environment) []object.Object {
	var result []object.Object

	for _, e := range exprs {
		evaluated := Eval(e, env)
		if isError(evaluated) {
			return []object.Object{evaluated}
		}
		result = append(result, evaluated)
	}

	return result
}

func applyFunction(fn object.Object, args []object.Object) object.Object {
	switch fn := fn.(type) {
	case *object.FunctionClosure:
		if extendedEnv := extendFunctionClosureEnv(fn, args); extendedEnv != nil {
			evaluated := Eval(fn.Body, extendedEnv)
			if isError(evaluated) {
				return evaluated
			}
			return unwrapReturnValue(fn, evaluated)
		}
		return newError("data type mismatch into function closure.")

	case *object.Function:
		if extendedEnv := extendFunctionEnv(fn, args); extendedEnv != nil {
			evaluated := Eval(fn.Body, extendedEnv)
			if isError(evaluated) {
				return evaluated
			}
			return unwrapReturnFunctionValue(fn, evaluated)
		}
		return newError("data type mismatch into function '%s'", fn.Name.Name)
	case *object.Builtin:
		return fn.Fn(args...)
	default:
		return newError("not a function: %s", fn.Type())
	}
}

func unwrapReturnFunctionValue(fn *object.Function, obj object.Object) object.Object {
	if fn.Return == nil {
		if _, ok := obj.(*object.ReturnStatement); ok {
			return newError("return value not mismatch, expected 'null'.")
		}
		return NIL
	}
	if returnValue, ok := obj.(*object.ReturnStatement); ok {
		if fn.Return.Name == "int" && returnValue.Value.Type() == object.INTEGER_OBJ {
			return returnValue.Value
		} else if fn.Return.Name == "double" && (returnValue.Value.Type() == object.DOUBLE_OBJ || returnValue.Value.Type() == object.INTEGER_OBJ) {
			if returnValue.Value.Type() == object.INTEGER_OBJ {
				return &object.Double{Value: float64(returnValue.Value.(*object.Integer).Value)}
			}
			return returnValue.Value
		} else if fn.Return.Name == "bool" && returnValue.Value.Type() == object.BOOLEAN_OBJ {
			return returnValue.Value
		} else if fn.Return.Name == "string" && returnValue.Value.Type() == object.STRING_OBJ {
			return returnValue.Value
		} else if fn.Return.Name == "list" && returnValue.Value.Type() == object.LIST_OBJ {
			return returnValue.Value
		} else if fn.Return.Name == "map" && returnValue.Value.Type() == object.HASH_OBJ {
			return returnValue.Value
		} else if fn.Return.Name == "func" && returnValue.Value.Type() == object.CLOSURE_OBJ {
			return returnValue.Value
		}
	}
	return newError("return value not mismatch, expected '%s'.", fn.Return.Name)
}

func extendFunctionEnv(fn *object.Function, args []object.Object) *object.Environment {
	env := object.NewEnvironmentFn(fn.Env)
	save := false

	for pos, param := range fn.Parameters {
		if param.Type.Name == "int" && args[pos].Type() == object.INTEGER_OBJ {
			save = true
		} else if param.Type.Name == "double" && (args[pos].Type() == object.DOUBLE_OBJ || args[pos].Type() == object.INTEGER_OBJ) {
			if args[pos].Type() == object.INTEGER_OBJ {
				args[pos] = &object.Double{Value: float64(args[pos].(*object.Integer).Value)}
			}
			save = true
		} else if param.Type.Name == "bool" && args[pos].Type() == object.BOOLEAN_OBJ {
			save = true
		} else if param.Type.Name == "string" && args[pos].Type() == object.STRING_OBJ {
			save = true
		} else if param.Type.Name == "list" && args[pos].Type() == object.LIST_OBJ {
			save = true
		} else if param.Type.Name == "map" && args[pos].Type() == object.HASH_OBJ {
			save = true
		} else if param.Type.Name == "func" && args[pos].Type() == object.CLOSURE_OBJ {
			save = true
		}

		if save {
			env.Save(param.Name.Name, args[pos])
			save = false
		} else {
			return nil
		}
	}

	return env
}

func extendFunctionClosureEnv(fn *object.FunctionClosure, args []object.Object) *object.Environment {
	env := object.NewEncloseEnvironment(fn.Env)
	save := false

	for pos, param := range fn.Parameters {
		if param.Type.Name == "int" && args[pos].Type() == object.INTEGER_OBJ {
			save = true
		} else if param.Type.Name == "double" && (args[pos].Type() == object.DOUBLE_OBJ || args[pos].Type() == object.INTEGER_OBJ) {
			if args[pos].Type() == object.INTEGER_OBJ {
				args[pos] = &object.Double{Value: float64(args[pos].(*object.Integer).Value)}
			}
			save = true
		} else if param.Type.Name == "bool" && args[pos].Type() == object.BOOLEAN_OBJ {
			save = true
		} else if param.Type.Name == "string" && args[pos].Type() == object.STRING_OBJ {
			save = true
		} else if param.Type.Name == "list" && args[pos].Type() == object.LIST_OBJ {
			save = true
		} else if param.Type.Name == "map" && args[pos].Type() == object.HASH_OBJ {
			save = true
		} else if param.Type.Name == "func" && args[pos].Type() == object.CLOSURE_OBJ {
			save = true
		}

		if save {
			env.Save(param.Name.Name, args[pos])
			save = false
		} else {
			return nil
		}
	}

	return env
}

func unwrapReturnValue(fn *object.FunctionClosure, obj object.Object) object.Object {
	if fn.Return == nil {
		if _, ok := obj.(*object.ReturnStatement); ok {
			return newError("return value not mismatch, expected 'null'.")
		}
		return NIL
	}
	if returnValue, ok := obj.(*object.ReturnStatement); ok {
		if fn.Return.Name == "int" && returnValue.Value.Type() == object.INTEGER_OBJ {
			return returnValue.Value
		} else if fn.Return.Name == "double" && (returnValue.Value.Type() == object.DOUBLE_OBJ || returnValue.Value.Type() == object.INTEGER_OBJ) {
			if returnValue.Value.Type() == object.INTEGER_OBJ {
				return &object.Double{Value: float64(returnValue.Value.(*object.Integer).Value)}
			}
			return returnValue.Value
		} else if fn.Return.Name == "bool" && returnValue.Value.Type() == object.BOOLEAN_OBJ {
			return returnValue.Value
		} else if fn.Return.Name == "string" && returnValue.Value.Type() == object.STRING_OBJ {
			return returnValue.Value
		} else if fn.Return.Name == "list" && returnValue.Value.Type() == object.LIST_OBJ {
			return returnValue.Value
		} else if fn.Return.Name == "map" && returnValue.Value.Type() == object.HASH_OBJ {
			return returnValue.Value
		}
	}
	return newError("return value not mismatch, expected '%s'.", fn.Return.Name)
}

func evalSetIndexExpression(left, index, value object.Object) object.Object {
	switch {
	case left.Type() == object.LIST_OBJ && index.Type() == object.INTEGER_OBJ:
		return evalSetListIndexExpression(left, index, value)
	case left.Type() == object.HASH_OBJ:
		return evalSetHashIndexExpression(left, index, value)
	default:
		return newError("Index operator not supported %s", left.Inspect())
	}
}

func evalSetHashIndexExpression(hash, index, value object.Object) object.Object {
	hashObj := hash.(*object.Hash)
	key, ok := index.(object.Hashable)
	if !ok {
		return newError("unusable as hash key: %s", index.Type())
	}

	pair, ok := hashObj.Pairs[key.HashKey()]
	if ok {
		delete(hashObj.Pairs, key.HashKey())
	}

	hashObj.Pairs[key.HashKey()] = object.HashPair{Key: index, Value: value}
	return pair.Value
}

func evalSetListIndexExpression(left, index, value object.Object) object.Object {
	listObject := left.(*object.List)
	pos := index.(*object.Integer).Value
	max := int64(len(listObject.Elements) - 1)

	if pos < 0 || pos > max {
		return newError("list index out of range.")
		// return NIL
	}

	listObject.Elements[pos] = value
	return listObject.Elements[pos]
}

func evalIndexExpression(left, index object.Object) object.Object {
	switch {
	case left.Type() == object.LIST_OBJ && index.Type() == object.INTEGER_OBJ:
		return evalListIndexExpression(left, index)
	case left.Type() == object.HASH_OBJ:
		return evalHashIndexExpression(left, index)
	case left.Type() == object.STRING_OBJ && index.Type() == object.INTEGER_OBJ:
		return evalStringIndexExpression(left, index)
	default:
		return newError("Index operator not supported %s", left.Inspect())
	}
}

func evalStringIndexExpression(left, index object.Object) object.Object {
	stringObject := left.(*object.String)
	pos := index.(*object.Integer).Value
	max := int64(len(stringObject.Value) - 1)

	if pos < 0 || pos > max {
		return newError("string index out of range.")
		// return NIL
	}

	return &object.String{Value: string(stringObject.Value[pos])}
}

func evalListIndexExpression(left, index object.Object) object.Object {
	listObject := left.(*object.List)
	pos := index.(*object.Integer).Value
	max := int64(len(listObject.Elements) - 1)

	if pos < 0 || pos > max {
		return newError("list index out of range.")
		// return NIL
	}

	return listObject.Elements[pos]
}

func evalHashExpression(node *ast.Hash, env *object.Environment) object.Object {
	pairs := make(map[object.HashKey]object.HashPair)

	for keyNode, valueNode := range node.Pairs {
		key := Eval(keyNode, env)
		if isError(key) {
			return key
		}

		hashkey, ok := key.(object.Hashable)
		if !ok {
			return newError("unusable as hash key: %s", key.Type())
		}

		value := Eval(valueNode, env)
		if isError(value) {
			return value
		}

		hashed := hashkey.HashKey()
		pairs[hashed] = object.HashPair{Key: key, Value: value}
	}

	return &object.Hash{Pairs: pairs}
}

func evalHashIndexExpression(hash, index object.Object) object.Object {
	hashObj := hash.(*object.Hash)

	key, ok := index.(object.Hashable)
	if !ok {
		return newError("unusable as hash key: %s", index.Type())
	}

	pair, ok := hashObj.Pairs[key.HashKey()]
	if !ok {
		return newError("key error: %s", index.Inspect())
	}

	return pair.Value
}

//***************************************************************************************
//***************************************************************************************
//***************************************************************************************

func newError(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR_OBJ
	}

	return false
}
