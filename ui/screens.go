package ui

import (
	"fmt"
	"strings"

	"github.com/twentyone/twentyone/styles"
)

type Screen int

const (
	ScreenWelcome Screen = iota
	ScreenBet
	ScreenTable
	ScreenResult
	ScreenSummary
	ScreenTooSmall
)

func renderWelcome(m Model) string {
	lines := []string{
		"",
		"",
		"",
		"",
		"",
		"           T W E N T Y   O N E                           ",
		"",
		"                Vegas Strip Blackjack                    ",
		"",
		fmt.Sprintf("              Starting balance: $%d                    ", m.game.SessionStart),
		"",
		"              Press any key to start…                    ",
		"",
	}

	width := m.width
	content := ""
	for _, line := range lines {
		padding := (width - len(stripAnsi(line))) / 2
		if padding < 0 {
			padding = 0
		}
		content += strings.Repeat(" ", padding) + line + "\n"
	}

	return styles.StyleBackground.Render(content)
}

func renderTooSmall(m Model) string {
	content := fmt.Sprintf("Terminal too small.\nMinimum: %d × %d\nCurrent: %d × %d\n\nPlease resize your terminal.", styles.MinTermWidth, styles.MinTermHeight, m.width, m.height)
	return renderModal(content, 50)
}