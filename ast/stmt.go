package ast

import "github.com/jordanwebster/golox/token"

type Stmt interface {
	Accept(visitor StmtVisitor) error
}

type ExprStmt struct {
	Expression Expr
}

type PrintStmt struct {
	Expression Expr
}

type VarStmt struct {
	Name        token.Token
	Initializer Expr
}
