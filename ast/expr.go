package ast

import "github.com/jordanwebster/golox/token"

type Expr interface {
	Accept(visitor ExprVisitor) interface{}
}

type BinaryExpr struct {
	operator token.Token
	left     Expr
	right    Expr
}

type GroupingExpr struct {
	expression Expr
}

type LiteralExpr struct {
	value interface{}
}

type UnaryExpr struct {
	operator token.Token
	right    Expr
}
