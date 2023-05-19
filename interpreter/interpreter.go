package interpreter

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/jordanwebster/golox/ast"
	"github.com/jordanwebster/golox/loxerror"
	"github.com/jordanwebster/golox/token"
)

type Interpreter struct{}

func NewInterpreter() *Interpreter {
	return &Interpreter{}
}

func (interpreter *Interpreter) Interpret(statements []ast.Stmt) {
	for _, stmt := range statements {
		err := interpreter.execute(stmt)
		if err != nil {
			switch err.(type) {
			case *loxerror.RuntimeError:
				loxerror.ReportRuntimeError(err.(*loxerror.RuntimeError))
			default:
				panic(err)
			}
		}
	}
}

func (interpreter *Interpreter) VisitLiteralExpr(expr *ast.LiteralExpr) (interface{}, error) {
	return expr.Value, nil
}

func (interpreter *Interpreter) VisitGroupingExpr(expr *ast.GroupingExpr) (interface{}, error) {
	return interpreter.evaluate(expr.Expression)
}

func (interpreter *Interpreter) VisitUnaryExpr(expr *ast.UnaryExpr) (interface{}, error) {
	right, err := interpreter.evaluate(expr.Right)
	if err != nil {
		return nil, err
	}

	switch expr.Operator.Type {
	case token.BANG:
		return !isTruthy(right), nil
	case token.MINUS:
		err := checkNumberOperand(expr.Operator, right)
		if err != nil {
			return nil, err
		}

		return -1 * right.(float64), nil
	}

	return nil, nil
}

func (interpreter *Interpreter) VisitBinaryExpr(expr *ast.BinaryExpr) (interface{}, error) {
	left, err := interpreter.evaluate(expr.Left)
	if err != nil {
		return nil, err
	}

	right, err := interpreter.evaluate(expr.Right)
	if err != nil {
		return nil, err
	}

	switch expr.Operator.Type {
	case token.GREATER:
		err := checkNumberOperands(expr.Operator, left, right)
		if err != nil {
			return nil, err
		}
		return left.(float64) > right.(float64), nil
	case token.GREATER_EQUAL:
		err := checkNumberOperands(expr.Operator, left, right)
		if err != nil {
			return nil, err
		}
		return left.(float64) >= right.(float64), nil
	case token.LESS:
		err := checkNumberOperands(expr.Operator, left, right)
		if err != nil {
			return nil, err
		}
		return left.(float64) < right.(float64), nil
	case token.LESS_EQUAL:
		err := checkNumberOperands(expr.Operator, left, right)
		if err != nil {
			return nil, err
		}
		return left.(float64) <= right.(float64), nil
	case token.BANG_EQUAL:
		return isEqual(left, right), nil
	case token.EQUAL_EQUAL:
		return isEqual(left, right), nil
	case token.MINUS:
		err := checkNumberOperands(expr.Operator, left, right)
		if err != nil {
			return nil, err
		}
		return left.(float64) - right.(float64), nil
	case token.PLUS:
		leftNumValue, isLeftNumber := left.(float64)
		rightNumValue, isRightNumber := right.(float64)
		if isLeftNumber && isRightNumber {
			return leftNumValue + rightNumValue, nil
		}

		leftStringValue, isLeftString := left.(string)
		rightStringValue, isRightString := right.(string)
		if isLeftString && isRightString {
			return leftStringValue + rightStringValue, nil
		}

		return nil, loxerror.NewRuntimeError(expr.Operator, "Operands must be two numbers or two strings.")
	case token.SLASH:
		err := checkNumberOperands(expr.Operator, left, right)
		if err != nil {
			return nil, err
		}
		return left.(float64) / right.(float64), nil
	case token.STAR:
		err := checkNumberOperands(expr.Operator, left, right)
		if err != nil {
			return nil, err
		}
		return left.(float64) * right.(float64), nil
	}

	// Unreachable
	return nil, nil
}

func (interpreter *Interpreter) VisitExprStmt(stmt *ast.ExprStmt) error {
	_, err := interpreter.evaluate(stmt.Expression)
	return err
}

func (interpreter *Interpreter) VisitPrintStmt(stmt *ast.PrintStmt) error {
	value, err := interpreter.evaluate(stmt.Expression)
	if err != nil {
		return err
	}
	fmt.Println(stringify(value))
	return nil
}

func (interpreter *Interpreter) evaluate(expr ast.Expr) (interface{}, error) {
	return expr.Accept(interpreter)
}

func (interpreter *Interpreter) execute(stmt ast.Stmt) error {
	return stmt.Accept(interpreter)
}

func isTruthy(object interface{}) bool {
	if object == nil {
		return false
	}
	if value, isBool := object.(bool); isBool {
		return value
	} else {
		return true
	}
}

func isEqual(a interface{}, b interface{}) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil {
		return false
	}

	return reflect.DeepEqual(a, b)
}

func checkNumberOperand(operator token.Token, operand interface{}) error {
	if _, isNumber := operand.(float64); isNumber {
		return nil
	}

	return loxerror.NewRuntimeError(operator, "Operand must be a number.")
}

func checkNumberOperands(operator token.Token, left interface{}, right interface{}) error {
	_, isLeftNumber := left.(float64)
	_, isRightNumber := right.(float64)

	if isLeftNumber && isRightNumber {
		return nil
	}

	return loxerror.NewRuntimeError(operator, "Operands must be numbers.")
}

func stringify(object interface{}) string {
	if object == nil {
		return "nil"
	}

	if number, isNumber := object.(float64); isNumber {
		return strconv.FormatFloat(number, 'f', -1, 64)
	}

	return fmt.Sprintf("%v", object)
}
