package main

import (
	"math"
	"strings"
)

func (g *Game) wrapText(value string, width int) []string {
	maxChars := int(math.Round(float64(width)/float64(g.Screen.FontWidth))) - g.Screen.FontWidth*2 - g.Screen.FontWidth/2
	var lines []string
	for len(value) > 0 {
		if newlineIndex := strings.Index(value, "\n"); newlineIndex != -1 && newlineIndex < maxChars {
			lines = append(lines, value[:newlineIndex])
			value = value[newlineIndex+1:]
			continue
		}
		
		if len(value) <= maxChars {
			break
		}
		
		lines = append(lines, value[:maxChars])
		value = value[maxChars:]
	}
	lines = append(lines, value)
	return lines
}

func (g *Game) AppendLine(value string, input bool) {
	wrapped := g.wrapText(strings.Replace(value, "\t", "", -1), g.Screen.Width)

	g.LinearBuffer = append(g.LinearBuffer, LinearBuffer{
		Content: wrapped,
		IsInput: input,
	})

	g.TruncateLines(g.Screen.Height / g.Screen.FontSize * 2)
}

func (g *Game) ModifyLine(index int, value string) {
	wrapped := g.wrapText(value, g.Screen.Width)
	
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