package ast

type StmtVisitor interface {
	VisitExprStmt(stmt *ExprStmt) error
	VisitPrintStmt(stmt *PrintStmt) error
}

func (stmt *ExprStmt) Accept(visitor StmtVisitor) error {
	return visitor.VisitExprStmt(stmt)
}
func (stmt *PrintStmt) Accept(visitor StmtVisitor) error {
	return visitor.VisitPrintStmt(stmt)
}
