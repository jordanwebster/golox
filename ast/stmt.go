package ast

type Stmt interface {
	Accept(visitor StmtVisitor) (interface{}, error)
}

type ExprStmt struct {
	Expression Expr
}

type PrintStmt struct {
	Expression Expr
}
