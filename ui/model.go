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
	reshuffling bool
	roundCount  int
	buttonAreas []ButtonArea
	resultFlash bool
	resultTicks int
}

type ButtonArea struct {
	Label  string
	X, Y   int
	W, H   int
	Action interface{}
}

type StartRoundMsg struct{}

type DealNextCardMsg struct{}

type InsuranceResponseMsg bool

type PlayerActionMsg game.Action

type DealerTickMsg struct{}

type ResolveRoundMsg struct{}

type NewBetMsg int

type QuitMsg struct{}

type AnimTickMsg struct{}

func New(startBalance, minBet, maxBet int, noSplash bool) *Model {
	ti := textinput.New()
	ti.Focus()
	ti.Placeholder = fmt.Sprintf("$%d", min(minBet, startBalance))
	ti.Width = 20
	ti.SetValue(fmt.Sprintf("%d", min(minBet, startBalance)))

	g := game.NewGame(startBalance)

	m := &Model{
		width:       80,
		height:      24,
		screen:      ScreenWelcome,
		game:        g,
		betInput:    ti,
		currentBet:  min(minBet, startBalance),
		minBet:      minBet,
		maxBet:      maxBet,
		noSplash:    noSplash,
		anim:        anim.NewManager(),
		reshuffling: false,
		roundCount:  0,
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
	return tea.Tick(time.Second+time.Millisecond*500, func(time.Time) tea.Msg {
		return StartRoundMsg{}
	})
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	m.buttonAreas = nil

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		if m.width < styles.MinTermWidth || m.height < styles.MinTermHeight {
			m.screen = ScreenTooSmall
		}
		return m, nil

	case tea.KeyMsg:
		return m.handleKeyMsg(msg)

	case tea.MouseMsg:
		return m.handleMouseMsg(msg)

	case AnimTickMsg:
		expired := m.anim.Tick()
		for _, id := range expired {
			if id == "result" {
				m.resultFlash = false
			}
		}
		if m.resultFlash && m.resultTicks > 0 {
			m.resultTicks--
		}
		return m, tea.Tick(150*time.Millisecond, func(t time.Time) tea.Msg { return AnimTickMsg{} })

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
		m.resultFlash = true
		m.resultTicks = 3
		return m, tea.Tick(150*time.Millisecond, func(t time.Time) tea.Msg { return AnimTickMsg{} })
	}

	return m, nil
}

func (m *Model) handleKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
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
		m.currentBet = adjustBet(m.currentBet, -25, m.minBet, m.maxBet, m.game.Balance)
		m.betInput.SetValue(fmt.Sprintf("%d", m.currentBet))
		return m, nil
	case tea.KeyRight:
		m.currentBet = adjustBet(m.currentBet, 25, m.minBet, m.maxBet, m.game.Balance)
		m.betInput.SetValue(fmt.Sprintf("%d", m.currentBet))
		return m, nil
	case tea.KeyUp:
		m.currentBet = adjustBet(m.currentBet, 5, m.minBet, m.maxBet, m.game.Balance)
		m.betInput.SetValue(fmt.Sprintf("%d", m.currentBet))
		return m, nil
	case tea.KeyDown:
		m.currentBet = adjustBet(m.currentBet, -5, m.minBet, m.maxBet, m.game.Balance)
		m.betInput.SetValue(fmt.Sprintf("%d", m.currentBet))
		return m, nil
	case tea.KeyBackspace:
		currentStr := fmt.Sprintf("%d", m.currentBet)
		if len(currentStr) > 1 {
			currentStr = currentStr[:len(currentStr)-1]
		} else {
			currentStr = "0"
		}
		var newBet int
		fmt.Sscanf(currentStr, "%d", &newBet)
		m.currentBet = adjustBet(newBet, 0, m.minBet, m.maxBet, m.game.Balance)
		m.betInput.SetValue(fmt.Sprintf("%d", m.currentBet))
		return m, nil
	case tea.KeyRunes:
		r := msg.Runes[0]
		if r >= '0' && r <= '9' {
			currentStr := fmt.Sprintf("%d", m.currentBet)
			if currentStr == "0" {
				currentStr = string(r)
			} else {
				currentStr += string(r)
			}
			var newBet int
			fmt.Sscanf(currentStr, "%d", &newBet)
			m.currentBet = adjustBet(newBet, 0, m.minBet, m.maxBet, m.game.Balance)
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

	if msg.Type == tea.KeyEsc {
		m.screen = ScreenSummary
		return m, nil
	}

	if msg.Type == tea.KeyRunes && (msg.Runes[0] == 'q' || msg.Runes[0] == 'Q') {
		m.screen = ScreenSummary
		return m, nil
	}

	if m.game.State == game.GameInsurance {
		if msg.Type == tea.KeyRunes {
			r := msg.Runes[0]
			if r == 'y' || r == 'Y' {
				return m.handleInsurance(true)
			}
			if r == 'n' || r == 'N' {
				return m.handleInsurance(false)
			}
		}
		return m, nil
	}

	actions := m.game.AvailableActions()
	if len(actions) == 0 {
		return m, nil
	}

	actionMap := make(map[string]game.Action)
	for _, a := range actions {
		switch a {
		case game.ActionHit:
			actionMap["h"] = game.ActionHit
			actionMap["H"] = game.ActionHit
		case game.ActionStand:
			actionMap["s"] = game.ActionStand
			actionMap["S"] = game.ActionStand
		case game.ActionDouble:
			actionMap["d"] = game.ActionDouble
			actionMap["D"] = game.ActionDouble
		case game.ActionSplit:
			actionMap["p"] = game.ActionSplit
			actionMap["P"] = game.ActionSplit
		}
	}

	if a, ok := actionMap[msg.String()]; ok {
		return m.applyPlayerAction(a)
	}

	return m, nil
}

func (m *Model) handleResultKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyEnter, tea.KeySpace:
		m.screen = ScreenBet
		m.roundCount++
		return m, nil
	case tea.KeyEsc:
		m.screen = ScreenSummary
		return m, nil
	}
	if msg.Type == tea.KeyRunes && (msg.Runes[0] == 'q' || msg.Runes[0] == 'Q') {
		m.screen = ScreenSummary
		return m, nil
	}
	return m, nil
}

func (m *Model) handleMouseMsg(msg tea.MouseMsg) (tea.Model, tea.Cmd) {
	x, y := msg.X, msg.Y

	switch m.screen {
	case ScreenBet:
		if y >= m.height-3 && y <= m.height-1 {
			if x >= 14 && x <= 26 {
				m.currentBet = adjustBet(m.currentBet, -25, m.minBet, m.maxBet, m.game.Balance)
				m.betInput.SetValue(fmt.Sprintf("%d", m.currentBet))
			} else if x >= 28 && x <= 40 {
				m.currentBet = adjustBet(m.currentBet, -5, m.minBet, m.maxBet, m.game.Balance)
				m.betInput.SetValue(fmt.Sprintf("%d", m.currentBet))
			} else if x >= 42 && x <= 54 {
				m.currentBet = adjustBet(m.currentBet, 5, m.minBet, m.maxBet, m.game.Balance)
				m.betInput.SetValue(fmt.Sprintf("%d", m.currentBet))
			} else if x >= 56 && x <= 68 {
				m.currentBet = adjustBet(m.currentBet, 25, m.minBet, m.maxBet, m.game.Balance)
				m.betInput.SetValue(fmt.Sprintf("%d", m.currentBet))
			}
		}
	}

	return m, nil
}

func (m *Model) applyPlayerAction(a game.Action) (tea.Model, tea.Cmd) {
	handIdx := m.game.ActiveHandIdx
	cardCount := len(m.game.PlayerHands[handIdx].Cards)

	m.game.ApplyAction(a)

	if a == game.ActionHit || a == game.ActionDouble || a == game.ActionSplit {
		flashID := fmt.Sprintf("player-%d-%d", handIdx, cardCount)
		m.anim.Trigger(anim.FlashID(flashID), 2)
	}

	if m.game.State == game.GameDealerTurn {
		return m, tea.Tick(400*time.Millisecond, func(t time.Time) tea.Msg { return DealerTickMsg{} })
	}

	if m.game.State == game.GameDone {
		cmd := tea.Tick(time.Millisecond, func(t time.Time) tea.Msg { return ResolveRoundMsg{} })
		return m, cmd
	}

	return m, nil
}

func (m *Model) handleInsurance(accept bool) (tea.Model, tea.Cmd) {
	if accept {
		m.game.ApplyAction(game.ActionInsuranceYes)
	} else {
		m.game.ApplyAction(game.ActionInsuranceNo)
	}

	if m.game.State == game.GameCheckBJ {
		m.checkBlackjack()
	}

	if m.game.State == game.GamePlayerTurn {
		return m, nil
	}

	if m.game.State == game.GameDone {
		cmd := tea.Tick(time.Millisecond, func(t time.Time) tea.Msg { return ResolveRoundMsg{} })
		return m, cmd
	}

	return m, nil
}

func (m *Model) handleDealNextCard() (tea.Model, tea.Cmd) {
	playerCardCount := len(m.game.PlayerHands[0].Cards)
	dealerCardCount := len(m.game.DealerHand.Cards)

	hasMore := m.game.PopDealStep()

	if hasMore {
		if len(m.game.PlayerHands) > 0 && len(m.game.PlayerHands[0].Cards) > playerCardCount {
			flashID := fmt.Sprintf("player-0-%d", len(m.game.PlayerHands[0].Cards)-1)
			m.anim.Trigger(anim.FlashID(flashID), 2)
		} else if len(m.game.DealerHand.Cards) > dealerCardCount {
			flashID := fmt.Sprintf("dealer-%d", len(m.game.DealerHand.Cards)-1)
			m.anim.Trigger(anim.FlashID(flashID), 2)
		}
		return m, tea.Batch(
			tea.Tick(120*time.Millisecond, func(t time.Time) tea.Msg { return AnimTickMsg{} }),
			tea.Tick(120*time.Millisecond, func(t time.Time) tea.Msg { return DealNextCardMsg{} }),
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

	if hasMore && len(m.game.DealerHand.Cards) > dealerCardCount {
		flashID := fmt.Sprintf("dealer-%d", len(m.game.DealerHand.Cards)-1)
		m.anim.Trigger(anim.FlashID(flashID), 2)
	}

	if hasMore {
		return m, tea.Tick(400*time.Millisecond, func(t time.Time) tea.Msg { return DealerTickMsg{} })
	}

	m.game.State = game.GameResolve
	return m, tea.Tick(time.Millisecond, func(t time.Time) tea.Msg { return ResolveRoundMsg{} })
}

func (m *Model) checkBlackjack() {
	if m.game.DealerHand.IsBlackjack() {
		if m.game.PlayerHands[0].InsuranceBet > 0 {
			m.game.Balance += m.game.PlayerHands[0].InsuranceBet * 3
		}
		if m.game.PlayerHands[0].IsBlackjack() {
			m.game.Balance += m.game.CurrentBet
			m.game.Results = []game.HandResult{
				{
					Hand:   m.game.PlayerHands[0],
					Result: game.ResultPush,
					Profit: 0,
				},
			}
		} else {
			m.game.Results = []game.HandResult{
				{
					Hand:   m.game.PlayerHands[0],
					Result: game.ResultLose,
					Profit: -m.game.CurrentBet,
				},
			}
		}
		m.game.State = game.GameDone
	} else {
		if m.game.PlayerHands[0].IsBlackjack() {
			profit := int(float64(m.game.CurrentBet) * 1.5)
			m.game.Balance += m.game.CurrentBet + profit
			m.game.State = game.GameDone
			m.game.Results = []game.HandResult{
				{
					Hand:   m.game.PlayerHands[0],
					Result: game.ResultBlackjack,
					Profit: profit,
				},
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

	flashID := "dealer-0"
	m.anim.Trigger(anim.FlashID(flashID), 2)

	return m, tea.Tick(400*time.Millisecond, func(t time.Time) tea.Msg { return DealNextCardMsg{} })
}

func (m *Model) View() string {
	if m.width < styles.MinTermWidth || m.height < styles.MinTermHeight {
		return renderTooSmall(*m)
	}

	styles.EnsureInit()

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

func adjustBet(current, delta, minBet, maxBet, balance int) int {
	newBet := current + delta
	if newBet < minBet {
		newBet = minBet
	}
	if newBet > balance {
		newBet = balance
	}
	if maxBet > 0 && newBet > maxBet {
		newBet = maxBet
	}
	return newBet
}