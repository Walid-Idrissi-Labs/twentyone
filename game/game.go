package game

type Action int

const (
	ActionHit Action = iota
	ActionStand
	ActionDouble
	ActionSplit
	ActionInsuranceYes
	ActionInsuranceNo
)

type RoundResult int

const (
	ResultWin RoundResult = iota
	ResultLose
	ResultPush
	ResultBlackjack
)

type HandResult struct {
	Hand   *Hand
	Result RoundResult
	Profit int
}

type GameState int

const (
	GameIdle GameState = iota
	GameDealing
	GameInsurance
	GameCheckBJ
	GamePlayerTurn
	GameDealerTurn
	GameResolve
	GameDone
)

type Game struct {
	Shoe          *Shoe
	PlayerHands   []*Hand
	DealerHand    *Hand
	ActiveHandIdx int
	SplitCount    int
	State         GameState
	Results       []HandResult
	SessionStart  int
	Balance       int
	CurrentBet    int
	DealQueue     []dealStep
}

type dealStep struct {
	Target string
	Hidden bool
}

func NewGame(startBalance int) *Game {
	return &Game{
		Shoe:         NewShoe(),
		DealerHand:   &Hand{},
		SessionStart:  startBalance,
		Balance:      startBalance,
		State:         GameIdle,
		ActiveHandIdx: 0,
		SplitCount:    0,
	}
}

func (g *Game) StartRound(bet int) error {
	if g.Shoe.NeedsReshuffle() {
		g.Shoe.Reshuffle()
	}

	g.PlayerHands = []*Hand{}
	g.DealerHand = &Hand{}
	g.Results = nil
	g.ActiveHandIdx = 0
	g.SplitCount = 0
	g.CurrentBet = bet
	g.Balance -= bet

	playerHand := &Hand{Bet: bet}
	g.PlayerHands = append(g.PlayerHands, playerHand)

	g.DealQueue = []dealStep{
		{Target: "player", Hidden: false},
		{Target: "dealer", Hidden: false},
		{Target: "player", Hidden: false},
		{Target: "dealer", Hidden: true},
	}

	g.State = GameDealing
	return nil
}

func (g *Game) PopDealStep() bool {
	if len(g.DealQueue) == 0 {
		return false
	}

	step := g.DealQueue[0]
	g.DealQueue = g.DealQueue[1:]

	var card Card
	if step.Hidden {
		card = g.Shoe.DealHidden()
	} else {
		card = g.Shoe.Deal()
	}

	if step.Target == "player" {
		g.PlayerHands[0].AddCard(card)
	} else {
		g.DealerHand.AddCard(card)
	}

	return len(g.DealQueue) > 0
}

func (g *Game) ApplyAction(a Action) error {
	if g.State != GamePlayerTurn {
		return nil
	}

	hand := g.PlayerHands[g.ActiveHandIdx]

	switch a {
	case ActionHit:
		if !hand.CanHit() {
			return nil
		}
		card := g.Shoe.Deal()
		hand.AddCard(card)
		if hand.IsBust() {
			g.advanceHand()
		}

	case ActionStand:
		hand.Status = StatusStood
		g.advanceHand()

	case ActionDouble:
		if !hand.CanDouble() || g.Balance < hand.Bet {
			return nil
		}
		g.Balance -= hand.Bet
		hand.Bet *= 2
		hand.IsDoubled = true
		card := g.Shoe.Deal()
		hand.AddCard(card)
		if hand.IsBust() {
			hand.Status = StatusBust
		} else {
			hand.Status = StatusStood
		}
		g.advanceHand()

	case ActionSplit:
		if !hand.CanSplit(g.SplitCount) || g.Balance < g.CurrentBet {
			return nil
		}

		c1 := hand.Cards[0]
		c2 := hand.Cards[1]

		hand1 := &Hand{Bet: g.CurrentBet, Status: StatusActive}
		hand2 := &Hand{Bet: g.CurrentBet, Status: StatusActive}

		hand1.AddCard(c1)
		hand2.AddCard(c2)

		if c1.Rank == Ace {
			hand1.IsSplitAce = true
			hand2.IsSplitAce = true
		}

		g.Balance -= g.CurrentBet
		g.SplitCount++

		newCard1 := g.Shoe.Deal()
		newCard2 := g.Shoe.Deal()
		hand1.AddCard(newCard1)
		hand2.AddCard(newCard2)

		if hand1.IsSplitAce {
			hand1.Status = StatusStood
			hand2.Status = StatusStood
		}

		g.PlayerHands[g.ActiveHandIdx] = hand1
		g.PlayerHands = append(g.PlayerHands[:g.ActiveHandIdx+1], hand2)

		if hand1.IsSplitAce && hand2.IsSplitAce {
			g.advanceHand()
		}

	case ActionInsuranceYes:
		if g.DealerHand.Cards[0].Rank != Ace {
			return nil
		}
		insuranceBet := g.CurrentBet / 2
		if insuranceBet < 1 {
			insuranceBet = 1
		}
		g.Balance -= insuranceBet
		g.PlayerHands[0].InsuranceBet = insuranceBet
		g.State = GameCheckBJ
		return nil

	case ActionInsuranceNo:
		g.State = GameCheckBJ
		return nil
	}

	return nil
}

func (g *Game) advanceHand() {
	g.ActiveHandIdx++
	if g.ActiveHandIdx >= len(g.PlayerHands) {
		allBust := true
		for _, h := range g.PlayerHands {
			if h.Status != StatusBust {
				allBust = false
				break
			}
		}
		if allBust {
			g.State = GameDone
			g.Results = make([]HandResult, len(g.PlayerHands))
			for i, h := range g.PlayerHands {
				g.Results[i] = HandResult{Hand: h, Result: ResultLose, Profit: -h.Bet}
			}
		} else {
			g.State = GameDealerTurn
		}
	}
}

func (g *Game) DealerPlay() bool {
	if g.State != GameDealerTurn {
		return false
	}

	if len(g.DealerHand.Cards) == 1 {
		g.DealerHand.Cards[0].FaceUp = true
		return true
	}

	if g.DealerHand.Cards[1].FaceUp == false {
		g.DealerHand.Cards[1].FaceUp = true
		return true
	}

	score := g.DealerHand.Score()
	if score >= 17 {
		return false
	}

	card := g.Shoe.Deal()
	g.DealerHand.AddCard(card)
	return g.DealerHand.Score() < 17
}

func (g *Game) Resolve() {
	g.Results = make([]HandResult, len(g.PlayerHands))

	for i, hand := range g.PlayerHands {
		result := ResultLose
		profit := -hand.Bet

		if hand.Status == StatusBust {
			result = ResultLose
			profit = -hand.Bet
		} else if hand.IsBlackjack() && !g.DealerHand.IsBlackjack() && !hand.IsSplitAce {
			result = ResultBlackjack
			profit = int(float64(hand.Bet) * 1.5)
			g.Balance += hand.Bet + profit
		} else if g.DealerHand.IsBust() {
			result = ResultWin
			profit = hand.Bet
			g.Balance += hand.Bet * 2
		} else if hand.Score() > g.DealerHand.Score() {
			result = ResultWin
			profit = hand.Bet
			g.Balance += hand.Bet * 2
		} else if hand.Score() == g.DealerHand.Score() {
			result = ResultPush
			profit = 0
			g.Balance += hand.Bet
		}

		g.Results[i] = HandResult{Hand: hand, Result: result, Profit: profit}
	}
}

func (g *Game) AvailableActions() []Action {
	if g.State != GamePlayerTurn {
		return nil
	}

	hand := g.PlayerHands[g.ActiveHandIdx]
	actions := []Action{}

	if hand.CanHit() {
		actions = append(actions, ActionHit)
	}
	actions = append(actions, ActionStand)

	if hand.CanDouble() && g.Balance >= hand.Bet {
		actions = append(actions, ActionDouble)
	}

	if hand.CanSplit(g.SplitCount) && g.Balance >= g.CurrentBet {
		actions = append(actions, ActionSplit)
	}

	return actions
}

func (g *Game) SessionProfit() int {
	return g.Balance - g.SessionStart
}

func (g *Game) InsuranceResolution() bool {
	if g.DealerHand.IsBlackjack() {
		insuranceBet := g.PlayerHands[0].InsuranceBet
		if insuranceBet > 0 {
			g.Balance += insuranceBet * 3
		}
		return true
	}
	return false
}

func (g *Game) CheckPlayerBlackjack() bool {
	return g.PlayerHands[0].IsBlackjack()
}