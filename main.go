package main

import (
	"bytes"
	"fmt"
	"image/color"
	_ "image/png"
	"log"
	"os"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"

	lua "github.com/yuin/gopher-lua"
)

type Game struct {
	Screen       ScreenSpecs
	LuaVM        *lua.LState
	Navbar       Navbar
	Input        Input
	LinearBuffer []LinearBuffer
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
	CliBgColor      color.RGBA
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

var (
	TextFaceSource *text.GoTextFaceSource
	TextFace       text.Face
	CodeEditors    = make(map[int]*CodeEditor) // map of file ids to CodeEditor instances
)

func (g *Game) HandleCommand(command string) {
	var err = g.LuaVM.DoString(command)
	if strings.Split(command, " ")[0] == "run" {
		err := g.LuaVM.DoFile(strings.TrimSpace(strings.Split(command, " ")[1]))
		if err != nil {
			log.Println("Error executing file:", err)
			g.AppendLine(err.Error(), false)
		}
		return
	}
	if err != nil {
		log.Println("Error executing command:", err)
		g.AppendLine(err.Error(), false)
	}
	g.AppendLine(g.Input.CurrentInputString, true)
}

func (g *Game) Update() error {
	if updateFn := g.LuaVM.GetGlobal("_update"); updateFn != lua.LNil {
		if err := g.LuaVM.CallByParam(lua.P{
			Fn:      updateFn,
			NRet:    0,
			Protect: true,
		}); err != nil {
			g.AppendLine("Lua error in _update: "+err.Error(), false)
		}
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
			g.HandleCommand(g.Input.CurrentInputString)
			g.Input.CurrentInputString = ""
		default:
			g.Input.CurrentInputString += string(r)
		}
	}

	// on enter, just clear the "buffer" (if that's the word) and append the contents as a new line.
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		g.Input.Keys = []ebiten.Key{}
		g.HandleCommand(g.Input.CurrentInputString)
		g.Input.CurrentInputString = ""
	}

	// if backspace'd, remove one character
	if inpututil.IsKeyJustPressed(ebiten.KeyBackspace) {
		if len(g.Input.CurrentInputString) > 0 {
			g.Input.CurrentInputString = g.Input.CurrentInputString[:len(g.Input.CurrentInputString)-1]
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

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Clear()

	if drawFn := g.LuaVM.GetGlobal("_draw"); drawFn != lua.LNil {
		if err := g.LuaVM.CallByParam(lua.P{
			Fn:      drawFn,
			NRet:    0,
			Protect: true,
		}); err != nil {
			g.AppendLine("Lua error in _draw: "+err.Error(), false)
		}
	}

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
				text.Draw(img, prefix+wrappedLine, TextFace, &text.DrawOptions{})

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
	screen.Fill(g.Screen.BgColor)

	lineHeight := g.Screen.FontSize + 1
	availableHeight := g.Screen.Height - navbarHeight
	maxVisibleLines := availableHeight / lineHeight

	startLine := 0
	if len(g.LinearBuffer) > maxVisibleLines {
		startLine = len(g.LinearBuffer) - maxVisibleLines
	}

	y := navbarHeight
	for i := startLine; i < len(g.LinearBuffer) && i < startLine+maxVisibleLines; i++ {
		line := g.LinearBuffer[i]
		lineNumber := i + 1

		for j, wrappedLine := range line.Content {
			img := ebiten.NewImage(g.Screen.Width, lineHeight)

			prefix := ""
			if j == 0 {
				prefix = fmt.Sprintf("%3d ", lineNumber)
			} else {
				prefix = "    "
			}

			content := prefix + wrappedLine

			if i == len(g.LinearBuffer)-1 && line.IsInput {
				content = prefix + g.Input.CurrentInputString
			}

			text.Draw(img, content, TextFace, &text.DrawOptions{})

			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(0, float64(y))
			screen.DrawImage(img, op)

			y += lineHeight
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

	// load every icon
	for i := range navbar.Tabs {
		iconData, err := os.ReadFile(navbar.Tabs[i].IconPath)
		if err != nil {
			log.Fatal("Error reading icon file:", err)
		}
		navbar.Tabs[i].Icon, _, err = ebitenutil.NewImageFromReader(bytes.NewReader(iconData))
		if err != nil {
			log.Fatal("Error loading icon:", err)
		}
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
	var mouse, _, err = ebitenutil.NewImageFromFile(game.Input.MouseSkin)
	if err != nil {
		log.Fatal("Error loading mouse skin:", err)
	}
	mouseShadow, _, err := ebitenutil.NewImageFromFile(game.Input.MouseShadowPath)
	if err != nil {
		log.Fatal("Error loading mouse shadow:", err)
	}

	game.Input.Mouse = mouse
	game.Input.MouseShadow = mouseShadow
	game.setupLuaAPI()
	game.AppendLine("", true)

	font, err := os.ReadFile(screen.Font)
	if err != nil {
		log.Fatal(err)
	}
	TextFaceSource, err = text.NewGoTextFaceSource(bytes.NewReader(font))
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
