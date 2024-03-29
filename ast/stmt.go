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

type BlockStmt struct {
	Statements []Stmt
}

type IfStmt struct {
	Condition  Expr
	ThenBranch Stmt
	ElseBranch Stmt
}

type WhileStmt struct {
	Condition Expr
	Body      Stmt
}

type FunctionStmt struct {
	Name       token.Token
	Parameters []token.Token
	Body       []Stmt
}

type ReturnStmt struct {
	Keyword token.Token
	Value   Expr
}
