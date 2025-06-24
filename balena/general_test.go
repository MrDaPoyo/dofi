package balena

import (
	"fmt"
	"log"
	"testing"
)
func TestLexerFloatTokens(t *testing.T) {
input := "2.5 * 4.0"

	lexer := NewLexer(input)
	
	fmt.Println("=== LEXER DEBUG ===")
	for {
		tok := lexer.NextToken()
		fmt.Printf("Token: Type=%s, Literal=%s\n", tok.Type, tok.Literal)
		if tok.Type == EOF {
			break
		}
	}
}
func TestFloatMultiplicationDebug(t *testing.T) {
input := "2.5 * 4.0"

	
	fmt.Println("\n=== STEP 1: LEXER TEST ===")
	lexer := NewLexer(input)
	tokens := []Token{}
	for {
		tok := lexer.NextToken()
		tokens = append(tokens, tok)
		fmt.Printf("Token: %s = %s\n", tok.Type, tok.Literal)
		if tok.Type == EOF {
			break
		}
	}
	
	fmt.Println("\n=== STEP 2: PARSER TEST ===")
	lexer2 := NewLexer(input)
	parser := NewParser(lexer2)
	program := parser.ParseProgram()
	
	if len(parser.Errors()) > 0 {
		fmt.Println("Parser errors:")
		for _, err := range parser.Errors() {
			fmt.Printf("  %s\n", err)
		}
	} else {
		fmt.Println("Parsed program successfully!")
		for _, stmt := range program.Statements {
			fmt.Printf("Statement: %T\n", stmt)
			if exprStmt, ok := stmt.(*ExpressionStatement); ok {
				fmt.Printf("  Expression: %T\n", exprStmt.Expression)
			}
		}
	}
	
	fmt.Println("\n=== STEP 3: EVALUATION TEST ===")
	env := NewEnvironment()
	result := Eval(program, env)
	fmt.Printf("Result type: %T\n", result)
	fmt.Printf("Result: %+v\n", result)
	
	if floatObj, ok := result.(*ObjectFloat); ok {
		fmt.Printf("Float value: %f\n", floatObj.Value)
	} else if intObj, ok := result.(*ObjectInteger); ok {
		fmt.Printf("Integer value: %d (THIS IS THE PROBLEM!)\n", intObj.Value)
	}
}

func TestForLoops(t *testing.T) {
	input := `
	let i = 0;
	for (let x = 0; x < 10; x += 1) {
		i * 2;
	}
	`
	lexer := NewLexer(input)
	fmt.Println("=== LEXER DEBUG ===")
	for {
		tok := lexer.NextToken()
		fmt.Printf("Token: Type=%s, Literal=%s\n", tok.Type, tok.Literal)
		if tok.Type == EOF {
			break
		}
	}
	fmt.Println("\n=== STEP 2: PARSER TEST ===")
	parser := NewParser(lexer)
	program := parser.ParseProgram()
	fmt.Printf("Parser output: %+v\n", program)

	fmt.Println(program)
	
	if len(parser.Errors()) > 0 {
		fmt.Println("Parser errors:")
		for _, err := range parser.Errors() {
			fmt.Printf("  %s\n", err)
		}
	} else {
		fmt.Println("Parsed program successfully!")
		for _, stmt := range program.Statements {
			fmt.Printf("Statement: %T\n", stmt)
			if exprStmt, ok := stmt.(*ExpressionStatement); ok {
				fmt.Printf("  Expression: %T\n", exprStmt.Expression)
			}
		}
	}
	fmt.Println("\n=== STEP 3: EVALUATION TEST ===")
	env := NewEnvironment()
	env.Set("i", &ObjectInteger{Value: 10})
	result := Eval(program, env)
	if result != nil {
		fmt.Printf("Result: %+v\n", result)
	}

	iVal, ok := env.Get("i")
	log.Print(iVal)
	if !ok {
		log.Printf("Modified variables: %v", env.modifiedVars)
		log.Printf("Environment store: %v", env.store)
		t.Errorf("Variable 'i' not found in environment")
	}

}

func TestWhileLoops(t *testing.T) {
	input := `
	let ibefore = 0;
	while (ibefore < 10) {
		ibefore += 1;
	}
	`
	lexer := NewLexer(input)
	fmt.Println("=== LEXER DEBUG ===")
	for {
		tok := lexer.NextToken()
		fmt.Printf("Token: Type=%s, Literal=%s\n", tok.Type, tok.Literal)
		if tok.Type == EOF {
			break
		}
	}
	fmt.Println("\n=== STEP 2: PARSER TEST ===")
	parser := NewParser(lexer)
	program := parser.ParseProgram()
	fmt.Printf("Parser output: %+v\n", program)
	if len(parser.Errors()) > 0 {
		fmt.Println("Parser errors:")
		for _, err := range parser.Errors() {
			fmt.Printf("  %s\n", err)
		}
	} else {
		fmt.Println("Parsed program successfully!")
		for _, stmt := range program.Statements {
			fmt.Printf("Statement: %T\n", stmt)
			if exprStmt, ok := stmt.(*ExpressionStatement); ok {
				fmt.Printf("  Expression: %T\n", exprStmt.Expression)
			}
		}
	}
	fmt.Println("\n=== STEP 3: EVALUATION TEST ===")
	env := NewEnvironment()
	result := Eval(program, env)
	log.Print(env.Get("i"))
	if result != nil {
		fmt.Printf("Result: %+v\n", result)
	} else {
		fmt.Println("Result is nil")
	}
	if _, ok := result.(*ObjectNull); ok {
		fmt.Println("Result is null")
	} else {
		fmt.Println("Result is not null")
	}
}