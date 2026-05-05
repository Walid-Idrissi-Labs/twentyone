package game

type Suit int

const (
	Spades Suit = iota
	Hearts
	Diamonds
	Clubs
)

type Rank int

const (
	Two Rank = iota + 2
	Three
	Four
	Five
	Six
	Seven
	Eight
	Nine
	Ten
	Jack
	Queen
	King
	Ace
)

type Card struct {
	Rank   Rank
	Suit   Suit
	FaceUp bool
}

func (c Card) Value() int {
	switch c.Rank {
	case Ace:
		return 11
	case Jack, Queen, King:
		return 10
	default:
		return int(c.Rank)
	}
}

func (c Card) String() string {
	rankStr := c.RankString()
	return rankStr + c.SuitSymbol()
}

func (c Card) RankString() string {
	switch c.Rank {
	case Ace:
		return "A"
	case Jack:
		return "J"
	case Queen:
		return "Q"
	case King:
		return "K"
	case Ten:
		return "10"
	default:
		return string(rune('0' + c.Rank))
	}
}

func (c Card) SuitSymbol() string {
	switch c.Suit {
	case Spades:
		return "♠"
	case Hearts:
		return "♥"
	case Diamonds:
		return "♦"
	case Clubs:
		return "♣"
	}
	return "?"
}

func (c Card) IsRedSuit() bool {
	return c.Suit == Hearts || c.Suit == Diamonds
}

func (c Card) IsTenValue() bool {
	return c.Rank == Ten || c.Rank == Jack || c.Rank == Queen || c.Rank == King
}