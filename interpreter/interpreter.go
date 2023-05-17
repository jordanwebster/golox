package interpreter

import (
	"reflect"

	"github.com/jordanwebster/golox/ast"
	"github.com/jordanwebster/golox/token"
)

type Interpreter struct{}

func (interpreter *Interpreter) VisitLiteralExpr(expr *ast.LiteralExpr) interface{} {
	return expr.Value
}

func (interpreter *Interpreter) VisitGroupingExpr(expr *ast.GroupingExpr) interface{} {
	return interpreter.evaluate(expr.Expression)
}

func (interpreter *Interpreter) VisitUnaryExpr(expr *ast.UnaryExpr) interface{} {
	right := interpreter.evaluate(expr.Right)

	switch expr.Operator.Type {
	case token.BANG:
		return !interpreter.isTruthy(right)
	case token.MINUS:
		return -1 * right.(float64)
	}

	return nil
}

func (interpreter *Interpreter) VisitBinaryExpr(expr *ast.BinaryExpr) interface{} {
	left := interpreter.evaluate(expr.Left)
	right := interpreter.evaluate(expr.Right)

	switch expr.Operator.Type {
	case token.GREATER:
		return left.(float64) > right.(float64)
	case token.GREATER_EQUAL:
		return left.(float64) >= right.(float64)
	case token.LESS:
		return left.(float64) < right.(float64)
	case token.LESS_EQUAL:
		return left.(float64) <= right.(float64)
	case token.BANG_EQUAL:
		return !interpreter.isEqual(left, right)
	case token.EQUAL_EQUAL:
		return interpreter.isEqual(left, right)
	case token.MINUS:
		return left.(float64) - right.(float64)
	case token.PLUS:
		leftNumValue, isLeftNumber := left.(float64)
		rightNumValue, isRightNumber := right.(float64)
		if isLeftNumber && isRightNumber {
			return leftNumValue + rightNumValue
		}

		leftStringValue, isLeftString := left.(string)
		rightStringValue, isRightString := right.(string)
		if isLeftString && isRightString {
			return leftStringValue + rightStringValue
		}

		break
	case token.SLASH:
		return left.(float64) / right.(float64)
	case token.STAR:
		return left.(float64) * right.(float64)
	}

	// Unreachable
	return nil
}

func (interpreter *Interpreter) evaluate(expr ast.Expr) interface{} {
	return expr.Accept(interpreter)
}

func (interpreter *Interpreter) isTruthy(object interface{}) bool {
	if object == nil {
		return false
	}
	if value, isBool := object.(bool); isBool {
		return value
	} else {
		return true
	}
}

func (interpreter *Interpreter) isEqual(a interface{}, b interface{}) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil {
		return false
	}

	return reflect.DeepEqual(a, b)
}
