package main

import (
	"math"
)

func (g *Game) wrapText(value string) []string {
	maxChars := int(math.Round(float64(g.Screen.Width)/float64(g.Screen.FontWidth))) - g.Screen.FontWidth*2 - g.Screen.FontWidth/2
	var lines []string
	for len(value) > maxChars {
		lines = append(lines, value[:maxChars])
		value = value[maxChars:]
	}
	lines = append(lines, value)
	return lines
}

func (g *Game) AppendLine(value string, input bool) {
	wrapped := g.wrapText(value)

	g.LinearBuffer = append(g.LinearBuffer, LinearBuffer{
		Content: wrapped,
		IsInput: input,
	})

	g.TruncateLines(g.Screen.Height / g.Screen.FontSize * 2)
}

func (g *Game) ModifyLine(index int, value string) {
	wrapped := g.wrapText(value)
	
	if len(wrapped) > 0 && index < len(g.LinearBuffer) {
		g.LinearBuffer[index].Content = wrapped
	}
	
	g.TruncateLines(g.Screen.Height / g.Screen.FontSize * 2)
}

func (g *Game) TruncateLines(length int) {
	if length < 0 {
		return
	}
	if length >= len(g.LinearBuffer) {
		return
	}
	g.LinearBuffer = g.LinearBuffer[len(g.LinearBuffer)-length:]
}