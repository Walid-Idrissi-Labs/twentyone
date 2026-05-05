package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/twentyone/twentyone/game"
)

func renderResult(m Model) string {
	tableView := renderTable(Model{
		width:        m.width,
		height:       m.height,
		screen:       ScreenTable,
		game:         m.game,
		anim:         m.anim,
		roundCount:   m.roundCount,
		buttonAreas:  m.buttonAreas,
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
		"       " + header,
		"",
	}

	for i, res := range m.game.Results {
		handNum := i + 1
		resultStr := getResultString(res.Result)
		profitStr := formatProfitWithSign(res.Profit)
		lines = append(lines, fmt.Sprintf("       Hand %d: %s   %s", handNum, profitStr, resultStr))
	}

	lines = append(lines, "")

	totalProfit := 0
	for _, res := range m.game.Results {
		totalProfit += res.Profit
	}
	lines = append(lines, fmt.Sprintf("       Round profit: %s", formatProfitWithSign(totalProfit)))
	lines = append(lines, fmt.Sprintf("       Balance: $%d", m.game.Balance))
	lines = append(lines, "")

	if m.game.Balance > 0 {
		lines = append(lines, "   [ Play Again  ↵ ]  [ Quit  Q ]")
	} else {
		lines = append(lines, "                    [ Quit  Q ]")
	}
	lines = append(lines, "")

	joined := strings.Join(lines, "\n")
	return renderModal(joined, 60)
}

func renderBrokeOverlay() string {
	lines := []string{
		"",
		"",
		"       Game Over — You're broke!",
		"",
		"                    [ Quit  Q ]",
		"",
	}
	joined := strings.Join(lines, "\n")
	return renderModal(joined, 50)
}

func getResultHeader(m Model) string {
	if len(m.game.Results) == 0 {
		return "PUSH"
	}

	allLose := true
	allPush := true
	anyWin := false
	anyBlackjack := false

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
		if res.Result == game.ResultBlackjack {
			anyBlackjack = true
		}
	}

	if allLose {
		return "\033[31mYOU LOSE\033[0m"
	}
	if allPush {
		return "PUSH"
	}
	if anyWin || anyBlackjack {
		return "\033[32m🏆 YOU WIN!\033[0m"
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