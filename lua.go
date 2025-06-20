package main

import (
	"image/color"
	"log"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	
	lua "github.com/yuin/gopher-lua"
)

func (g *Game) setupLuaAPI() {
	dofiTable := g.LuaVM.NewTable()
	g.LuaVM.SetGlobal("dofi", dofiTable)

	// cls function - clear screen
	g.LuaVM.SetField(dofiTable, "cls", g.LuaVM.NewFunction(func(L *lua.LState) int {
		for x := 0; x < 128; x++ {
			for y := 0; y < 128; y++ {
				g.Screen.Buffer[x][y] = color.RGBA{0, 0, 0, 255}
			}
		}
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

		if x >= 0 && x < 128 && y >= 0 && y < 128 {
			g.Screen.Buffer[x][y] = color.RGBA{r, green, b, 255}
		}
		return 0
	}))

	g.LuaVM.SetGlobal("print", g.LuaVM.NewFunction(func(L *lua.LState) int {
		top := L.GetTop()
		var parts []string
		for i := 1; i <= top; i++ {
			parts = append(parts, L.ToString(i))
		}
		log.Println(strings.Join(parts, " "))
		g.AppendLine(strings.Join(parts, " "), false)
		return 0
	}))
}

func (g *Game) ClearLines() {
	g.LinearBuffer = []LinearBuffer{}
	g.Input.CurrentInputString = ""
	g.Input.Keys = []ebiten.Key{}
}