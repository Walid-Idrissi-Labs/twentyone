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
	var sb strings.Builder
	sb.WriteString(renderHUD(m))
	sb.WriteString("\n\n")
	sb.WriteString(renderDealerArea(m))
	sb.WriteString("\n\n")
	sb.WriteString(renderPlayerArea(m))
	sb.WriteString("\n\n")
	sb.WriteString(renderActionBar(m))

	bg := lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Top, sb.String())
	return styles.GetStyleBackground().Render(bg)
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

	used := len(balanceStr) + len(betStr) + len(shoeStr) + len(roundStr) + 12
	padding := m.width - used
	if padding < 3 {
		padding = 3
	}
	centerPadding := padding / 2

	left := balanceStr + "    " + betStr
	center := strings.Repeat(" ", centerPadding) + "twentyone"
	right := shoeStr + "    " + roundStr

	return styles.GetHUDStyle().Render(left + center + right)
}

func renderDealerArea(m Model) string {
	hand := m.game.DealerHand
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("  DEALER"))
	if len(hand.Cards) > 0 {
		if hand.Cards[0].FaceUp {
			sb.WriteString(fmt.Sprintf("                                           Score: %s", hand.ScoreString()))
		} else {
			sb.WriteString(fmt.Sprintf("                                           Score: ?"))
		}
	}
	sb.WriteString("\n\n")

	cards := make([]string, 0)
	for i, c := range hand.Cards {
		flashID := fmt.Sprintf("dealer-%d", i)
		highlighted := m.anim.IsActive(anim.FlashID(flashID))
		cards = append(cards, renderCard(c, highlighted, false))
	}
	if len(cards) == 0 {
		cards = append(cards, renderCard(game.Card{FaceUp: false}, false, false))
	}

	for i := 0; i < 5; i++ {
		line := "    "
		for j, card := range cards {
			cardLines := strings.Split(card, "\n")
			if i < len(cardLines) {
				line += cardLines[i]
				if j < len(cards)-1 {
					line += " "
				}
			}
		}
		sb.WriteString(line + "\n")
	}

	return sb.String()
}

func renderPlayerArea(m Model) string {
	var sb strings.Builder
	hands := m.game.PlayerHands
	activeIdx := m.game.ActiveHandIdx

	if len(hands) == 0 {
		return sb.String()
	}

	if len(hands) == 1 {
		hand := hands[0]
		active := m.game.State == game.GamePlayerTurn && activeIdx == 0

		sb.WriteString(fmt.Sprintf("  YOU"))
		if len(hand.Cards) > 0 {
			sb.WriteString(fmt.Sprintf("                                           Score: %s", hand.ScoreString()))
		}
		sb.WriteString("\n\n")

		cards := make([]string, 0)
		for i, c := range hand.Cards {
			flashID := fmt.Sprintf("player-0-%d", i)
			highlighted := m.anim.IsActive(anim.FlashID(flashID))
			dimmed := !active && m.game.State == game.GamePlayerTurn
			cards = append(cards, renderCard(c, highlighted, dimmed))
		}

		for i := 0; i < 5; i++ {
			line := "    "
			for j, card := range cards {
				cardLines := strings.Split(card, "\n")
				if i < len(cardLines) {
					line += cardLines[i]
					if j < len(cards)-1 {
						line += " "
					}
				}
			}
			sb.WriteString(line + "\n")
		}
	} else {
		sb.WriteString(fmt.Sprintf("  YOU\n\n"))

		handLines := make([][]string, len(hands))
		for hIdx, hand := range hands {
			active := m.game.State == game.GamePlayerTurn && activeIdx == hIdx
			cards := make([]string, 0)
			for i, c := range hand.Cards {
				flashID := fmt.Sprintf("player-%d-%d", hIdx, i)
				highlighted := m.anim.IsActive(anim.FlashID(flashID))
				dimmed := !active
				cards = append(cards, renderCard(c, highlighted, dimmed))
			}
			handLines[hIdx] = cards
		}

		maxCards := 0
		for _, cards := range handLines {
			if len(cards) > maxCards {
				maxCards = len(cards)
			}
		}

		for row := 0; row < 5; row++ {
			line := "    "
			for hIdx, cards := range handLines {
				if hIdx > 0 {
					line += "   "
				}
				active := m.game.State == game.GamePlayerTurn && activeIdx == hIdx
				if row == 0 {
					prefix := "  "
					if active {
						prefix = "▶ "
					}
					hand := hands[hIdx]
					handInfo := fmt.Sprintf("%sHand %d $%d   %s", prefix, hIdx+1, hand.Bet, hand.ScoreString())
					line += handInfo
					spaceNeeded := 40 - len(handInfo)
					if spaceNeeded > 0 {
						line += strings.Repeat(" ", spaceNeeded)
					}
				} else if row-1 < len(cards) {
					cardLines := strings.Split(cards[row-1], "\n")
					if len(cardLines) > 0 {
						line += cardLines[0]
					}
				}
			}
			sb.WriteString(line + "\n")
		}
	}

	return sb.String()
}

func renderActionBar(m Model) string {
	state := m.game.State

	if state == game.GameDealing {
		return styles.GetHUDStyle().Render("  Dealing…")
	}

	if state == game.GameInsurance {
		return styles.GetHUDStyle().Render("  Insurance offered.   [ Take  Y ]   [ Decline  N ]")
	}

	if state == game.GameDealerTurn {
		return styles.GetHUDStyle().Render("  Dealer playing…")
	}

	actions := m.game.AvailableActions()
	if len(actions) == 0 {
		return styles.GetHUDStyle().Render("")
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

	return styles.GetHUDStyle().Render(sb.String())
}