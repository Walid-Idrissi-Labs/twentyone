package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/twentyone/twentyone/game"
	"github.com/twentyone/twentyone/styles"
)

func renderResult(m Model) string {
	styles.EnsureInit()

	tableView := renderTable(Model{
		width:      m.width,
		height:     m.height,
		screen:     ScreenTable,
		game:       m.game,
		anim:       m.anim,
		roundCount: m.roundCount,
	})

	overlay := renderResultOverlay(m)
	content := lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, overlay)

	return tableView + "\n" + content
}

func renderResultOverlay(m Model) string {
	if m.game.Balance == 0 {
		return renderBrokeOverlay()
	}

	header := getResultHeader(m)
	lines := []string{
		"",
		"",
		"  " + header,
		"",
	}

	for i, res := range m.game.Results {
		handNum := i + 1
		resultStr := getResultString(res.Result)
		profitStr := formatProfitWithSign(res.Profit)
		lines = append(lines, fmt.Sprintf("  Hand %d: %s   %s", handNum, profitStr, resultStr))
	}

	lines = append(lines, "")

	totalProfit := 0
	for _, res := range m.game.Results {
		totalProfit += res.Profit
	}
	lines = append(lines, fmt.Sprintf("  Round profit: %s", formatProfitWithSign(totalProfit)))
	lines = append(lines, fmt.Sprintf("  Balance: $%d", m.game.Balance))
	lines = append(lines, "")

	if m.game.Balance > 0 {
		lines = append(lines, "  [ Play Again  Enter ]  [ Quit  Q ]")
	} else {
		lines = append(lines, "  [ Quit  Q ]")
	}
	lines = append(lines, "")

	joined := strings.Join(lines, "\n")
	return renderModal(joined, 55)
}

func renderBrokeOverlay() string {
	lines := []string{
		"",
		"",
		"  Game Over - You're broke!",
		"",
		"  [ Quit  Q ]",
		"",
	}
	joined := strings.Join(lines, "\n")
	return renderModal(joined, 35)
}

func getResultHeader(m Model) string {
	if len(m.game.Results) == 0 {
		return "PUSH"
	}

	allLose := true
	allPush := true
	anyWin := false

	for _, res := range m.game.Results {
		if res.Result != game.ResultLose {
			allLose = false
		}
		if res.Result != game.ResultPush {
			allPush = false
		}
		if res.Result == game.ResultWin || res.Result == game.ResultBlackjack {
			anyWin = true
		}
	}

	if allLose {
		return lipgloss.NewStyle().Foreground(styles.ColorDanger).Bold(true).Render("YOU LOSE")
	}
	if allPush {
		return lipgloss.NewStyle().Foreground(styles.ColorNeutral).Render("PUSH")
	}
	if anyWin {
		return lipgloss.NewStyle().Foreground(styles.ColorSuccess).Bold(true).Render("YOU WIN!")
	}
	return "ROUND OVER"
}

func getResultString(r game.RoundResult) string {
	switch r {
	case game.ResultWin:
		return "Win"
	case game.ResultLose:
		return "Lose"
	case game.ResultPush:
		return "Push"
	case game.ResultBlackjack:
		return "Blackjack"
	default:
		return ""
	}
}

func formatProfitWithSign(p int) string {
	if p >= 0 {
		return fmt.Sprintf("+$%d", p)
	}
	return fmt.Sprintf("-$%d", -p)
}