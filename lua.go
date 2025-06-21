package main

import (
	"image/color"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"

	lua "github.com/yuin/gopher-lua"
)

func (g *Game) setupLuaAPI() {
	dofiTable := g.LuaVM.NewTable()
	g.LuaVM.SetGlobal("dofi", dofiTable)

	// cls function - clear screen
	g.LuaVM.SetGlobal("cls", g.LuaVM.NewFunction(func(L *lua.LState) int {
		g.Screen.Buffer = [128][128]color.RGBA{}
		g.ClearLines()
		return 0
	}))

	g.LuaVM.SetField(dofiTable, "cls", g.LuaVM.NewFunction(func(L *lua.LState) int {
		g.Screen.Buffer = [128][128]color.RGBA{}
		g.ClearLines()
		return 0
	}))

	// pset function - set pixel
	g.LuaVM.SetField(dofiTable, "pset", g.LuaVM.NewFunction(func(L *lua.LState) int {
		x := int(L.CheckNumber(1))
		y := int(L.CheckNumber(2))
		r := uint8(L.CheckNumber(3))
		green := uint8(L.CheckNumber(4))
		b := uint8(L.CheckNumber(5))

		g.DrawPixel(x, y, color.RGBA{r, green, b, 255})
		return 0
	}))

	g.LuaVM.SetGlobal("print", g.LuaVM.NewFunction(func(L *lua.LState) int {
		top := L.GetTop()
		if top == 0 {
			return 0
		}
		
		if L.GetTop() >= 1 && L.CheckAny(1).Type() == lua.LTString {
			var parts []string
			for i := 1; i <= top; i++ {
				parts = append(parts, L.ToString(i))
			}
			g.AppendLine(strings.Join(parts, " "), false)
			return 0
		}
		
		x := int(L.OptNumber(1, 0))
		y := int(L.OptNumber(2, 0))
		var parts []string
		for i := 3; i <= top; i++ {
			parts = append(parts, L.ToString(i))
		}
		g.DrawText(x, y, strings.Join(parts, " "), color.RGBA{255, 255, 255, 255})
		g.AppendLine(strings.Join(parts, " "), false)
		return 0
	}))
}

func (g *Game) RunLuaScript(script string) error {
	defer func() {
		g.ScriptRunning = false
	}()
	err := g.LuaVM.DoString(script)
	return err
}

func (g *Game) ClearLines() {
	g.LinearBuffer = []LinearBuffer{}
	g.Input.CurrentInputString = ""
}

func (g *Game) DrawPixel(x, y int, c color.RGBA) {
	if x >= 0 && x < 128 && y >= 0 && y < 128 {
		g.Screen.Buffer[y][x] = c
	}
	if x < 0 || x >= 128 || y < 0 || y >= 128 {
		g.AppendLine("Error: Pixel out of bounds", true)
		return
	}
	g.Screen.Buffer[y][x] = c
}

func (g *Game) DrawText(x, y int, value string, c color.RGBA) {
	var op = &text.DrawOptions{}
	op.ColorScale.Scale(float32(c.R)/255, float32(c.G)/255, float32(c.B)/255, float32(c.A)/255)
	op.GeoM.Translate(float64(x), float64(y))
	image := ebiten.NewImage(g.Screen.Width, g.Screen.Height)
	text.Draw(image, value, TextFace, op)
	buffer := make([]byte, 4*g.Screen.Width*g.Screen.Height)
	image.ReadPixels(buffer)
	for i := 0; i < len(buffer); i += 4 {
		r := buffer[i]
		green := buffer[i+1]
		b := buffer[i+2]
		a := buffer[i+3]
		g.Screen.Buffer[y+i/4/128][x+i/4%128] = color.RGBA{r, green, b, a}
	}
}