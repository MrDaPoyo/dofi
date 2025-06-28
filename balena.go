package main

import (
	"fmt"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"

	balena "github.com/mrdapoyo/dofi/balena/interpreter"
	parser "github.com/mrdapoyo/dofi/balena/parser"
	scanner "github.com/mrdapoyo/dofi/balena/scanner"
)

// Global error state variables
var hadError = false
var hadRuntimeError = false

func (g *Game) setupBalenaAPI() {
	g.BalenaEnv.Globals.Define("clear", func(_ ...interface{}) interface{} {
		g.LinearBuffer = nil
		g.Input.CurrentInputString = ""
		return nil
	})

	// pset(x, y, r, g,b): set a pixel in the 128×128 back-buffer
	g.BalenaEnv.Globals.Define("pset", func(a ...interface{}) interface{} {
		if len(a) != 5 {
			return nil
		}
		x, y := int(a[0].(float64)), int(a[1].(float64))
		clr := color.RGBA{uint8(a[2].(float64)), uint8(a[3].(float64)), uint8(a[4].(float64)), 255}
		if x >= 0 && x < 128 && y >= 0 && y < 128 {
			g.Screen.Buffer[y][x] = clr
		}
		return nil
	})

	g.BalenaEnv.Globals.Define("cls", func(_ ...interface{}) interface{} {
		g.ClearScreenBuffer()
		return nil
	})

	g.BalenaEnv.Globals.Define("print", func(a ...interface{}) interface{} {
		for _, arg := range a {
			g.AppendLine(arg.(string), false)
		}
		return nil
	})

	g.BalenaEnv.Globals.Define("txt", g.bindingDrawText)

	g.BalenaEnv.Globals.Define("sin", func(a ...interface{}) interface{} {
		if len(a) != 1 {
			return nil
		}
		return math.Sin(a[0].(float64))
	})

	g.BalenaEnv.Globals.Define("cos", func(a ...interface{}) interface{} {
		if len(a) != 1 {
			return nil
		}
		return math.Cos(a[0].(float64))
	})

	g.BalenaEnv.Globals.Define("floor", func(a ...interface{}) interface{} {
		if len(a) != 1 {
			return nil
		}
		return math.Floor(a[0].(float64))
	})

	g.BalenaEnv.Globals.Define("sqrt", func(a ...interface{}) interface{} {
		if len(a) != 1 {
			return nil
		}
		return math.Sqrt(a[0].(float64))
	})

	g.BalenaEnv.Globals.Define("time", 0.0)
	g.BalenaEnv.Globals.Define("screen_width", 128.0)
	g.BalenaEnv.Globals.Define("screen_height", 128.0)
	g.BalenaEnv.Globals.Define("cx", 64.0)
	g.BalenaEnv.Globals.Define("cy", 64.0)
	g.BalenaEnv.Globals.Define("R1", 20.0)
	g.BalenaEnv.Globals.Define("R2", 40.0)
	g.BalenaEnv.Globals.Define("K2", 200.0)
	g.BalenaEnv.Globals.Define("K1", 128.0*200.0*3.0/(8.0*(20.0+40.0)))
	g.BalenaEnv.Globals.Define("A_speed", 0.07)
	g.BalenaEnv.Globals.Define("B_speed", 0.03)
}

func (g *Game) RunBalenaScript(code string) {
	g.ClearScreenBuffer()

	// Reset error state
	hadError = false
	hadRuntimeError = false

	// Scan tokens
	tokens := scanner.NewScanner(code).ScanTokens()

	// Parse statements with error handling
	defer func() {
		if r := recover(); r != nil {
			if parseErr, ok := r.(parser.ParseError); ok {
				g.AppendLine(fmt.Sprintf("Parse error at line %d: %s", parseErr.Line, parseErr.Message), false)
			} else {
				g.AppendLine(fmt.Sprintf("Unexpected error: %v", r), false)
			}
		}
	}()

	raw := parser.NewParser(tokens).Parse()
	if hadError {
		return
	}

	stmts := make([]parser.Stmt, 0, len(raw))
	for _, s := range raw {
		if s != nil {
			stmts = append(stmts, s)
		}
	}
	if len(stmts) == 0 {
		return
	}

	// Resolve statements with error handling
	resolver := balena.NewResolver(g.BalenaEnv)
	defer func() {
		if r := recover(); r != nil {
			if runtimeErr, ok := r.(*balena.RuntimeError); ok {
				g.AppendLine(fmt.Sprintf("Resolution error: %s", runtimeErr.Error()), false)
			} else {
				g.AppendLine(fmt.Sprintf("Unexpected resolution error: %v", r), false)
			}
		}
	}()

	resolver.Resolve(stmts)
	if hadError {
		return
	}

	// Execute statements with error handling
	defer func() {
		if r := recover(); r != nil {
			if runtimeErr, ok := r.(*balena.RuntimeError); ok {
				g.AppendLine(fmt.Sprintf("Runtime error: %s", runtimeErr.Error()), false)
			} else if returnErr, ok := r.(balena.Return); ok {
				if returnErr.Value != nil {
					g.AppendLine(fmt.Sprintf("Return value: %v", returnErr.Value), false)
				}
				g.AppendLine("Unexpected return statement at top level.", false)
			} else {
				g.AppendLine(fmt.Sprintf("Unexpected execution error: %v", r), false)
			}
		}
	}()

	g.BalenaEnv.Execute(stmts)
}

func (g *Game) bindingDrawText(a ...interface{}) interface{} {
	if len(a) != 4 {
		return nil
	}
	x := int(a[0].(float64))
	y := int(a[1].(float64))
	s := a[2].(string)
	shade := uint8(a[3].(float64))

	col := color.RGBA{shade, shade, shade, 255}
	op := &text.DrawOptions{}
	op.ColorScale.Scale(float32(col.R)/255, float32(col.G)/255, float32(col.B)/255, 1)
	op.GeoM.Translate(float64(x), float64(y))

	img := ebiten.NewImage(g.Screen.Width, g.Screen.Height)
	text.Draw(img, s, TextFace, op)

	buf := make([]byte, 4*g.Screen.Width*g.Screen.Height)
	img.ReadPixels(buf)
	for i := 0; i < len(buf); i += 4 {
		r, gC, bC, aC := buf[i], buf[i+1], buf[i+2], buf[i+3]
		yy := y + (i/4)/128
		xx := x + (i/4)%128
		if xx >= 0 && xx < 128 && yy >= 0 && yy < 128 {
			g.Screen.Buffer[yy][xx] = color.RGBA{r, gC, bC, aC}
		}
	}
	return nil
}

func (g *Game) ClearScreenBuffer() {
	for y := 0; y < len(g.Screen.Buffer); y++ {
		for x := 0; x < len(g.Screen.Buffer[y]); x++ {
			g.Screen.Buffer[y][x] = color.RGBA{0, 0, 0, 0}
		}
	}
}
