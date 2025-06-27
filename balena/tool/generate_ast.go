package balena

import (
	"fmt"
	"os"
	"strings"
)

type Ast struct{}

func (ast *Ast) Run(args []string) (string, error) {
	if len(args) != 2 {
		fmt.Fprintln(os.Stderr, "Usage: generate_ast <output directory>")
		os.Exit(64)
	}
	outputDir := args[1]
	err := ast.defineAst(outputDir, "Expr", []string{
		"Assign   : Token name, Expr value",
		"Binary   : Expr left, Token operator, Expr right",
		"Call     : Expr callee, Token paren, List<Expr> arguments",
		"Grouping : Expr expression",
		"Literal  : Object value",
		"Logical  : Expr left, Token operator, Expr right",
		"Unary    : Token operator, Expr right",
		"Variable : Token name",
	})

	err = ast.defineAst(outputDir, "Stmt", []string{
		"Block      : List<Stmt> statements",
		"Expression : Expr expression",
		"Function   : Token name, List<Token> params, List<Stmt> body",
		"If         : Expr condition, Stmt thenBranch, Stmt elseBranch",
		"Print      : Expr expression",
		"Return     : Token keyword, Expr value",
		"Var        : Token name, Expr initializer",
		"While      : Expr condition, Stmt body",
	})

	if err != nil {
		return "", err
	}
	return args[1], nil
}

func (ast *Ast) defineAst(outputDir, baseName string, types []string) error {
	path := fmt.Sprintf("%s/%s.go", outputDir, strings.ToLower(baseName))
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	fmt.Fprintf(file, "package balena\n\n")
	fmt.Fprintf(file, "// %s is the base interface for AST nodes.\n", baseName)
	fmt.Fprintf(file, "type %s interface {\n", baseName)
	fmt.Fprintf(file, "\tAccept(visitor %sVisitor) error\n", baseName)
	fmt.Fprintf(file, "}\n\n")

	defineVisitor(file, baseName, types)

	for _, t := range types {
		parts := strings.SplitN(t, ":", 2)
		className := strings.TrimSpace(parts[0])
		fields := ""
		if len(parts) > 1 {
			fields = strings.TrimSpace(parts[1])
		}
		defineType(file, baseName, className, fields)

		fmt.Fprintf(file, "func (e *%s) Accept(visitor %sVisitor) error {\n", className, baseName)
		fmt.Fprintf(file, "\treturn visitor.Visit%s(*e)\n", className)
		fmt.Fprintf(file, "}\n\n")
	}

	return nil
}

func defineType(file *os.File, baseName, className, fieldList string) {
	fmt.Fprintf(file, "type %s struct {\n", className)
	fields := []string{}
	if fieldList != "" {
		fields = strings.Split(fieldList, ", ")
		for _, field := range fields {
			parts := strings.SplitN(field, " ", 2)
			if len(parts) == 2 {
				fieldType := strings.TrimSpace(parts[0])
				fieldName := strings.TrimSpace(parts[1])
				fmt.Fprintf(file, "\t%s %s\n", fieldName, fieldType)
			}
		}
	}
	fmt.Fprintf(file, "}\n\n")
}

func defineVisitor(file *os.File, baseName string, types []string) {
	fmt.Fprintf(file, "type %sVisitor interface {\n", baseName)
	for _, t := range types {
		className := strings.TrimSpace(strings.SplitN(t, ":", 2)[0])
		fmt.Fprintf(file, "\tVisit%s(%s) error\n", className, className)
	}
	fmt.Fprintf(file, "}\n\n")
}
