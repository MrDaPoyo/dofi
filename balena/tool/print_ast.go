package balena


type AstPrinter struct{}

func (a *AstPrinter) Print(expr Expr) string {
	return expr.Accept(a)
}

func (a *AstPrinter) VisitBinaryExpr(expr *Binary) string {
	return a.parenthesize(expr.Operator.Lexeme, expr.Left, expr.Right)
}

func (a *AstPrinter) VisitGroupingExpr(expr *Grouping) string {
	return a.parenthesize("group", expr.Expression)
}

func (a *AstPrinter) VisitLiteralExpr(expr *Literal) string {
	if expr.Value == nil {
		return "nil"
	}
	return expr.Value.(string)
}

func (a *AstPrinter) VisitUnaryExpr(expr *Unary) string {
	return a.parenthesize(expr.Operator.Lexeme, expr.Right)
}

func (a *AstPrinter) parenthesize(name string, exprs ...Expr) string {
	result := "(" + name
	for _, expr := range exprs {
		result += " " + expr.Accept(a)
	}
	result += ")"
	return result
}