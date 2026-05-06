package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/twentyone/twentyone/game"
	"github.com/twentyone/twentyone/styles"
)

func renderCard(c game.Card, highlighted bool, dimmed bool) string {
	var style lipgloss.Style

	if !c.FaceUp {
		style = styles.StyleCard
		if dimmed {
			style = styles.StyleCardDimmed
		}
		return style.Render("┌───┐\n│▓▓▓│\n│▓▓▓│\n│▓▓▓│\n└───┘")
	}

	if c.IsRedSuit() {
		if highlighted {
			style = styles.StyleCardRedHighlighted
		} else {
			style = styles.StyleCardRed
		}
	} else {
		if highlighted {
			style = styles.StyleCardHighlighted
		} else {
			style = styles.StyleCard
		}
	}

	if dimmed {
		style = styles.StyleCardDimmed
	}

	rank := c.RankString()
	suit := c.SuitSymbol()

	return style.Render(fmt.Sprintf("┌───┐\n│ %s │\n│ %s │\n└───┘", rank, suit))
}

func renderButton(label, shortcut string, focused, disabled bool) string {
	if disabled {
		return styles.StyleButtonDisabled.Render(fmt.Sprintf("%s %s", label, shortcut))
	}
	if focused {
		return styles.StyleButtonFocused.Render(fmt.Sprintf("%s %s", label, shortcut))
	}
	return styles.StyleButton.Render(fmt.Sprintf("%s %s", label, shortcut))
}

func renderModal(content string, width int) string {
	lines := strings.Split(content, "\n")
	maxLen := 0
	for _, line := range lines {
		if len(line) > maxLen {
			maxLen = len(line)
		}
	}

	if width > maxLen {
		width = maxLen
	}

	innerWidth := width - 4
	topBorder := "┌" + strings.Repeat("─", innerWidth) + "┐"
	botBorder := "└" + strings.Repeat("─", innerWidth) + "┘"

	var b strings.Builder
	b.WriteString(topBorder + "\n")
	for _, line := range lines {
		padding := innerWidth - len(line)
		if padding < 0 {
			padding = 0
		}
		leftPad := padding / 2
		rightPad := padding - leftPad
		b.WriteString("│")
		b.WriteString(strings.Repeat(" ", leftPad))
		b.WriteString(line)
		b.WriteString(strings.Repeat(" ", rightPad))
		b.WriteString("│\n")
	}
	b.WriteString(botBorder)

	return styles.StyleModal.Render(b.String())
}

func renderBox(content string, width int) string {
	lines := strings.Split(content, "\n")
	maxLen := 0
	for _, line := range lines {
		if len(line) > maxLen {
			maxLen = len(line)
		}
	}

	if width > maxLen {
		width = maxLen
	}

	innerWidth := width - 4
	topBorder := "┌" + strings.Repeat("─", innerWidth) + "┐"
	botBorder := "└" + strings.Repeat("─", innerWidth) + "┘"

	var b strings.Builder
	b.WriteString(topBorder + "\n")
	for _, line := range lines {
		padding := innerWidth - len(line)
		if padding < 0 {
			padding = 0
		}
		leftPad := padding / 2
		rightPad := padding - leftPad
		b.WriteString("│")
		b.WriteString(strings.Repeat(" ", leftPad))
		b.WriteString(line)
		b.WriteString(strings.Repeat(" ", rightPad))
		b.WriteString("│\n")
	}
	b.WriteString(botBorder)

	return b.String()
}