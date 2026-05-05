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
	balance := m.game.Balance

	content := strings.Builder{}
	content.WriteString(hud)
	content.WriteString("\n\n")
	content.WriteString("                    Place Your Bet\n\n")
	content.WriteString("                   ┌────────────┐\n")
	content.WriteString(fmt.Sprintf("                   │  $%-10s│\n", m.betInput.Value))
	content.WriteString("                   └────────────┘\n\n")

	minDisplay := m.minBet
	maxDisplay := m.maxBet
	if maxDisplay == 0 {
		maxDisplay = balance
	}
	content.WriteString(fmt.Sprintf("          Min: $%d          Max: $%d\n\n", minDisplay, maxDisplay))

	content.WriteString("    [ - $25 ]  [ - $5 ]  [ + $5 ]  [ + $25 ]\n\n")

	dealEnabled := m.currentBet >= m.minBet && m.currentBet <= balance && m.currentBet <= balance
	dealBtn := renderButton("Deal", "↵", false, !dealEnabled)
	content.WriteString(fmt.Sprintf("              %s\n\n", dealBtn))

	quitBtn := renderButton("Quit", "Q", false, false)
	content.WriteString(fmt.Sprintf("              %s\n", quitBtn))

	bg := lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, content.String())
	return styles.GetStyleBackground().Render(bg)
}

func renderBetHUD(m Model) string {
	balanceStr := fmt.Sprintf("Balance: $%d", m.game.Balance)
	roundStr := fmt.Sprintf("Round %d", m.roundCount+1)

	left := balanceStr
	center := "twentyone"
	right := roundStr

	bar := fmt.Sprintf("%s    %s    %s", left, center, right)
	return styles.GetHUDStyle().Render(bar)
}