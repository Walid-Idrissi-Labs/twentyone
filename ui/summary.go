package ui

import (
	"fmt"
	"strings"

	"github.com/twentyone/twentyone/styles"
)

func renderSummary(m Model) string {
	styles.EnsureInit()

	startingBalance := m.game.SessionStart
	finalBalance := m.game.Balance
	netProfit := finalBalance - startingBalance
	roundsPlayed := m.roundCount

	profitStr := formatProfit(netProfit)

	lines := []string{
		"",
		"",
		"               Session Complete                          ",
		"",
		fmt.Sprintf("           Starting balance:    $%d                   ", startingBalance),
		fmt.Sprintf("           Final balance:       $%d                   ", finalBalance),
		"           ─────────────────────────────                 ",
		fmt.Sprintf("           Net profit:          %s", profitStr),
		"",
		fmt.Sprintf("           Rounds played:       %d                       ", roundsPlayed),
		"",
		"           Thanks for playing twentyone!                 ",
		"",
		"                    [ Exit ]                             ",
		"",
	}

	joined := strings.Join(lines, "\n")
	return styles.GetStyleBackground().Render(joined)
}