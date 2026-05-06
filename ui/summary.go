package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/twentyone/twentyone/styles"
)

func renderSummary(m Model) string {
	styles.EnsureInit()

	startingBalance := m.game.SessionStart
	finalBalance := m.game.Balance
	netProfit := finalBalance - startingBalance
	roundsPlayed := m.roundCount

	var profitStr string
	if netProfit > 0 {
		profitStr = fmt.Sprintf("+$%d ▲", netProfit)
	} else if netProfit < 0 {
		profitStr = fmt.Sprintf("-$%d ▼", -netProfit)
	} else {
		profitStr = "$0 ="
	}

	lines := []string{
		"",
		"",
		"  Session Complete",
		"",
		fmt.Sprintf("  Starting balance:    $%d", startingBalance),
		fmt.Sprintf("  Final balance:       $%d", finalBalance),
		"  ────────────────────────────",
		fmt.Sprintf("  Net profit:          %s", profitStr),
		"",
		fmt.Sprintf("  Rounds played:       %d", roundsPlayed),
		"",
		"  Thanks for playing twentyone!",
		"",
		"  [ Exit ]",
		"",
	}

	joined := strings.Join(lines, "\n")
	boxed := renderModal(joined, 45)
	return styles.StyleBackground.Render(lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, boxed))
}