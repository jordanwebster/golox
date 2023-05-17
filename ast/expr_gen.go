package ast

type ExprVisitor interface {
	VisitBinaryExpr(expr *BinaryExpr) interface{}
	VisitGroupingExpr(expr *GroupingExpr) interface{}
	VisitLiteralExpr(expr *LiteralExpr) interface{}
	VisitUnaryExpr(expr *UnaryExpr) interface{}
}

func (expr *BinaryExpr) Accept(visitor ExprVisitor) interface{} {
	return visitor.VisitBinaryExpr(expr)
}
func (expr *GroupingExpr) Accept(visitor ExprVisitor) interface{} {
	return visitor.VisitGroupingExpr(expr)
}
func (expr *LiteralExpr) Accept(visitor ExprVisitor) interface{} {
	return visitor.VisitLiteralExpr(expr)
}
func (expr *UnaryExpr) Accept(visitor ExprVisitor) interface{} {
	return visitor.VisitUnaryExpr(expr)
}
