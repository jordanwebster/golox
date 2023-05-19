package ast

type Stmt interface {
	Accept(visitor StmtVisitor) error
}

type ExprStmt struct {
	Expression Expr
}

type PrintStmt struct {
	Expression Expr
}
