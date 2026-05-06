package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/twentyone/twentyone/styles"
)

func renderBet(m Model) string {
	styles.EnsureInit()

	hud := renderBetHUD(m)

	content := strings.Builder{}
	content.WriteString(hud)
	content.WriteString("\n\n")
	content.WriteString("                    Place Your Bet\n\n")

	minDisplay := m.minBet
	maxDisplay := m.maxBet
	if maxDisplay == 0 {
		maxDisplay = m.game.Balance
	}

	content.WriteString(fmt.Sprintf("                   Min: $%d  Max: $%d\n\n", minDisplay, maxDisplay))

	dealEnabled := m.currentBet >= m.minBet && m.currentBet <= m.game.Balance
	dealBtn := renderButton("Deal", "Enter", false, !dealEnabled)
	quitBtn := renderButton("Quit", "Q", false, false)

	content.WriteString(fmt.Sprintf("                   %s\n\n", dealBtn))
	content.WriteString(fmt.Sprintf("                   %s\n", quitBtn))

	placed := lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, content.String())
	return styles.StyleBackground.Render(placed)
}

func renderBetHUD(m Model) string {
	balanceStr := fmt.Sprintf("Balance: $%d", m.game.Balance)
	roundStr := fmt.Sprintf("Round %d", m.roundCount+1)

	return lipgloss.NewStyle().
		Background(styles.ColorSurface).
		Foreground(styles.ColorText).
		Render(fmt.Sprintf("%s              twentyone              %s", balanceStr, roundStr))
}