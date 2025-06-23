package balena
// 
// import (
// 	"testing"
// 	"fmt"
// )
// func TestLexerFloatTokens(t *testing.T) {
// input := "2.5 * 4.0"
// 
// 	lexer := NewLexer(input)
// 	
// 	fmt.Println("=== LEXER DEBUG ===")
// 	for {
// 		tok := lexer.NextToken()
// 		fmt.Printf("Token: Type=%s, Literal=%s\n", tok.Type, tok.Literal)
// 		if tok.Type == EOF {
// 			break
// 		}
// 	}
// }
// func TestFloatMultiplicationDebug(t *testing.T) {
// input := "2.5 * 4.0"
// 
// 	
// 	fmt.Println("\n=== STEP 1: LEXER TEST ===")
// 	lexer := NewLexer(input)
// 	tokens := []Token{}
// 	for {
// 		tok := lexer.NextToken()
// 		tokens = append(tokens, tok)
// 		fmt.Printf("Token: %s = %s\n", tok.Type, tok.Literal)
// 		if tok.Type == EOF {
// 			break
// 		}
// 	}
// 	
// 	fmt.Println("\n=== STEP 2: PARSER TEST ===")
// 	lexer2 := NewLexer(input)
// 	parser := NewParser(lexer2)
// 	program := parser.ParseProgram()
// 	
// 	if len(parser.Errors()) > 0 {
// 		fmt.Println("Parser errors:")
// 		for _, err := range parser.Errors() {
// 			fmt.Printf("  %s\n", err)
// 		}
// 	} else {
// 		fmt.Println("Parsed program successfully!")
// 		for _, stmt := range program.Statements {
// 			fmt.Printf("Statement: %T\n", stmt)
// 			if exprStmt, ok := stmt.(*ExpressionStatement); ok {
// 				fmt.Printf("  Expression: %T\n", exprStmt.Expression)
// 			}
// 		}
// 	}
// 	
// 	fmt.Println("\n=== STEP 3: EVALUATION TEST ===")
// 	env := NewEnvironment()
// 	result := Eval(program, env)
// 	fmt.Printf("Result type: %T\n", result)
// 	fmt.Printf("Result: %+v\n", result)
// 	
// 	if floatObj, ok := result.(*ObjectFloat); ok {
// 		fmt.Printf("Float value: %f\n", floatObj.Value)
// 	} else if intObj, ok := result.(*ObjectInteger); ok {
// 		fmt.Printf("Integer value: %d (THIS IS THE PROBLEM!)\n", intObj.Value)
// 	}
// }