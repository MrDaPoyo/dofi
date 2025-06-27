package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"

	env "github.com/mrdapoyo/dofi/balena/env"
	parser "github.com/mrdapoyo/dofi/balena/parser"
)

func (g *Game) SetupBalenaAPI() *env.Environment {
	var env = env.NewEnvironment()

	env.SetUserData("game", env.NewGoObject(g))

	env.RegisterExternalBinding("clear", g.bindingClearLines)
	env.RegisterExternalBinding("pset", g.bindingDrawPixel)
	env.RegisterExternalBinding("txt", g.bindingDrawText)

	return env
}

func (g *Game) RunBalenaScript(env *env.Environment, code string) []string {
	p := parser.NewParser(code)
	program := p.ParseProgram()
	if len(p.Errors()) != 0 {
		errors := []string{}
		for _, err := range p.Errors() {
			errors = append(errors, err)
		}
		return errors
	}
	evaluated := balena.Eval(program, env)
	if evaluated != nil {
		return []string{evaluated.Inspect()}
	}
	return nil
}

func (g *Game) bindingClearLines(args ...balena.Object) balena.Object {
	g.LinearBuffer = []LinearBuffer{}
	g.Input.CurrentInputString = ""
	return balena.EVAL_NULL
}

func (g *Game) bindingDrawPixel(args ...balena.Object) balena.Object {
	if len(args) != 5 {
		return balena.NewError("wrong number of arguments. got=%d, want=5", len(args))
	}

	x, ok1 := args[0].(*balena.ObjectInteger)
	y, ok2 := args[1].(*balena.ObjectInteger)
	r, ok3 := args[2].(*balena.ObjectInteger)
	gc, ok4 := args[3].(*balena.ObjectInteger)
	b, ok5 := args[4].(*balena.ObjectInteger)
	if !ok1 || !ok2 || !ok3 || !ok4 || !ok5 {
		return balena.NewError("arguments 1-5 must be integers")
	}
	color := color.RGBA{uint8(r.Value), uint8(gc.Value), uint8(b.Value), 255}
	if x.Value >= 0 && x.Value < 128 && y.Value >= 0 && y.Value < 128 {
		g.Screen.Buffer[y.Value][x.Value] = color
		return balena.EVAL_NULL
	}
	g.AppendLine("Error: Pixel out of bounds", true)
	return balena.EVAL_NULL
}

func (g *Game) bindingDrawText(args ...balena.Object) balena.Object {
	if len(args) != 5 {
		return balena.NewError("wrong number of arguments. got=%d, want=5", len(args))
	}

	x, ok1 := args[1].(*balena.ObjectInteger)
	y, ok2 := args[2].(*balena.ObjectInteger)
	value, ok3 := args[3].(*balena.String)
	c, ok4 := args[4].(*balena.ObjectInteger)
	if !ok1 || !ok2 || !ok3 || !ok4 {
		return balena.NewError("arguments 2-5 must be correct types")
	}
	col := color.RGBA{uint8(c.Value), uint8(c.Value), uint8(c.Value), 255}
	var op = &text.DrawOptions{}
	op.ColorScale.Scale(float32(col.R)/255, float32(col.G)/255, float32(col.B)/255, float32(col.A)/255)
	op.GeoM.Translate(float64(x.Value), float64(y.Value))
	image := ebiten.NewImage(g.Screen.Width, g.Screen.Height)
	text.Draw(image, value.Value, TextFace, op)
	buffer := make([]byte, 4*g.Screen.Width*g.Screen.Height)
	image.ReadPixels(buffer)
	for i := 0; i < len(buffer); i += 4 {
		r := buffer[i]
		green := buffer[i+1]
		b := buffer[i+2]
		a := buffer[i+3]
		g.Screen.Buffer[int(y.Value)+i/4/128][int(x.Value)+i/4%128] = color.RGBA{r, green, b, a}
	}
	return balena.EVAL_NULL
}

