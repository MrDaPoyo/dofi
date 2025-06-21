package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"image/color"
	_ "image/png"
	"log"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"

	lua "github.com/yuin/gopher-lua"
)

type Game struct {
	Screen        ScreenSpecs
	LuaVM         *lua.LState
	Navbar        Navbar
	Input         Input
	LinearBuffer  []LinearBuffer
	ScriptRunning bool
}

type ScreenSpecs = struct {
	Width           int
	Height          int
	UpscalingFactor int
	Font            string
	FontSize        int
	FontWidth       int
	Buffer          [128][128]color.RGBA
	ImageBuffer     []*ebiten.Image
	BgColor         color.RGBA
	BgTextColor     color.RGBA
	CliBgColor      color.RGBA
	CliColor        color.RGBA
}

type Navbar = struct {
	Tabs         []Tab
	CurrentTab   int
	CliEnabled   bool
	NavbarColor  color.RGBA
	TabColor     color.RGBA
	NavbarHeight int
}

type Tab = struct {
	Name     string
	Enabled  bool
	IconPath string
	Icon     *ebiten.Image
	Function func(*Game) error
}

type LinearBuffer = struct {
	Content []string
	IsInput bool
}

type Input = struct {
	Keys               []ebiten.Key
	CurrentInputString string
	Mouse              *ebiten.Image
	MouseX             int
	MouseY             int
	MouseSkin          string
	MouseShadowPath    string
	MouseShadow        *ebiten.Image
	IsMouseDown        bool
}

type CodeEditor struct {
	Content []string
	Line    int
	Column  int
	ScrollY int
	Saved   bool
}

//go:embed donut.lua
var donutLua []byte

//go:embed resources/cg-pixel-4x5-mono.otf
var fontBytes []byte

//go:embed resources/icons/code.png
var codeIcon []byte

//go:embed resources/icons/brush.png
var brushIcon []byte

//go:embed resources/icons/tile.png
var tileIcon []byte

//go:embed resources/icons/play.png
var playIcon []byte

//go:embed resources/icons/music.png
var musicIcon []byte

//go:embed resources/icons/mouse.png
var mouseIcon []byte

//go:embed resources/icons/mouse_shadow.png
var mouseShadowIcon []byte

var (
	TextFaceSource    *text.GoTextFaceSource
	TextFace          text.Face
	CodeEditors       = make(map[int]*CodeEditor) // map of file ids to CodeEditor instances
	CodeEditorIndex   = 0
	CursorBlinkFrames = 0
	CursorBlinkRate   = 30
	LuaExamples       = map[string]string{
		"donut": string(donutLua),
		"print": `print("hello world")`,
	}
)

func (g *Game) HandleCommand(command string) {
	command = strings.TrimSpace(command)
	if command == "" {
		return
	}

	if command == "help" {
		g.AppendLine("Available commands:", false)
		g.AppendLine("help - Show this help message", false)
		g.AppendLine("cls - Clear the screen", false)
		g.AppendLine("dofi.pset(x,y,r,g,b) - Set pixel at (x, y) to color (r, g, b)", false)
		g.AppendLine(command, true)
		return
	}

	if strings.HasPrefix(command, "example") {
		exampleName := strings.TrimPrefix(command, "example ")
		if exampleLua, exists := LuaExamples[exampleName]; exists {
			if err := g.RunLuaScript(exampleLua); err != nil {
				g.AppendLine("Error running example: "+err.Error(), false)
			} else {
				g.ScriptRunning = true
				g.AppendLine("Running example: "+exampleName, false)
			}
		} else {
			g.AppendLine("Example not found: "+exampleName, false)
		}
		return
	}
}

func (g *Game) Update() (err error) {
	if g.ScriptRunning {
		if updateFn := g.LuaVM.GetGlobal("_update"); updateFn != lua.LNil {
			if err := g.LuaVM.CallByParam(lua.P{
				Fn:      updateFn,
				NRet:    0,
				Protect: true,
			}); err != nil {
				g.AppendLine("Lua error in _update: "+err.Error(), false)
			}
		}
		return nil
	}

	if !g.Navbar.CliEnabled {
		g.Input.MouseX, g.Input.MouseY = ebiten.CursorPosition()
		g.Input.IsMouseDown = ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)
	}

	var inputChars []rune
	inputChars = ebiten.AppendInputChars(inputChars[:0])

	for _, r := range inputChars {
		switch r {
		case '\r', '\n':
			if g.Navbar.CliEnabled {
				g.HandleCommand(g.Input.CurrentInputString)
				g.Input.CurrentInputString = ""
			}
		default:
			if g.Navbar.CliEnabled {
				g.Input.CurrentInputString += string(r)
			} else {
				if editor, exists := CodeEditors[CodeEditorIndex]; exists {
					if editor.Line < len(editor.Content) {
						line := editor.Content[editor.Line]
						if editor.Column <= len(line) {
							editor.Content[editor.Line] = line[:editor.Column] + string(r) + line[editor.Column:]
							editor.Column++
						}
					}
				}
			}
		}
	}

	// on enter, just clear the "buffer" (if that's the word) and append the contents as a new line.
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		if g.Navbar.CliEnabled {
			g.Input.Keys = []ebiten.Key{}
			g.HandleCommand(g.Input.CurrentInputString)
			g.Input.CurrentInputString = ""
		} else {
			if editor, exists := CodeEditors[CodeEditorIndex]; exists {
				editor.Content = append(editor.Content, "")
				editor.Line++
				editor.Column = 0
			}
		}
	}

	// if backspace'd, remove one character
	if inpututil.IsKeyJustPressed(ebiten.KeyBackspace) {
		if g.Navbar.CliEnabled {
			if len(g.Input.CurrentInputString) > 0 {
				g.Input.CurrentInputString = g.Input.CurrentInputString[:len(g.Input.CurrentInputString)-1]
			}
		} else {
			if editor, exists := CodeEditors[CodeEditorIndex]; exists {
				if editor.Line < len(editor.Content) && editor.Column > 0 {
					line := editor.Content[editor.Line]
					editor.Content[editor.Line] = line[:editor.Column-1] + line[editor.Column:]
					editor.Column--
				} else if editor.Line > 0 {
					prevLine := editor.Content[editor.Line-1]
					editor.Column = len(prevLine)
					if editor.Line < len(editor.Content)-1 {
						editor.Content = append(editor.Content[:editor.Line], editor.Content[editor.Line+1:]...)
					} else {
						editor.Content = editor.Content[:editor.Line]
					}
					editor.Line--
				}
			}
		}
	}

	// if escaped, switch tabs (POYO (yes me) THIS IS A TODO) (switching tabs here)
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		g.Navbar.CliEnabled = !g.Navbar.CliEnabled
	}

	// if input and CliEnabled, change the contents
	if len(g.LinearBuffer) > 0 && g.Navbar.CliEnabled {
		g.ModifyLine(len(g.LinearBuffer)-1, g.Input.CurrentInputString)
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyLeft) {
		if !g.Navbar.CliEnabled {
			if editor, exists := CodeEditors[CodeEditorIndex]; exists {
				if editor.Column > 0 {
					editor.Column--
				} else if editor.Line > 0 {
					editor.Line--
					editor.Column = len(editor.Content[editor.Line])
				}
			}
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyRight) {
		if !g.Navbar.CliEnabled {
			if editor, exists := CodeEditors[CodeEditorIndex]; exists {
				if editor.Column < len(editor.Content[editor.Line]) {
					editor.Column++
				} else if editor.Line < len(editor.Content)-1 {
					editor.Line++
					editor.Column = 0
				}
			}
		}
	}

	CursorBlinkFrames++
	if CursorBlinkFrames > 60 {
		CursorBlinkFrames = 0
	}

	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("recovered: %v", r)
			log.Println(err)
		}
	}()

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Clear()

	bufferImg := ebiten.NewImage(len(g.Screen.Buffer[0]), len(g.Screen.Buffer))

	pixels := make([]byte, len(g.Screen.Buffer)*len(g.Screen.Buffer[0])*4)
	for y := 0; y < len(g.Screen.Buffer); y++ {
		for x := 0; x < len(g.Screen.Buffer[y]); x++ {
			pixel := g.Screen.Buffer[y][x]
			idx := (y*len(g.Screen.Buffer[0]) + x) * 4
			pixels[idx] = pixel.R
			pixels[idx+1] = pixel.G
			pixels[idx+2] = pixel.B
			pixels[idx+3] = pixel.A
		}
	}

	bufferImg.WritePixels(pixels)
	screen.DrawImage(bufferImg, &ebiten.DrawImageOptions{})

	if !g.ScriptRunning {
		if g.Navbar.CliEnabled {
			screen.Fill(g.Screen.CliBgColor)
			lineHeight := g.Screen.FontSize + 1
			var totalLines int
			for _, line := range g.LinearBuffer {
				totalLines += len(line.Content)
			}

			y := g.Screen.Height - totalLines*lineHeight
			for _, line := range g.LinearBuffer {
				prefix := "- "
				if line.IsInput {
					prefix = "> "
				}
				for _, wrappedLine := range line.Content {
					img := ebiten.NewImage(g.Screen.Width, lineHeight)
					textOP := &text.DrawOptions{}
					textOP.ColorScale.ScaleWithColor(g.Screen.CliColor)
					text.Draw(img, prefix+wrappedLine, TextFace, textOP)

					op := &ebiten.DrawImageOptions{}
					op.GeoM.Translate(0, float64(y))
					screen.DrawImage(img, op)

					y += lineHeight
					prefix = "  "
				}
			}
		} else {
			screen.Fill(g.Screen.BgColor)
			navbarHeight := g.Navbar.NavbarHeight
			navbarImg := ebiten.NewImage(g.Screen.Width, navbarHeight)
			navbarImg.Fill(g.Navbar.NavbarColor)
			screen.DrawImage(navbarImg, nil)

			totalTabWidth := 0
			enabledTabs := 0
			for _, tab := range g.Navbar.Tabs {
				if tab.Enabled {
					enabledTabs++
					totalTabWidth += tab.Icon.Bounds().Dx() + 2
				}
			}
			if enabledTabs > 0 {
				totalTabWidth += (enabledTabs - 1) * 2 // spacing between tabs
			}

			xPosition := g.Screen.Width - totalTabWidth - 1
			for _, tab := range g.Navbar.Tabs {
				if !tab.Enabled {
					continue
				}

				iconWidth := tab.Icon.Bounds().Dx() + 2
				iconHeight := tab.Icon.Bounds().Dy() + 2

				tabImg := ebiten.NewImage(iconWidth, iconHeight)
				tabImg.Fill(g.Navbar.TabColor)

				iconOP := &ebiten.DrawImageOptions{}
				iconOP.GeoM.Translate(float64((iconWidth-tab.Icon.Bounds().Dx())/2), float64((iconHeight-tab.Icon.Bounds().Dy())/2))
				tabImg.DrawImage(tab.Icon, iconOP)

				tabOp := &ebiten.DrawImageOptions{}
				tabOp.GeoM.Translate(float64(xPosition), float64((navbarHeight-iconHeight)/2))
				screen.DrawImage(tabImg, tabOp)

				xPosition += iconWidth + 2
			}

			var contentImage = ebiten.NewImage(g.Screen.Width, g.Screen.Height-navbarHeight)
			contentImage.Fill(g.Screen.BgColor)

			if g.Navbar.Tabs[g.Navbar.CurrentTab].Name == "code" {
				if editor, exists := CodeEditors[g.Navbar.CurrentTab]; exists {
					g.CodeEditor(contentImage, editor, navbarHeight)
				} else {
					CodeEditors[g.Navbar.CurrentTab] = &CodeEditor{
						Content: []string{""},
						Line:    0,
						Column:  0,
						ScrollY: 0,
						Saved:   false,
					}
					g.CodeEditor(screen, CodeEditors[g.Navbar.CurrentTab], navbarHeight)
				}
			}
			var contentImageOp = &ebiten.DrawImageOptions{}
			contentImageOp.GeoM.Translate(0, float64(navbarHeight))
			screen.DrawImage(contentImage, contentImageOp)

			g.DrawMouse(screen)
		}
	} else {
		if drawFn := g.LuaVM.GetGlobal("_draw"); drawFn != lua.LNil {
			if err := g.LuaVM.CallByParam(lua.P{
				Fn:      drawFn,
				NRet:    0,
				Protect: true,
			}); err != nil {
				log.Println("Lua error in _draw:", err)
				g.AppendLine("Lua error in _draw: "+err.Error(), false)
			}
			screen.DrawImage(bufferImg, nil)
		}
	}
}

func (g *Game) DrawMouse(screen *ebiten.Image) {
	mouseOp := &ebiten.DrawImageOptions{}
	mouseOp.GeoM.Translate(float64(g.Input.MouseX)-1, float64(g.Input.MouseY)-1)
	mouseShadowOp := &ebiten.DrawImageOptions{}
	mouseShadowOp.GeoM.Translate(float64(g.Input.MouseX)-1, float64(g.Input.MouseY)-1)
	screen.DrawImage(g.Input.MouseShadow, mouseShadowOp)
	screen.DrawImage(g.Input.Mouse, mouseOp)
}

func (g *Game) CodeEditor(screen *ebiten.Image, editor *CodeEditor, navbarHeight int) *Game {
	screen.Clear()
	screen.Fill(g.Screen.CliBgColor)

	lineHeight := g.Screen.FontSize + 2
	availableHeight := g.Screen.Height - navbarHeight
	maxVisibleLines := availableHeight / lineHeight

	startLine := 0
	if len(g.LinearBuffer) > maxVisibleLines {
		startLine = len(g.LinearBuffer) - maxVisibleLines
	}

	var y = 0

	for i := startLine; i < len(CodeEditors[CodeEditorIndex].Content) && i < startLine+maxVisibleLines; i++ {
		line := CodeEditors[CodeEditorIndex].Content[i]
		wrappedLines := g.wrapText(line, g.Screen.Width)
		for _, wrappedLine := range wrappedLines {
			img := ebiten.NewImage(g.Screen.Width, lineHeight)

			if i%2 == 0 {
				img.Fill(color.RGBA{g.Screen.CliBgColor.R - 10, g.Screen.CliBgColor.G - 10, g.Screen.CliBgColor.B - 10, g.Screen.CliBgColor.A})
			} else {
				img.Fill(color.RGBA{g.Screen.CliBgColor.R + 20, g.Screen.CliBgColor.G + 20, g.Screen.CliBgColor.B + 20, g.Screen.CliBgColor.A})
			}
			textOP := &text.DrawOptions{}
			textOP.GeoM.Translate(0, 1)
			textOP.ColorScale.ScaleWithColor(g.Screen.CliColor)
			text.Draw(img, wrappedLine, TextFace, textOP)

			screenOP := &ebiten.DrawImageOptions{}
			screenOP.GeoM.Translate(0, float64(y*lineHeight))
			screen.DrawImage(img, screenOP)
			y++
		}
	}

	cursorVisualLine := 0
	cursorVisualColumn := editor.Column

	for i := startLine; i < editor.Line && i < len(editor.Content); i++ {
		wrappedLines := g.wrapText(editor.Content[i], g.Screen.Width)
		cursorVisualLine += len(wrappedLines)
	}

	if editor.Line < len(editor.Content) {
		currentLineWrapped := g.wrapText(editor.Content[editor.Line], g.Screen.Width)
		charsProcessed := 0
		for wrapIndex, wrappedSegment := range currentLineWrapped {
			if editor.Column <= charsProcessed+len(wrappedSegment) {
				cursorVisualLine += wrapIndex
				cursorVisualColumn = editor.Column - charsProcessed
				break
			}
			charsProcessed += len(wrappedSegment)
		}
	}

	// only draw the cursor if it's on the screen
	if cursorVisualLine < maxVisibleLines && cursorVisualLine >= 0 {
		cursorX := cursorVisualColumn * (g.Screen.FontWidth + 1)
		if cursorX <= g.Screen.Width {
			cursorImg := ebiten.NewImage(1, lineHeight-2)
			cursorImg.Fill(color.White)
			cursorOP := &ebiten.DrawImageOptions{}
			if CursorBlinkFrames < CursorBlinkRate {
				cursorOP.ColorScale.ScaleWithColor(color.White)
			} else {
				cursorOP.ColorScale.ScaleWithColor(color.RGBA{180, 180, 180, 255})
			}
			cursorOP.GeoM.Translate(float64(cursorX), float64(cursorVisualLine*lineHeight+1))
			screen.DrawImage(cursorImg, cursorOP)
		}
	}

	return g
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return g.Screen.Width, g.Screen.Height
}

func MakeGame() *Game {
	var screen = ScreenSpecs{
		Width:           128,
		Height:          128,
		UpscalingFactor: 4,
		Font:            "resources/cg-pixel-4x5-mono.otf",
		FontSize:        5,
		FontWidth:       4,
		BgColor:         color.RGBA{255, 169, 133, 255},
		CliBgColor:      color.RGBA{70, 82, 113, 255},
		CliColor:        color.RGBA{255, 255, 255, 255},
		BgTextColor:     color.RGBA{0, 86, 122, 255},
	}

	var navbar = Navbar{
		Tabs: []Tab{
			{Name: "code", Enabled: true, IconPath: "resources/icons/code.png"},
			{Name: "draw", Enabled: true, IconPath: "resources/icons/brush.png"},
			{Name: "tile", Enabled: true, IconPath: "resources/icons/tile.png"},
			{Name: "play", Enabled: true, IconPath: "resources/icons/play.png"},
			{Name: "music", Enabled: true, IconPath: "resources/icons/music.png"},
		},
		CurrentTab:   0,
		NavbarColor:  color.RGBA{204, 116, 83, 255},
		TabColor:     color.RGBA{154, 56, 63, 255},
		NavbarHeight: 10,
		CliEnabled:   true,
	}

	var lineBuffer = make([]*ebiten.Image, 1)
	lineBuffer[0] = ebiten.NewImage(screen.Width, 128)

	iconMap := map[string][]byte{
		"resources/icons/code.png":  codeIcon,
		"resources/icons/brush.png": brushIcon,
		"resources/icons/tile.png":  tileIcon,
		"resources/icons/play.png":  playIcon,
		"resources/icons/music.png": musicIcon,
	}

	for i := range navbar.Tabs {
		iconData, _, err := ebitenutil.NewImageFromReader(bytes.NewReader(iconMap[navbar.Tabs[i].IconPath]))
		if err != nil {
			log.Fatal("Error loading icon:", err)
		}
		navbar.Tabs[i].Icon = iconData
	}

	var game = Game{
		Navbar: navbar,
		Screen: screen,
		LuaVM:  lua.NewState(),
		Input: Input{
			CurrentInputString: "",
			MouseX:             0,
			MouseY:             0,
			MouseSkin:          "resources/icons/mouse.png",
			MouseShadowPath:    "resources/icons/mouse_shadow.png",
			IsMouseDown:        false,
		},
	}

	ebiten.SetCursorMode(ebiten.CursorModeHidden)

	mouse, _, err := ebitenutil.NewImageFromReader(bytes.NewReader(mouseIcon))
	if err != nil {
		log.Fatal("Error loading mouse skin:", err)
	}
	mouseShadow, _, err := ebitenutil.NewImageFromReader(bytes.NewReader(mouseShadowIcon))
	if err != nil {
		log.Fatal("Error loading mouse shadow:", err)
	}

	game.Input.Mouse = mouse
	game.Input.MouseShadow = mouseShadow
	game.setupLuaAPI()
	game.AppendLine("", true)

	TextFaceSource, err = text.NewGoTextFaceSource(bytes.NewReader(fontBytes))
	if err != nil {
		log.Fatal(err)
	}
	TextFace = &text.GoTextFace{
		Source: TextFaceSource,
		Size:   float64(game.Screen.FontSize),
	}

	return &game
}

func main() {
	game := MakeGame()
	defer game.LuaVM.Close()

	ebiten.SetWindowSize(game.Screen.Width*game.Screen.UpscalingFactor, game.Screen.Height*game.Screen.UpscalingFactor)
	ebiten.SetWindowTitle("Dofi! :3")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
