package balena

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	parser "github.com/mrdapoyo/dofi/balena/parser"
	scanner "github.com/mrdapoyo/dofi/balena/scanner"
)

func runScript(t *testing.T, interp *Interpreter, src string) string {
	t.Helper()

	origStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	defer func() {
		_ = w.Close()
		os.Stdout = origStdout
	}()

	defer func() {
		if rcv := recover(); rcv != nil {
			t.Fatalf("script panicked: %v", rcv)
		}
	}()

	s := scanner.NewScanner(src)
	tokens := s.ScanTokens()

	p := parser.NewParser(tokens)
	stmts := p.Parse()

	resolver := NewResolver(interp)
	resolver.Resolve(stmts)

	interp.Execute(stmts)

	_ = w.Close()
	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)
	return strings.TrimSpace(buf.String())
}

func TestSimpleVarPrint(t *testing.T) {
	interp := NewInterpreter()
	out := runScript(t, interp, "var x = 1;\nprint x;")
	if out != "1" {
		t.Fatalf("expected output 1, got %q", out)
	}
}

func TestVariablePersistence(t *testing.T) {
	interp := NewInterpreter()
	_ = runScript(t, interp, "var x = 1;")
	out := runScript(t, interp, "x = x + 1;\nprint x;")
	if out != "2" {
		t.Fatalf("expected output 2, got %q", out)
	}
}

func TestFunctionLoop(t *testing.T) {
	script := `fn count(n) {
  while (n < 3) {
    print n;
    n = n + 1;
  }
}
count(1);`

	interp := NewInterpreter()
	out := runScript(t, interp, script)
	if out != "1\n2" && out != "1\n2\n" {
		t.Fatalf("unexpected output: %q", out)
	}
}

func TestBuiltinAdd(t *testing.T) {
	RegisterBuiltin("add", 2, func(args ...interface{}) interface{} {
		return args[0].(float64) + args[1].(float64)
	})

	interp := NewInterpreter()
	out := runScript(t, interp, "print add(1, 2);")
	if out != "3" {
		t.Fatalf("expected 3, got %q", out)
	}
}
