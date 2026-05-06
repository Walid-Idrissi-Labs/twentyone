package ui

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/twentyone/twentyone/anim"
	"github.com/twentyone/twentyone/game"
	"github.com/twentyone/twentyone/styles"
)

type Model struct {
	width       int
	height      int
	screen      Screen
	game        *game.Game
	betInput    textinput.Model
	currentBet  int
	minBet      int
	maxBet      int
	noSplash    bool
	anim        *anim.Manager
	roundCount  int
	resultFlash bool
	resultTicks int
}

type StartRoundMsg     struct{}
type DealNextCardMsg   struct{}
type DealerTickMsg     struct{}
type ResolveRoundMsg   struct{}
type AnimTickMsg       struct{}
type QuitMsg           struct{}

func New(startBalance, minBet, maxBet int, noSplash bool) *Model {
	ti := textinput.New()
	ti.Focus()
	ti.Placeholder = fmt.Sprintf("$%d", min(minBet, startBalance))
	ti.Width = 15
	ti.SetValue(fmt.Sprintf("%d", min(minBet, startBalance)))

	g := game.NewGame(startBalance)

	m := &Model{
		width:      80,
		height:     24,
		screen:     ScreenWelcome,
		game:       g,
		betInput:   ti,
		currentBet: min(minBet, startBalance),
		minBet:     minBet,
		maxBet:     maxBet,
		noSplash:   noSplash,
		anim:       anim.NewManager(),
		roundCount: 0,
	}

	if noSplash {
		m.screen = ScreenBet
	}

	return m
}

func (m *Model) Init() tea.Cmd {
	styles.EnsureInit()
	if m.noSplash {
		return nil
	}
	return tea.Tick(1500*time.Millisecond, func(time.Time) tea.Msg {
		return StartRoundMsg{}
	})
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		if m.width < styles.MinTermWidth || m.height < styles.MinTermHeight {
			m.screen = ScreenTooSmall
		} else if m.screen == ScreenTooSmall {
			m.screen = ScreenBet
		}
		return m, nil

	case tea.KeyMsg:
		return m.handleKey(msg)

	case AnimTickMsg:
		m.anim.Tick()
		return m, tea.Tick(150*time.Millisecond, func(time.Time) tea.Msg { return AnimTickMsg{} })

	case StartRoundMsg:
		if m.screen == ScreenWelcome {
			m.screen = ScreenBet
		}
		return m, nil

	case DealNextCardMsg:
		return m.handleDealNextCard()

	case DealerTickMsg:
		return m.handleDealerTick()

	case ResolveRoundMsg:
		m.game.Resolve()
		m.screen = ScreenResult
		m.anim.Trigger("result", 3)
		return m, tea.Tick(150*time.Millisecond, func(time.Time) tea.Msg { return AnimTickMsg{} })
	}

	return m, nil
}

func (m *Model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch m.screen {
	case ScreenWelcome:
		if msg.Type == tea.KeyCtrlC {
			return m, tea.Quit
		}
		m.screen = ScreenBet
		return m, nil

	case ScreenBet:
		return m.handleBetKey(msg)

	case ScreenTable:
		return m.handleTableKey(msg)

	case ScreenResult:
		return m.handleResultKey(msg)

	case ScreenSummary:
		if msg.Type == tea.KeyEnter || msg.Type == tea.KeySpace || msg.Type == tea.KeyCtrlC {
			return m, tea.Quit
		}
		return m, nil
	}

	return m, nil
}

func (m *Model) handleBetKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyCtrlC, tea.KeyEsc:
		m.screen = ScreenSummary
		return m, nil
	case tea.KeyEnter:
		if m.currentBet >= m.minBet && m.currentBet <= m.game.Balance {
			return m.startRound()
		}
		return m, nil
	case tea.KeyLeft:
		m.currentBet = adjustBet(m.currentBet, -25)
		m.betInput.SetValue(fmt.Sprintf("%d", m.currentBet))
		return m, nil
	case tea.KeyRight:
		m.currentBet = adjustBet(m.currentBet, 25)
		m.betInput.SetValue(fmt.Sprintf("%d", m.currentBet))
		return m, nil
	case tea.KeyUp:
		m.currentBet = adjustBet(m.currentBet, 5)
		m.betInput.SetValue(fmt.Sprintf("%d", m.currentBet))
		return m, nil
	case tea.KeyDown:
		m.currentBet = adjustBet(m.currentBet, -5)
		m.betInput.SetValue(fmt.Sprintf("%d", m.currentBet))
		return m, nil
	case tea.KeyBackspace:
		s := fmt.Sprintf("%d", m.currentBet)
		if len(s) > 1 {
			s = s[:len(s)-1]
		} else {
			s = "0"
		}
		var n int
		fmt.Sscanf(s, "%d", &n)
		m.currentBet = adjustBet(n, 0)
		m.betInput.SetValue(fmt.Sprintf("%d", m.currentBet))
		return m, nil
	case tea.KeyRunes:
		r := msg.Runes[0]
		if r >= '0' && r <= '9' {
			s := fmt.Sprintf("%d", m.currentBet)
			if s == "0" {
				s = string(r)
			} else {
				s += string(r)
			}
			var n int
			fmt.Sscanf(s, "%d", &n)
			m.currentBet = adjustBet(n, 0)
			m.betInput.SetValue(fmt.Sprintf("%d", m.currentBet))
		}
		return m, nil
	}
	return m, nil
}

func (m *Model) handleTableKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if msg.Type == tea.KeyCtrlC {
		m.screen = ScreenSummary
		return m, nil
	}
	if msg.Type == tea.KeyEsc || (msg.Type == tea.KeyRunes && msg.Runes[0] == 'q') {
		m.screen = ScreenSummary
		return m, nil
	}

	if m.game.State == game.GameInsurance {
		if msg.Type == tea.KeyRunes {
			r := msg.Runes[0]
			if r == 'y' || r == 'Y' {
				m.game.ApplyAction(game.ActionInsuranceYes)
				m.checkBlackjack()
			} else if r == 'n' || r == 'N' {
				m.game.ApplyAction(game.ActionInsuranceNo)
				m.checkBlackjack()
			}
		}
		return m, nil
	}

	actions := m.game.AvailableActions()
	if len(actions) == 0 {
		return m, nil
	}

	actionMap := map[string]game.Action{
		"h": game.ActionHit, "H": game.ActionHit,
		"s": game.ActionStand, "S": game.ActionStand,
		"d": game.ActionDouble, "D": game.ActionDouble,
		"p": game.ActionSplit, "P": game.ActionSplit,
	}

	if a, ok := actionMap[msg.String()]; ok {
		return m.applyAction(a)
	}

	return m, nil
}

func (m *Model) handleResultKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyEnter, tea.KeySpace:
		m.screen = ScreenBet
		m.roundCount++
		return m, nil
	case tea.KeyEsc, tea.KeyCtrlC:
		m.screen = ScreenSummary
		return m, nil
	}
	if msg.Type == tea.KeyRunes && (msg.Runes[0] == 'q' || msg.Runes[0] == 'Q') {
		m.screen = ScreenSummary
		return m, nil
	}
	return m, nil
}

func (m *Model) applyAction(a game.Action) (tea.Model, tea.Cmd) {
	handIdx := m.game.ActiveHandIdx
	cardCount := len(m.game.PlayerHands[handIdx].Cards)

	m.game.ApplyAction(a)

	newCardCount := len(m.game.PlayerHands[handIdx].Cards)
	if newCardCount > cardCount {
		flashID := fmt.Sprintf("player-%d-%d", handIdx, newCardCount-1)
		m.anim.Trigger(anim.FlashID(flashID), 2)
	}

	if m.game.State == game.GameDealerTurn {
		flashID := fmt.Sprintf("dealer-%d", len(m.game.DealerHand.Cards)-1)
		m.anim.Trigger(anim.FlashID(flashID), 2)
		return m, tea.Tick(400*time.Millisecond, func(time.Time) tea.Msg { return DealerTickMsg{} })
	}

	if m.game.State == game.GameDone {
		return m, tea.Tick(time.Millisecond, func(time.Time) tea.Msg { return ResolveRoundMsg{} })
	}

	return m, nil
}

func (m *Model) handleDealNextCard() (tea.Model, tea.Cmd) {
	playerCardCount := len(m.game.PlayerHands[0].Cards)
	dealerCardCount := len(m.game.DealerHand.Cards)

	hasMore := m.game.PopDealStep()

	if hasMore {
		newPlayerCount := len(m.game.PlayerHands[0].Cards)
		newDealerCount := len(m.game.DealerHand.Cards)

		if newPlayerCount > playerCardCount {
			flashID := fmt.Sprintf("player-0-%d", newPlayerCount-1)
			m.anim.Trigger(anim.FlashID(flashID), 2)
		} else if newDealerCount > dealerCardCount {
			flashID := fmt.Sprintf("dealer-%d", newDealerCount-1)
			m.anim.Trigger(anim.FlashID(flashID), 2)
		}
		return m, tea.Batch(
			tea.Tick(150*time.Millisecond, func(time.Time) tea.Msg { return AnimTickMsg{} }),
			tea.Tick(150*time.Millisecond, func(time.Time) tea.Msg { return DealNextCardMsg{} }),
		)
	}

	if m.game.DealerHand.Cards[0].Rank == game.Ace {
		m.game.State = game.GameInsurance
	} else {
		m.game.State = game.GameCheckBJ
		m.checkBlackjack()
	}

	return m, nil
}

func (m *Model) handleDealerTick() (tea.Model, tea.Cmd) {
	dealerCardCount := len(m.game.DealerHand.Cards)

	hasMore := m.game.DealerPlay()

	if hasMore {
		newDealerCount := len(m.game.DealerHand.Cards)
		if newDealerCount > dealerCardCount {
			flashID := fmt.Sprintf("dealer-%d", newDealerCount-1)
			m.anim.Trigger(anim.FlashID(flashID), 2)
		}
		return m, tea.Tick(400*time.Millisecond, func(time.Time) tea.Msg { return DealerTickMsg{} })
	}

	m.game.State = game.GameResolve
	return m, tea.Tick(time.Millisecond, func(time.Time) tea.Msg { return ResolveRoundMsg{} })
}

func (m *Model) checkBlackjack() {
	if m.game.DealerHand.IsBlackjack() {
		if m.game.PlayerHands[0].InsuranceBet > 0 {
			m.game.Balance += m.game.PlayerHands[0].InsuranceBet * 3
		}
		if m.game.PlayerHands[0].IsBlackjack() {
			m.game.Balance += m.game.CurrentBet
			m.game.Results = []game.HandResult{
				{Hand: m.game.PlayerHands[0], Result: game.ResultPush, Profit: 0},
			}
		} else {
			m.game.Results = []game.HandResult{
				{Hand: m.game.PlayerHands[0], Result: game.ResultLose, Profit: -m.game.CurrentBet},
			}
		}
		m.game.State = game.GameDone
	} else {
		if m.game.PlayerHands[0].IsBlackjack() {
			profit := int(float64(m.game.CurrentBet) * 1.5)
			m.game.Balance += m.game.CurrentBet + profit
			m.game.State = game.GameDone
			m.game.Results = []game.HandResult{
				{Hand: m.game.PlayerHands[0], Result: game.ResultBlackjack, Profit: profit},
			}
		} else {
			m.game.State = game.GamePlayerTurn
		}
	}
}

func (m *Model) startRound() (tea.Model, tea.Cmd) {
	m.roundCount++
	m.game.StartRound(m.currentBet)
	m.screen = ScreenTable

	flashID := fmt.Sprintf("player-0-%d", len(m.game.PlayerHands[0].Cards)-1)
	m.anim.Trigger(anim.FlashID(flashID), 2)

	return m, tea.Tick(300*time.Millisecond, func(time.Time) tea.Msg { return DealNextCardMsg{} })
}

func (m *Model) View() string {
	styles.EnsureInit()

	if m.width < styles.MinTermWidth || m.height < styles.MinTermHeight {
		return renderTooSmall(*m)
	}

	switch m.screen {
	case ScreenWelcome:
		return renderWelcome(*m)
	case ScreenBet:
		return renderBet(*m)
	case ScreenTable:
		return renderTable(*m)
	case ScreenResult:
		return renderResult(*m)
	case ScreenSummary:
		return renderSummary(*m)
	default:
		return renderWelcome(*m)
	}
}

func adjustBet(current, delta int) int {
	newBet := current + delta
	if newBet < 1 {
		newBet = 1
	}
	return newBet
}