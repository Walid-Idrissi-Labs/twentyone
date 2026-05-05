package game

type HandStatus int

const (
	StatusActive HandStatus = iota
	StatusStood
	StatusBust
	StatusBlackjack
	StatusSplit
)

type Hand struct {
	Cards        []Card
	Status       HandStatus
	Bet          int
	IsDoubled    bool
	IsSplitAce   bool
	InsuranceBet int
}

func (h *Hand) Score() int {
	sum := 0
	aces := 0
	for _, c := range h.Cards {
		if c.Rank == Ace {
			aces++
			sum += 11
		} else {
			sum += c.Value()
		}
	}
	for sum > 21 && aces > 0 {
		sum -= 10
		aces--
	}
	return sum
}

func (h *Hand) IsSoft() bool {
	sum := 0
	aces := 0
	for _, c := range h.Cards {
		if c.Rank == Ace {
			aces++
			sum += 11
		} else {
			sum += c.Value()
		}
	}
	for sum > 21 && aces > 0 {
		sum -= 10
		aces--
	}
	return aces > 0 && sum <= 21
}

func (h *Hand) IsBust() bool {
	return h.Score() > 21
}

func (h *Hand) IsBlackjack() bool {
	if len(h.Cards) != 2 {
		return false
	}
	hasAce := false
	hasTenValue := false
	for _, c := range h.Cards {
		if c.Rank == Ace {
			hasAce = true
		}
		if c.IsTenValue() {
			hasTenValue = true
		}
	}
	return hasAce && hasTenValue
}

func (h *Hand) CanSplit(splitCount int) bool {
	if len(h.Cards) != 2 {
		return false
	}
	if splitCount >= 3 {
		return false
	}
	c1 := h.Cards[0]
	c2 := h.Cards[1]
	return c1.Rank == c2.Rank || (c1.IsTenValue() && c2.IsTenValue())
}

func (h *Hand) CanDouble() bool {
	return len(h.Cards) == 2 && !h.IsSplitAce
}

func (h *Hand) CanHit() bool {
	return h.Status == StatusActive && !h.IsSplitAce
}

func (h *Hand) AddCard(c Card) {
	h.Cards = append(h.Cards, c)
	if h.IsBust() {
		h.Status = StatusBust
	}
}

func (h *Hand) ScoreString() string {
	if h.IsBlackjack() {
		return "Blackjack"
	}
	score := h.Score()
	if h.IsBust() {
		return "Bust"
	}
	if h.IsSoft() {
		return "Soft " + itoa(score)
	}
	return itoa(score)
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	if n < 0 {
		return "-" + itoa(-n)
	}
	digits := []byte{}
	for n > 0 {
		digits = append([]byte{byte('0' + n%10)}, digits...)
		n /= 10
	}
	return string(digits)
}