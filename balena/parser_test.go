package main

import (
	"fmt"
	"testing"
)

func TestLetStatements(t *testing.T) {
	tests := []struct {
		input              string
		expectedIdentifier string
		expectedValue      interface{}
	}{
		{"let x = 5;", "x", 5},
		{"let y = true;", "y", true},
		{"let foobar = y;", "foobar", "y"},
	}
	for _, tt := range tests {
		l := NewLexer(tt.input)
		p := NewParser(l)
		program := p.ParseProgram()
		CheckParserErrors(t, p)
		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain 1 statements. got=%d",
				len(program.Statements))
		}
		stmt := program.Statements[0]
		if !testLetStatement(t, stmt, tt.expectedIdentifier) {
			return
		}
		val := stmt.(*LetStatement).Value
		if !testLiteralExpression(t, val, tt.expectedValue) {
			return
		}
	}
}

func testLetStatement(t *testing.T, s Statement, name string) bool {
	if s.TokenLiteral() != "let" {
		t.Errorf("s.TokenLiteral not 'let'. got=%q", s.TokenLiteral())
		return false
	}
	letStmt, ok := s.(*LetStatement)
	if !ok {
		t.Errorf("s not *LetStatement. got=%T", s)
		return false
	}
	if letStmt.Name.Value != name {
		t.Errorf("letStmt.Name.Value not '%s'. got=%s", name, letStmt.Name.Value)
		return false
	}
	if letStmt.Name.TokenLiteral() != name {
		t.Errorf("s.Name not '%s'. got=%s", name, letStmt.Name)
		return false
	}
	return true
}

func TestReturnStatements(t *testing.T) {
	input := `
return 5;
return 10;
return 993322;
`
	l := NewLexer(input)
	p := NewParser(l)
	program := p.ParseProgram()
	CheckParserErrors(t, p)
	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements. got=%d",
			len(program.Statements))
	}
	for _, stmt := range program.Statements {
		returnStmt, ok := stmt.(*ReturnStatement)
		if !ok {
			t.Errorf("stmt not *ast.returnStatement. got=%T", stmt)
			continue
		}
		if returnStmt.TokenLiteral() != "return" {
			t.Errorf("returnStmt.TokenLiteral not 'return', got %q",
				returnStmt.TokenLiteral())
		}
	}
}

func TestIdentifierExpression(t *testing.T) {
	input := "foobar;"
	l := NewLexer(input)
	p := NewParser(l)
	program := p.ParseProgram()
	CheckParserErrors(t, p)
	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d",
			len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}
	ident, ok := stmt.Expression.(*Identifier)
	if !ok {
		t.Fatalf("exp not *Identifier. got=%T", stmt.Expression)
	}
	if ident.Value != "foobar" {
		t.Errorf("ident.Value not %s. got=%s", "foobar", ident.Value)
	}
	if ident.TokenLiteral() != "foobar" {
		t.Errorf("ident.TokenLiteral not %s. got=%s", "foobar",
			ident.TokenLiteral())
	}
}

func TestIntegerLiteralExpression(t *testing.T) {
	input := "5;"
	l := NewLexer(input)
	p := NewParser(l)
	program := p.ParseProgram()
	CheckParserErrors(t, p)
	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d",
			len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}
	literal, ok := stmt.Expression.(*IntegerLiteral)
	if !ok {
		t.Fatalf("exp not *IntegerLiteral. got=%T", stmt.Expression)
	}
	if literal.Value != 5 {
		t.Errorf("literal.Value not %d. got=%d", 5, literal.Value)
	}
	if literal.TokenLiteral() != "5" {
		t.Errorf("literal.TokenLiteral not %s. got=%s", "5",
			literal.TokenLiteral())
	}
}

func TestParsingPrefixExpressions(t *testing.T) {
	prefixTests := []struct {
		input        string
		operator     string
		integerValue int64
	}{
		{"!5;", "!", 5},
		{"-15;", "-", 15},
	}
	for _, tt := range prefixTests {
		l := NewLexer(tt.input)
		p := NewParser(l)
		program := p.ParseProgram()
		CheckParserErrors(t, p)
		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
				1, len(program.Statements))
		}
		stmt, ok := program.Statements[0].(*ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ExpressionStatement. got=%T",
				program.Statements[0])
		}
		exp, ok := stmt.Expression.(*PrefixExpression)
		if !ok {
			t.Fatalf("stmt is not PrefixExpression. got=%T", stmt.Expression)
		}
		if exp.Operator != tt.operator {
			t.Fatalf("exp.Operator is not '%s'. got=%s",
				tt.operator, exp.Operator)
		}
		if !testIntegerLiteral(t, exp.Right, tt.integerValue) {
			return
		}
	}
}

func testIntegerLiteral(t *testing.T, il Expression, value int64) bool {
	integ, ok := il.(*IntegerLiteral)
	if !ok {
		t.Errorf("il not *IntegerLiteral. got=%T", il)
		return false
	}
	if integ.Value != value {
		t.Errorf("integ.Value not %d. got=%d", value, integ.Value)
		return false
	}
	if integ.TokenLiteral() != fmt.Sprintf("%d", value) {
		t.Errorf("integ.TokenLiteral not %d. got=%s", value,
			integ.TokenLiteral())
		return false
	}
	return true
}

func TestParsingInfixExpressions(t *testing.T) {
	infixTests := []struct {
		input      string
		leftValue  int64
		operator   string
		rightValue int64
	}{
		{"5 + 5;", 5, "+", 5},
		{"5 - 5;", 5, "-", 5},
		{"5 * 5;", 5, "*", 5},
		{"5 / 5;", 5, "/", 5},
		{"5 > 5;", 5, ">", 5},
		{"5 < 5;", 5, "<", 5},
		{"5 == 5;", 5, "==", 5},
		{"5 != 5;", 5, "!=", 5},
	}
	for _, tt := range infixTests {
		l := NewLexer(tt.input)
		p := NewParser(l)
		program := p.ParseProgram()
		CheckParserErrors(t, p)
		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
				1, len(program.Statements))
		}
		stmt, ok := program.Statements[0].(*ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ExpressionStatement. got=%T",
				program.Statements[0])
		}
		exp, ok := stmt.Expression.(*InfixExpression)
		if !ok {
			t.Fatalf("exp is not InfixExpression. got=%T", stmt.Expression)
		}
		if !testIntegerLiteral(t, exp.Left, tt.leftValue) {
			return
		}
		if exp.Operator != tt.operator {
			t.Fatalf("exp.Operator is not '%s'. got=%s",
				tt.operator, exp.Operator)
		}
		if !testIntegerLiteral(t, exp.Right, tt.rightValue) {
			return
		}
	}
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"-a * b",
			"((-a) * b)",
		},
		{
			"!-a",
			"(!(-a))",
		},
		{
			"a + b + c",
			"((a + b) + c)",
		},
		{
			"a + b - c",
			"((a + b) - c)",
		},
		{
			"a * b * c",
			"((a * b) * c)",
		},
		{
			"a * b / c",
			"((a * b) / c)",
		},
		{
			"a + b / c",
			"(a + (b / c))",
		},
		{
			"a + b * c + d / e - f",
			"(((a + (b * c)) + (d / e)) - f)",
		},
		{
			"3 + 4; -5 * 5",
			"(3 + 4)((-5) * 5)",
		},
		{
			"5 > 4 == 3 < 4",
			"((5 > 4) == (3 < 4))",
		},
		{
			"5 < 4 != 3 > 4",
			"((5 < 4) != (3 > 4))",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
		{
			"1 + (2 + 3) + 4",
			"((1 + (2 + 3)) + 4)",
		},
		{
			"(5 + 5) * 2",
			"((5 + 5) * 2)",
		},
		{
			"2 / (5 + 5)",
			"(2 / (5 + 5))",
		},
		{
			"-(5 + 5)",
			"(-(5 + 5))",
		},
		{
			"!(true == true)",
			"(!(true == true))",
		},
		{
			"a + add(b * c) + d",
			"((a + add((b * c))) + d)",
		},
		{
			"add(a, b, 1, 2 * 3, 4 + 5, add(6, 7 * 8))",
			"add(a, b, 1, (2 * 3), (4 + 5), add(6, (7 * 8)))",
		},
		{
			"add(a + b + c * d / f + g)",
			"add((((a + b) + ((c * d) / f)) + g))",
		},
	}
	for _, tt := range tests {
		l := NewLexer(tt.input)
		p := NewParser(l)
		program := p.ParseProgram()
		CheckParserErrors(t, p)
		actual := program.String()
		if actual != tt.expected {
			t.Errorf("expected=%q, got=%q", tt.expected, actual)
		}
	}
}

func TestIfExpression(t *testing.T) {
	input := `if (x < y) { x }`
	l := NewLexer(input)
	p := NewParser(l)
	program := p.ParseProgram()
	CheckParserErrors(t, p)
	if len(program.Statements) != 1 {
		t.Fatalf("program.Body does not contain %d statements. got=%d\n",
			1, len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ExpressionStatement. got=%T",
			program.Statements[0])
	}
	exp, ok := stmt.Expression.(*IfExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not IfExpression. got=%T",
			stmt.Expression)
	}
	if !testInfixExpression(t, exp.Condition, "x", "<", "y") {
		return
	}
	if len(exp.Consequence.Statements) != 1 {
		t.Errorf("consequence is not 1 statements. got=%d\n",
			len(exp.Consequence.Statements))
	}
	consequence, ok := exp.Consequence.Statements[0].(*ExpressionStatement)
	if !ok {
		t.Fatalf("Statements[0] is not ExpressionStatement. got=%T",
			exp.Consequence.Statements[0])
	}
	if !testIdentifier(t, consequence.Expression, "x") {
		return
	}
	if exp.Alternative != nil {
		t.Errorf("exp.Alternative.Statements was not nil. got=%+v", exp.Alternative)
	}
}

func testIdentifier(t *testing.T, exp Expression, value string) bool {
	ident, ok := exp.(*Identifier)
	if !ok {
		t.Errorf("exp not *Identifier. got=%T", exp)
		return false
	}
	if ident.Value != value {
		t.Errorf("ident.Value not %s. got=%s", value, ident.Value)
		return false
	}
	if ident.TokenLiteral() != value {
		t.Errorf("ident.TokenLiteral not %s. got=%s", value,
			ident.TokenLiteral())
		return false
	}
	return true
}

func testInfixExpression(t *testing.T, exp Expression, left interface{},
	operator string, right interface{}) bool {
	opExp, ok := exp.(*InfixExpression)
	if !ok {
		t.Errorf("exp is not InfixExpression. got=%T(%s)", exp, exp)
		return false
	}
	if !testLiteralExpression(t, opExp.Left, left) {
		return false
	}
	if opExp.Operator != operator {
		t.Errorf("exp.Operator is not '%s'. got=%q", operator, opExp.Operator)
		return false
	}
	if !testLiteralExpression(t, opExp.Right, right) {
		return false
	}
	return true
}

func testLiteralExpression(
	t *testing.T,
	exp Expression,
	expected interface{},
) bool {
	switch v := expected.(type) {
	case int:
		return testIntegerLiteral(t, exp, int64(v))
	case int64:
		return testIntegerLiteral(t, exp, v)
	case string:
		return testIdentifier(t, exp, v)
	case bool:
		return testBooleanLiteral(t, exp, v)
	}
	t.Errorf("type of exp not handled. got=%T", exp)
	return false
}

func TestFunctionLiteralParsing(t *testing.T) {
	input := `fn(x, y) { x + y; }`
	l := NewLexer(input)
	p := NewParser(l)
	program := p.ParseProgram()
	CheckParserErrors(t, p)
	if len(program.Statements) != 1 {
		t.Fatalf("program.Body does not contain %d statements. got=%d\n",
			1, len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ExpressionStatement. got=%T",
			program.Statements[0])
	}
	function, ok := stmt.Expression.(*FunctionLiteral)
	if !ok {
		t.Fatalf("stmt.Expression is not FunctionLiteral. got=%T",
			stmt.Expression)
	}
	if len(function.Parameters) != 2 {
		t.Fatalf("function literal parameters wrong. want 2, got=%d\n",
			len(function.Parameters))
	}
	testLiteralExpression(t, function.Parameters[0], "x")
	testLiteralExpression(t, function.Parameters[1], "y")
	if len(function.Body.Statements) != 1 {
		t.Fatalf("function.Body.Statements has not 1 statements. got=%d\n",
			len(function.Body.Statements))
	}
	bodyStmt, ok := function.Body.Statements[0].(*ExpressionStatement)
	if !ok {
		t.Fatalf("function body stmt is not ExpressionStatement. got=%T",
			function.Body.Statements[0])
	}
	testInfixExpression(t, bodyStmt.Expression, "x", "+", "y")
}

func TestFunctionParameterParsing(t *testing.T) {
	tests := []struct {
		input          string
		expectedParams []string
	}{
		{input: "fn() {};", expectedParams: []string{}},
		{input: "fn(x) {};", expectedParams: []string{"x"}},
		{input: "fn(x, y, z) {};", expectedParams: []string{"x", "y", "z"}},
	}
	for _, tt := range tests {
		l := NewLexer(tt.input)
		p := NewParser(l)
		program := p.ParseProgram()
		CheckParserErrors(t, p)
		stmt := program.Statements[0].(*ExpressionStatement)
		function := stmt.Expression.(*FunctionLiteral)
		if len(function.Parameters) != len(tt.expectedParams) {
			t.Errorf("length parameters wrong. want %d, got=%d\n",
				len(tt.expectedParams), len(function.Parameters))
		}
		for i, ident := range tt.expectedParams {
			testLiteralExpression(t, function.Parameters[i], ident)
		}
	}
}

func testBooleanLiteral(t *testing.T, exp Expression, value bool) bool {
	bo, ok := exp.(*Boolean)
	if !ok {
		t.Errorf("exp not *Boolean. got=%T", exp)
		return false
	}
	if bo.Value != value {
		t.Errorf("bo.Value not %t. got=%t", value, bo.Value)
		return false
	}
	if bo.TokenLiteral() != fmt.Sprintf("%t", value) {
		t.Errorf("bo.TokenLiteral not %t. got=%s", value,
			bo.TokenLiteral())
		return false
	}
	return true
}

func TestCallExpressionParsing(t *testing.T) {
	input := "add(1, 2 * 3, 4 + 5);"
	l := NewLexer(input)
	p := NewParser(l)
	program := p.ParseProgram()
	CheckParserErrors(t, p)
	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
			1, len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ExpressionStatement)
	if !ok {
		t.Fatalf("stmt is not ExpressionStatement. got=%T",
			program.Statements[0])
	}
	exp, ok := stmt.Expression.(*CallExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not CallExpression. got=%T",
			stmt.Expression)
	}
	if !testIdentifier(t, exp.Function, "add") {
		return
	}
	if len(exp.Arguments) != 3 {
		t.Fatalf("wrong length of arguments. got=%d", len(exp.Arguments))
	}
	testLiteralExpression(t, exp.Arguments[0], 1)
	testInfixExpression(t, exp.Arguments[1], 2, "*", 3)
	testInfixExpression(t, exp.Arguments[2], 4, "+", 5)
}

func TestStringLiteralExpression(t *testing.T) {
	input := `"hello world";`
	l := NewLexer(input)
	p := NewParser(l)
	program := p.ParseProgram()
	CheckParserErrors(t, p)
	stmt := program.Statements[0].(*ExpressionStatement)
	literal, ok := stmt.Expression.(*StringLiteral)
	if !ok {
		t.Fatalf("exp not *StringLiteral. got=%T", stmt.Expression)
	}
	if literal.Value != "hello world" {
		t.Errorf("literal.Value not %q. got=%q", "hello world", literal.Value)
	}
}
