package ast

type StmtVisitor interface {
	VisitExprStmt(stmt *ExprStmt) error
	VisitPrintStmt(stmt *PrintStmt) error
	VisitVarStmt(stmt *VarStmt) error
	VisitBlockStmt(stmt *BlockStmt) error
}

func (stmt *ExprStmt) Accept(visitor StmtVisitor) error {
	return visitor.VisitExprStmt(stmt)
}
func (stmt *PrintStmt) Accept(visitor StmtVisitor) error {
	return visitor.VisitPrintStmt(stmt)
}
func (stmt *VarStmt) Accept(visitor StmtVisitor) error {
	return visitor.VisitVarStmt(stmt)
}
func (stmt *BlockStmt) Accept(visitor StmtVisitor) error {
	return visitor.VisitBlockStmt(stmt)
}
