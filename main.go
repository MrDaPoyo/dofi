package main

import (
	"flag"
	"fmt"
	"image/color"
	"log"
	"os"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/yuin/gopher-lua"
)

const version = "1.0.0"

type Game struct {
	inputBuf         string
	outputs          []string
	previousCommands []string
	historyIndex     int // -1 means not browsing history
	luaState         *lua.LState
	gamePixels       [128][128]color.RGBA // 128x128 pixel buffer
	gameLoaded       bool
}

func NewGame() *Game {
	g := &Game{
		outputs: []string{
			"Dofi Fantasy Console v" + version,
			"Type 'help' for commands",
		},
		historyIndex: -1,
		luaState:     lua.NewState(),
	}
	g.setupLuaAPI()
	return g
}

func (g *Game) setupLuaAPI() {
	// Create dofi table
	dofiTable := g.luaState.NewTable()
	g.luaState.SetGlobal("dofi", dofiTable)

	// cls function - clear screen
	g.luaState.SetField(dofiTable, "cls", g.luaState.NewFunction(func(L *lua.LState) int {
		for x := 0; x < 128; x++ {
			for y := 0; y < 128; y++ {
				g.gamePixels[x][y] = color.RGBA{0, 0, 0, 255}
			}
		}
		return 0
	}))

	// pset function - set pixel
	g.luaState.SetField(dofiTable, "pset", g.luaState.NewFunction(func(L *lua.LState) int {
		x := int(L.CheckNumber(1))
		y := int(L.CheckNumber(2))
		r := uint8(L.CheckNumber(3))
		green := uint8(L.CheckNumber(4))
		b := uint8(L.CheckNumber(5))

		if x >= 0 && x < 128 && y >= 0 && y < 128 {
			g.gamePixels[x][y] = color.RGBA{r, green, b, 255}
		}
		return 0
	}))

	// Override print to capture output
	g.luaState.SetGlobal("print", g.luaState.NewFunction(func(L *lua.LState) int {
		top := L.GetTop()
		var parts []string
		for i := 1; i <= top; i++ {
			parts = append(parts, L.ToString(i))
		}
		g.outputs = append(g.outputs, strings.Join(parts, " "))
		return 0
	}))
}

// Update handles key presses and text input
func (g *Game) Update() error {
	// If a game is loaded, call its _update function
	if g.gameLoaded {
		if updateFn := g.luaState.GetGlobal("_update"); updateFn != lua.LNil {
			if err := g.luaState.CallByParam(lua.P{
				Fn:      updateFn,
				NRet:    0,
				Protect: true,
			}); err != nil {
				g.outputs = append(g.outputs, "Lua error in _update: "+err.Error())
			}
		}
	}

	// Arrow key navigation for command history
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) {
		if len(g.previousCommands) > 0 {
			if g.historyIndex == -1 {
				g.historyIndex = len(g.previousCommands) - 1
			} else if g.historyIndex > 0 {
				g.historyIndex--
			}
			g.inputBuf = g.previousCommands[g.historyIndex]
		}
		return nil
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) {
		if g.historyIndex != -1 {
			if g.historyIndex < len(g.previousCommands)-1 {
				g.historyIndex++
				g.inputBuf = g.previousCommands[g.historyIndex]
			} else {
				g.historyIndex = -1
				g.inputBuf = ""
			}
		}
		return nil
	}

	// Append any typed runes
	for _, r := range ebiten.AppendInputChars(nil) {
		g.inputBuf += string(r)
		g.historyIndex = -1 // Reset history browsing when typing
	}

	// Backspace
	if inpututil.IsKeyJustPressed(ebiten.KeyBackspace) && len(g.inputBuf) > 0 {
		g.inputBuf = g.inputBuf[:len(g.inputBuf)-1]
		g.historyIndex = -1 // Reset history browsing when editing
	}

	// Enter = execute command
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		cmd := strings.TrimSpace(g.inputBuf)
		if cmd != "" {
			// Add to command history
			g.previousCommands = append(g.previousCommands, cmd)
			// Keep only last 50 commands
			if len(g.previousCommands) > 50 {
				g.previousCommands = g.previousCommands[1:]
			}
		}
		g.execute(cmd)
		g.inputBuf = ""
		g.historyIndex = -1
	}

	// ESC = quit
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		return fmt.Errorf("user quit")
	}

	return nil
}

func (g *Game) runGameFile(filename string) error {
	// Reset game state
	g.gameLoaded = false
	for x := 0; x < 128; x++ {
		for y := 0; y < 128; y++ {
			g.gamePixels[x][y] = color.RGBA{0, 0, 0, 255}
		}
	}

	// Execute the Lua file
	if err := g.luaState.DoFile(filename); err != nil {
		return err
	}

	g.gameLoaded = true
	g.outputs = append(g.outputs, "Game loaded successfully: "+filename)
	return nil
}

func (g *Game) execute(cmdLine string) {
	if cmdLine == "" {
		return
	}
	g.outputs = append(g.outputs, "> "+cmdLine)

	parts := strings.Fields(cmdLine)
	switch parts[0] {
	case "help":
		g.outputs = append(g.outputs,
			"help            Show this help",
			"version         Show version",
			"list            List available games",
			"run <file>      Run a game file",
			"stop            Stop current game",
			"exit            Quit the console",
			"",
			"Use Up/Down arrow keys to browse command history",
		)
	case "version":
		g.outputs = append(g.outputs, "Version: "+version)
	case "list":
		g.outputs = append(g.outputs, "Available games:")
		files, err := os.ReadDir(".")
		if err != nil {
			g.outputs = append(g.outputs, "Error reading directory:", err.Error())
			return
		}
		for _, file := range files {
			if !file.IsDir() && strings.HasSuffix(file.Name(), ".lua") {
				g.outputs = append(g.outputs, " - "+file.Name())
			}
		}
	case "run":
		if len(parts) < 2 {
			g.outputs = append(g.outputs, "Usage: run <game.lua>")
			return
		}
		fn := parts[1]
		if _, err := os.Stat(fn); os.IsNotExist(err) {
			g.outputs = append(g.outputs, "Error: file not found:", fn)
		} else {
			if err := g.runGameFile(fn); err != nil {
				g.outputs = append(g.outputs, "Error loading game:", err.Error())
			}
		}
	case "stop":
		g.gameLoaded = false
		g.outputs = append(g.outputs, "Game stopped")
	case "exit":
		g.luaState.Close()
		os.Exit(0)
	default:
		g.outputs = append(g.outputs, "Unknown command:", cmdLine)
	}

	// keep last 100 lines
	if len(g.outputs) > 100 {
		g.outputs = g.outputs[len(g.outputs)-100:]
	}
}

// Draw renders the console text and game pixels
func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.Black)

	// If a game is loaded, draw game pixels first
	if g.gameLoaded {
		// Call _draw function if it exists
		if drawFn := g.luaState.GetGlobal("_draw"); drawFn != lua.LNil {
			if err := g.luaState.CallByParam(lua.P{
				Fn:      drawFn,
				NRet:    0,
				Protect: true,
			}); err != nil {
				g.outputs = append(g.outputs, "Lua error in _draw: "+err.Error())
			}
		}

		// Render the pixel buffer to screen (scaled up)
		scale := 4 // Each game pixel is 4x4 screen pixels
		for x := 0; x < 128; x++ {
			for y := 0; y < 128; y++ {
				pixelColor := g.gamePixels[x][y]
				if pixelColor.R != 0 || pixelColor.G != 0 || pixelColor.B != 0 {
					for dx := 0; dx < scale; dx++ {
						for dy := 0; dy < scale; dy++ {
							screen.Set(x*scale+dx, y*scale+dy, pixelColor)
						}
					}
				}
			}
		}
	} else {
		// Show console text when no game is running - create a virtual 128x128 screen
		virtualScreen := ebiten.NewImage(128, 128)
		virtualScreen.Fill(color.Black)
		
		all := strings.Join(g.outputs, "\n")
		all += "\n> " + g.inputBuf
		ebitenutil.DebugPrint(virtualScreen, all)
		
		// Scale up the virtual screen 4x to fill the 512x512 window
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(4, 4)
		screen.DrawImage(virtualScreen, op)
	}
}

// Layout just uses window size as logical size.
func (g *Game) Layout(w, h int) (int, int) {
	return w, h
}

func main() {
	flag.Parse()
	ebiten.SetWindowSize(480, 480)
	ebiten.SetWindowTitle("Dofi Fantasy Console")
	game := NewGame()
	defer game.luaState.Close()
	
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}