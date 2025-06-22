package main

import (
	"testing"
)

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"5", 5},
		{"10", 10},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}
func testEval(input string) Object {
	l := NewLexer(input)
	p := NewParser(l)
	program := p.ParseProgram()
	return Eval(program)
}
func testIntegerObject(t *testing.T, obj Object, expected int64) bool {
	result, ok := obj.(*ObjectInteger)
	if !ok {
		t.Errorf("object is not Integer. got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value. got=%d, want=%d",
			result.Value, expected)
		return false
	}
	return true
}
