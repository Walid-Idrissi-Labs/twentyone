package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
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
	styles.EnsureInit()

	lines := []string{
		"",
		"",
		"",
		"",
		"",
		"",
		"              T W E N T Y   O N E",
		"",
		"                Vegas Strip Blackjack",
		"",
		fmt.Sprintf("              Starting balance: $%d", m.game.SessionStart),
		"",
		"              Press any key to start...",
		"",
	}

	width := m.width
	if width < 60 {
		width = 60
	}

	content := make([]string, 0, len(lines))
	for _, line := range lines {
		padding := (width - len(line)) / 2
		if padding < 0 {
			padding = 0
		}
		content = append(content, strings.Repeat(" ", padding)+line)
	}

	joined := strings.Join(content, "\n")
	placed := lipgloss.Place(width, m.height, lipgloss.Center, lipgloss.Center, joined)

	return styles.StyleBackground.Render(placed)
}

func renderTooSmall(m Model) string {
	styles.EnsureInit()
	content := fmt.Sprintf("Terminal too small.\nMinimum: %d x %d\nCurrent: %d x %d\n\nPlease resize your terminal.", styles.MinTermWidth, styles.MinTermHeight, m.width, m.height)
	return styles.StyleBackground.Render(lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, renderModal(content, 50)))
}

func EnsureInit() {
	styles.EnsureInit()
}

func GetStyleBackground() lipgloss.Style {
	return styles.StyleBackground
}

func GetHUDStyle() lipgloss.Style {
	return styles.StyleHUD
}

func GetSuccessStyle() lipgloss.Style {
	return styles.StyleSuccess
}

func GetDangerStyle() lipgloss.Style {
	return styles.StyleDanger
}