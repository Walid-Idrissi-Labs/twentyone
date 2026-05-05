package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/twentyone/twentyone/game"
	"github.com/twentyone/twentyone/styles"
)

func renderCard(c game.Card, highlighted bool, dimmed bool) string {
	var cardStyle lipgloss.Style
	if c.IsRedSuit() {
		if highlighted {
			cardStyle = styles.RedCardHighlightedStyle
		} else {
			cardStyle = styles.RedCardStyle
		}
	} else {
		if highlighted {
			cardStyle = styles.BlackCardHighlightedStyle
		} else {
			cardStyle = styles.BlackCardStyle
		}
	}

	if dimmed {
		cardStyle = styles.CardDimmedStyle
	}

	if !c.FaceUp {
		return cardStyle.Render(
			fmt.Sprintf("┌───┐\n│▓▓▓│\n│▓▓▓│\n│▓▓▓│\n└───┘"),
		)
	}

	rank := c.RankString()
	suitSymbol := c.SuitSymbol()

	lines := []string{
		"┌───┐",
		fmt.Sprintf("│ %s │", rank),
		fmt.Sprintf("│ %s │", suitSymbol),
		"└───┘",
	}

	joined := strings.Join(lines, "\n")
	return cardStyle.Render(joined)
}

func renderButton(label string, shortcut string, focused bool, disabled bool) string {
	if disabled {
		return styles.ButtonDisabledStyle.Render(fmt.Sprintf("%s %s", label, shortcut))
	}
	if focused {
		return styles.ButtonFocusedStyle.Render(fmt.Sprintf("%s %s", label, shortcut))
	}
	return styles.ButtonStyle.Render(fmt.Sprintf("%s %s", label, shortcut))
}

func renderModal(content string, width int) string {
	lines := strings.Split(content, "\n")
	maxLen := 0
	for _, line := range lines {
		l := stripAnsi(line)
		if len(l) > maxLen {
			maxLen = len(l)
		}
	}

	borderWidth := maxLen + 4
	topBorder := "┌" + strings.Repeat("─", borderWidth-2) + "┐"
	bottomBorder := "└" + strings.Repeat("─", borderWidth-2) + "┘"

	result := topBorder + "\n"
	for _, line := range lines {
		l := stripAnsi(line)
		padding := borderWidth - 2 - len(l)
		leftPad := padding / 2
		rightPad := padding - leftPad
		result += "│" + strings.Repeat(" ", leftPad) + line + strings.Repeat(" ", rightPad) + "│\n"
	}
	result += bottomBorder

	return result
}

func stripAnsi(s string) string {
	var result strings.Builder
	inEscape := false
	for _, r := range s {
		if r == '\033' {
			inEscape = true
			continue
		}
		if inEscape && r == 'm' {
			inEscape = false
			continue
		}
		if !inEscape {
			result.WriteRune(r)
		}
	}
	return result.String()
}

func formatMoney(amount int, showSign bool) string {
	sign := ""
	if showSign && amount > 0 {
		sign = "+"
	}
	return fmt.Sprintf("%s$%d", sign, amount)
}

func formatBalance(amount int) string {
	return fmt.Sprintf("$%d", amount)
}

func formatProfit(amount int) string {
	if amount > 0 {
		return fmt.Sprintf("+$%d ▲", amount)
	} else if amount < 0 {
		return fmt.Sprintf("-$%d ▼", -amount)
	}
	return "$0  ="
}