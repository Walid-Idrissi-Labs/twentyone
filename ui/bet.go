package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/twentyone/twentyone/styles"
)

func renderBet(m Model) string {
	hud := renderBetHUD(m)
	balance := m.balance

	content := strings.Builder{}
	content.WriteString(hud)
	content.WriteString("\n")

	content.WriteString(fmt.Sprintf("                    Place Your Bet\n"))
	content.WriteString(fmt.Sprintf("\n"))
	content.WriteString(fmt.Sprintf("                   ┌────────────┐\n"))
	content.WriteString(fmt.Sprintf("                   │  $%-10s│\n", m.betInput.Value))
	content.WriteString(fmt.Sprintf("                   └────────────┘\n"))
	content.WriteString(fmt.Sprintf("\n"))

	minDisplay := m.minBet
	maxDisplay := m.maxBet
	if maxDisplay == 0 {
		maxDisplay = balance
	}
	content.WriteString(fmt.Sprintf("          Min: $%d          Max: $%d\n", minDisplay, maxDisplay))
	content.WriteString(fmt.Sprintf("\n"))

	content.WriteString(fmt.Sprintf("    [ - $25 ]  [ - $5 ]  [ + $5 ]  [ + $25 ]\n"))
	content.WriteString(fmt.Sprintf("\n"))

	dealEnabled := m.currentBet >= m.minBet && m.currentBet <= balance && m.currentBet <= m.balance
	dealBtn := renderButton("Deal", "↵", false, !dealEnabled)
	content.WriteString(fmt.Sprintf("              %s\n", dealBtn))
	content.WriteString(fmt.Sprintf("\n"))

	quitBtn := renderButton("Quit", "Q", false, false)
	content.WriteString(fmt.Sprintf("              %s\n", quitBtn))

	bg := lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, content.String())
	return styles.StyleBackground.Render(bg)
}

func renderBetHUD(m Model) string {
	balanceStr := fmt.Sprintf("Balance: $%d", m.balance)
	roundStr := fmt.Sprintf("Round %d", m.roundCount+1)

	left := balanceStr
	center := "twentyone"
	right := roundStr

	bar := fmt.Sprintf("%s    %s    %s", left, center, right)
	return styles.HUDStyle.Render(bar)
}