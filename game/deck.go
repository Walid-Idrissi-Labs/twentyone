package game

import (
	"math/rand/v2"
	"time"
)

type Shoe struct {
	Cards     []Card
	DealIndex int
}

func NewShoe() *Shoe {
	cards := make([]Card, 0, 312)
	for d := 0; d < 6; d++ {
		for s := Spades; s <= Clubs; s++ {
			for r := Two; r <= Ace; r++ {
				cards = append(cards, Card{Rank: r, Suit: s, FaceUp: true})
			}
		}
	}

	s := rand.New(rand.NewPCG(uint64(time.Now().UnixNano()), 0))
	for i := len(cards) - 1; i > 0; i-- {
		j := s.IntN(i + 1)
		cards[i], cards[j] = cards[j], cards[i]
	}

	return &Shoe{
		Cards:     cards,
		DealIndex: 0,
	}
}

func (s *Shoe) Deal() Card {
	card := s.Cards[s.DealIndex]
	card.FaceUp = true
	s.DealIndex++
	return card
}

func (s *Shoe) DealHidden() Card {
	card := s.Cards[s.DealIndex]
	card.FaceUp = false
	s.DealIndex++
	return card
}

func (s *Shoe) NeedsReshuffle() bool {
	return len(s.Cards)-s.DealIndex < 52
}

func (s *Shoe) Reshuffle() {
	newCards := make([]Card, 0, 312)
	for d := 0; d < 6; d++ {
		for s := Spades; s <= Clubs; s++ {
			for r := Two; r <= Ace; r++ {
				newCards = append(newCards, Card{Rank: r, Suit: s, FaceUp: true})
			}
		}
	}

	r := rand.New(rand.NewPCG(uint64(time.Now().UnixNano()), 0))
	for i := len(newCards) - 1; i > 0; i-- {
		j := r.IntN(i + 1)
		newCards[i], newCards[j] = newCards[j], newCards[i]
	}

	s.Cards = newCards
	s.DealIndex = 0
}

func (s *Shoe) CardsRemaining() int {
	return len(s.Cards) - s.DealIndex
}