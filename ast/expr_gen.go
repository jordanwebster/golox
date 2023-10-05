package ast

type ExprVisitor interface {
	VisitBinaryExpr(expr *BinaryExpr) (interface{}, error)
	VisitGroupingExpr(expr *GroupingExpr) (interface{}, error)
	VisitLiteralExpr(expr *LiteralExpr) (interface{}, error)
	VisitUnaryExpr(expr *UnaryExpr) (interface{}, error)
	VisitVariableExpr(expr *VariableExpr) (interface{}, error)
	VisitAssignExpr(expr *AssignExpr) (interface{}, error)
	VisitLogicalExpr(expr *LogicalExpr) (interface{}, error)
	VisitCallExpr(expr *CallExpr) (interface{}, error)
}

func (expr *BinaryExpr) Accept(visitor ExprVisitor) (interface{}, error) {
	return visitor.VisitBinaryExpr(expr)
}
func (expr *GroupingExpr) Accept(visitor ExprVisitor) (interface{}, error) {
	return visitor.VisitGroupingExpr(expr)
}
func (expr *LiteralExpr) Accept(visitor ExprVisitor) (interface{}, error) {
	return visitor.VisitLiteralExpr(expr)
}
func (expr *UnaryExpr) Accept(visitor ExprVisitor) (interface{}, error) {
	return visitor.VisitUnaryExpr(expr)
}
func (expr *VariableExpr) Accept(visitor ExprVisitor) (interface{}, error) {
	return visitor.VisitVariableExpr(expr)
}
func (expr *AssignExpr) Accept(visitor ExprVisitor) (interface{}, error) {
	return visitor.VisitAssignExpr(expr)
}
func (expr *LogicalExpr) Accept(visitor ExprVisitor) (interface{}, error) {
	return visitor.VisitLogicalExpr(expr)
}
func (expr *CallExpr) Accept(visitor ExprVisitor) (interface{}, error) {
	return visitor.VisitCallExpr(expr)
}
