package ast

type StmtVisitor interface {
	VisitExprStmt(stmt *ExprStmt) (interface{}, error)
	VisitPrintStmt(stmt *PrintStmt) (interface{}, error)
}

func (stmt *ExprStmt) Accept(visitor StmtVisitor) (interface{}, error) {
	return visitor.VisitExprStmt(stmt)
}
func (stmt *PrintStmt) Accept(visitor StmtVisitor) (interface{}, error) {
	return visitor.VisitPrintStmt(stmt)
}
