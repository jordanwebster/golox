package ast

import "github.com/jordanwebster/golox/token"

type Expr interface {
	Accept(visitor ExprVisitor) (interface{}, error)
}

type BinaryExpr struct {
	Operator token.Token
	Left     Expr
	Right    Expr
}

type GroupingExpr struct {
	Expression Expr
}

type LiteralExpr struct {
	Value interface{}
}

type UnaryExpr struct {
	Operator token.Token
	Right    Expr
}
