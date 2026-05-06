package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/twentyone/twentyone/anim"
	"github.com/twentyone/twentyone/game"
	"github.com/twentyone/twentyone/styles"
)

func renderTable(m Model) string {
	styles.EnsureInit()

	var sb strings.Builder

	sb.WriteString(renderHUD(m))
	sb.WriteString("\n\n")
	sb.WriteString(renderDealerArea(m))
	sb.WriteString("\n\n")
	sb.WriteString(renderPlayerArea(m))
	sb.WriteString("\n")
	sb.WriteString(renderActionBar(m))

	return styles.StyleBackground.Render(sb.String())
}

func renderHUD(m Model) string {
	balance := m.game.Balance
	bet := m.game.CurrentBet
	shoeCards := 0
	if m.game.Shoe != nil {
		shoeCards = m.game.Shoe.CardsRemaining()
	}
	round := m.roundCount + 1

	balanceStr := fmt.Sprintf("Balance: $%d", balance)
	betStr := fmt.Sprintf("Bet: $%d", bet)
	shoeStr := fmt.Sprintf("Shoe: %d cards", shoeCards)
	roundStr := fmt.Sprintf("Round %d", round)

	bar := fmt.Sprintf("%s    %s    %s    %s", balanceStr, betStr, shoeStr, roundStr)

	return lipgloss.NewStyle().
		Background(styles.ColorSurface).
		Foreground(styles.ColorText).
		Width(m.width).
		Render(bar)
}

func renderDealerArea(m Model) string {
	hand := m.game.DealerHand
	var sb strings.Builder

	dealerLabel := "DEALER"
	scoreStr := "Score: ?"
	if len(hand.Cards) > 0 && hand.Cards[0].FaceUp {
		scoreStr = "Score: " + hand.ScoreString()
	}

	sb.WriteString(fmt.Sprintf("  %s                                             %s\n\n", dealerLabel, scoreStr))

	cards := make([]string, 0)
	for i, c := range hand.Cards {
		flashID := fmt.Sprintf("dealer-%d", i)
		highlighted := m.anim.IsActive(anim.FlashID(flashID))
		cards = append(cards, renderCard(c, highlighted, false))
	}
	if len(cards) == 0 {
		cards = append(cards, renderCard(game.Card{FaceUp: false}, false, false))
	}

	for row := 0; row < 5; row++ {
		line := "    "
		for col, card := range cards {
			cardLines := strings.Split(card, "\n")
			if row < len(cardLines) {
				line += cardLines[row]
				if col < len(cards)-1 {
					line += " "
				}
			}
		}
		sb.WriteString(line + "\n")
	}

	return sb.String()
}

func renderPlayerArea(m Model) string {
	hand := m.game.PlayerHands
	activeIdx := m.game.ActiveHandIdx
	var sb strings.Builder

	if len(hand) == 0 {
		return sb.String()
	}

	if len(hand) == 1 {
		h := hand[0]
		active := m.game.State == game.GamePlayerTurn && activeIdx == 0

		scoreStr := ""
		if len(h.Cards) > 0 {
			scoreStr = "Score: " + h.ScoreString()
		}
		sb.WriteString(fmt.Sprintf("  YOU                                             %s\n\n", scoreStr))

		cards := make([]string, 0)
		for i, c := range h.Cards {
			flashID := fmt.Sprintf("player-0-%d", i)
			highlighted := m.anim.IsActive(anim.FlashID(flashID))
			dimmed := !active && m.game.State == game.GamePlayerTurn
			cards = append(cards, renderCard(c, highlighted, dimmed))
		}

		for row := 0; row < 5; row++ {
			line := "    "
			for col, card := range cards {
				cardLines := strings.Split(card, "\n")
				if row < len(cardLines) {
					line += cardLines[row]
					if col < len(cards)-1 {
						line += " "
					}
				}
			}
			sb.WriteString(line + "\n")
		}
	} else {
		sb.WriteString("  YOU\n\n")

		handInfos := make([]string, 0, len(hand))
		for i, h := range hand {
			active := m.game.State == game.GamePlayerTurn && activeIdx == i
			prefix := "  "
			if active {
				prefix = "▶ "
			}
			handInfo := fmt.Sprintf("%sHand %d $%d", prefix, i+1, h.Bet)
			if len(h.Cards) > 0 {
				handInfo += " " + h.ScoreString()
			}
			handInfos = append(handInfos, handInfo)
		}

		maxLen := 0
		for _, info := range handInfos {
			if len(info) > maxLen {
				maxLen = len(info)
			}
		}

		for i, info := range handInfos {
			padding := maxLen - len(info)
			if padding > 0 {
				info += strings.Repeat(" ", padding)
			}
			if i > 0 {
				sb.WriteString("    ")
			}
			sb.WriteString(info + "   ")
		}
		sb.WriteString("\n\n")

		for row := 0; row < 5; row++ {
			line := "    "
			for hIdx, h := range hand {
				active := m.game.State == game.GamePlayerTurn && activeIdx == hIdx
				cards := make([]string, 0)
				for i, c := range h.Cards {
					flashID := fmt.Sprintf("player-%d-%d", hIdx, i)
					highlighted := m.anim.IsActive(anim.FlashID(flashID))
					dimmed := !active
					cards = append(cards, renderCard(c, highlighted, dimmed))
				}

				for col, card := range cards {
					cardLines := strings.Split(card, "\n")
					if row < len(cardLines) {
						line += cardLines[row]
						if col < len(cards)-1 {
							line += " "
						}
					}
				}
				if hIdx < len(hand)-1 {
					line += "   "
				}
			}
			sb.WriteString(line + "\n")
		}
	}

	return sb.String()
}

func renderActionBar(m Model) string {
	state := m.game.State

	hudStyle := lipgloss.NewStyle().
		Background(styles.ColorSurface).
		Foreground(styles.ColorText).
		Width(m.width)

	if state == game.GameDealing {
		return hudStyle.Render("  Dealing...")
	}

	if state == game.GameInsurance {
		return hudStyle.Render("  Insurance offered.   [ Take  Y ]   [ Decline  N ]")
	}

	if state == game.GameDealerTurn {
		return hudStyle.Render("  Dealer playing...")
	}

	if state == game.GamePlayerTurn {
		actions := m.game.AvailableActions()
		if len(actions) == 0 {
			return hudStyle.Render("")
		}

		var sb strings.Builder
		sb.WriteString("  ")
		for _, a := range actions {
			switch a {
			case game.ActionHit:
				sb.WriteString(renderButton("Hit", "H", false, false) + "  ")
			case game.ActionStand:
				sb.WriteString(renderButton("Stand", "S", false, false) + "  ")
			case game.ActionDouble:
				sb.WriteString(renderButton("Double", "D", false, false) + "  ")
			case game.ActionSplit:
				sb.WriteString(renderButton("Split", "P", false, false) + "  ")
			}
		}

		return hudStyle.Render(sb.String())
	}

	return hudStyle.Render("")
}