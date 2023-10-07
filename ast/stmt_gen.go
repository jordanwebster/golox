package ast

type StmtVisitor interface {
	VisitExprStmt(stmt *ExprStmt) error
	VisitPrintStmt(stmt *PrintStmt) error
	VisitVarStmt(stmt *VarStmt) error
	VisitBlockStmt(stmt *BlockStmt) error
	VisitIfStmt(stmt *IfStmt) error
	VisitWhileStmt(stmt *WhileStmt) error
	VisitFunctionStmt(stmt *FunctionStmt) error
	VisitReturnStmt(stmt *ReturnStmt) error
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
func (stmt *IfStmt) Accept(visitor StmtVisitor) error {
	return visitor.VisitIfStmt(stmt)
}
func (stmt *WhileStmt) Accept(visitor StmtVisitor) error {
	return visitor.VisitWhileStmt(stmt)
}
func (stmt *FunctionStmt) Accept(visitor StmtVisitor) error {
	return visitor.VisitFunctionStmt(stmt)
}
func (stmt *ReturnStmt) Accept(visitor StmtVisitor) error {
	return visitor.VisitReturnStmt(stmt)
}
