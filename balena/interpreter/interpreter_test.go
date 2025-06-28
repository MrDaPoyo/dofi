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

func TestWhileLoop(t *testing.T) {
	script := `var i = 0;
while (i < 3) {
  print i;
  i = i + 1;
}`

	interp := NewInterpreter()
	out := runScript(t, interp, script)
	if out != "0\n1\n2" && out != "0\n1\n2\n" {
		t.Fatalf("unexpected output: %q", out)
	}
}

func TestForLoop(t *testing.T) {
	script := `for (var i = 0; i < 3; i = i + 1) {
  print i;
}`

	interp := NewInterpreter()
	out := runScript(t, interp, script)
	if out != "0\n1\n2" && out != "0\n1\n2\n" {
		t.Fatalf("unexpected output: %q", out)
	}
}

func TestForLoopScope(t *testing.T) {
	script := `for (var i = 0; i < 3; i = i + 1) {
	}
	print i;`
	interp := NewInterpreter()
	out := runScript(t, interp, script)
	if out != "" {
		t.Fatalf("unexpected output: %q", out)
	}
}

func TestNestedWhileLoopScope(t *testing.T) {
	script := `var i = 0;
while (i < 2) {
  var j = 0;
  while (j < 2) {
    print j;
    j = j + 1;
  }
  i = i + 1;
}`

	interp := NewInterpreter()
	out := runScript(t, interp, script)
	if out != "0\n1\n0\n1" && out != "0\n1\n0\n1\n" {
		t.Fatalf("unexpected output: %q", out)
	}
}

func TestDonutStyleLoop(t *testing.T) {
	script := `var theta = 0.0;
while (theta < 6.28) {
  var phi = 0.0;
  while (phi < 6.28) {
    print phi;
    phi = phi + 0.1;
  }
  theta = theta + 0.1;
}`

	interp := NewInterpreter()
	out := runScript(t, interp, script)
	// Should print phi values from 0.0 to 6.2 in steps of 0.1, repeated for each theta
	if len(out) == 0 {
		t.Fatalf("expected output, got empty string")
	}
}

func TestSimpleScope(t *testing.T) {
	script := `var x = 1;
{
  var y = 2;
  print y;
}
print x;`

	interp := NewInterpreter()
	out := runScript(t, interp, script)
	if out != "2\n1" && out != "2\n1\n" {
		t.Fatalf("unexpected output: %q", out)
	}
}

func TestWhileLoopWithBraces(t *testing.T) {
	script := `var i = 0;
while (i < 2) {
  var j = 0;
  while (j < 2) {
    print j;
    j = j + 1;
  }
  i = i + 1;
}`

	interp := NewInterpreter()
	out := runScript(t, interp, script)
	if out != "0\n1\n0\n1" && out != "0\n1\n0\n1\n" {
		t.Fatalf("unexpected output: %q", out)
	}
}

func TestWhileLoopVariableScope(t *testing.T) {
	script := `var i = 0;
while (i < 1) {
  var x = 42;
  print x;
  i = i + 1;
}
print i;`

	interp := NewInterpreter()
	out := runScript(t, interp, script)
	if out != "42\n1" && out != "42\n1\n" {
		t.Fatalf("unexpected output: %q", out)
	}
}

func TestSimpleForLoop(t *testing.T) {
	script := `for (var i = 0; i < 2; i = i + 1) {
  print i;
}`

	interp := NewInterpreter()
	out := runScript(t, interp, script)
	if out != "0\n1" && out != "0\n1\n" {
		t.Fatalf("unexpected output: %q", out)
	}
}

func TestForLoopFeatures(t *testing.T) {
	// Test 1: Basic for loop with variable declaration
	t.Run("BasicForLoop", func(t *testing.T) {
		script := `for (var i = 0; i < 3; i = i + 1) {
  print i;
}`
		interp := NewInterpreter()
		out := runScript(t, interp, script)
		if out != "0\n1\n2" && out != "0\n1\n2\n" {
			t.Fatalf("unexpected output: %q", out)
		}
	})

	// Test 2: For loop with expression initializer
	t.Run("ExpressionInitializer", func(t *testing.T) {
		script := `var x = 5;
for (x = 0; x < 2; x = x + 1) {
  print x;
}`
		interp := NewInterpreter()
		out := runScript(t, interp, script)
		if out != "0\n1" && out != "0\n1\n" {
			t.Fatalf("unexpected output: %q", out)
		}
	})

	// Test 3: For loop with empty initializer
	t.Run("EmptyInitializer", func(t *testing.T) {
		script := `var i = 0;
for (; i < 2; i = i + 1) {
  print i;
}`
		interp := NewInterpreter()
		out := runScript(t, interp, script)
		if out != "0\n1" && out != "0\n1\n" {
			t.Fatalf("unexpected output: %q", out)
		}
	})

	// Test 4: For loop with empty condition (infinite loop with break)
	t.Run("EmptyCondition", func(t *testing.T) {
		script := `var i = 0;
for (; ; i = i + 1) {
  print i;
  if (i >= 1) {
    break;
  }
}`
		interp := NewInterpreter()
		out := runScript(t, interp, script)
		if out != "0\n1" && out != "0\n1\n" {
			t.Fatalf("unexpected output: %q", out)
		}
	})

	// Test 5: For loop with empty increment
	t.Run("EmptyIncrement", func(t *testing.T) {
		script := `for (var i = 0; i < 2; ) {
  print i;
  i = i + 1;
}`
		interp := NewInterpreter()
		out := runScript(t, interp, script)
		if out != "0\n1" && out != "0\n1\n" {
			t.Fatalf("unexpected output: %q", out)
		}
	})
}

func TestForLoopSummary(t *testing.T) {
	// This test demonstrates all the working for loop features in Balena
	script := `// Basic for loop with variable declaration
for (var i = 0; i < 2; i = i + 1) {
  print i;
}

// For loop with expression initializer
var x = 5;
for (x = 0; x < 2; x = x + 1) {
  print x;
}

// For loop with empty initializer
var j = 0;
for (; j < 2; j = j + 1) {
  print j;
}

// For loop with empty increment
for (var k = 0; k < 2; ) {
  print k;
  k = k + 1;
}`

	interp := NewInterpreter()
	out := runScript(t, interp, script)
	expected := "0\n1\n0\n1\n0\n1\n0\n1"
	if out != expected && out != expected+"\n" {
		t.Fatalf("unexpected output: %q", out)
	}
}
