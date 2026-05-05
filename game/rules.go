package game

func CanPlayerInsure(dealerUpCard Card) bool {
	return dealerUpCard.Rank == Ace
}

func IsNaturalBlackjack(hand *Hand) bool {
	return hand.IsBlackjack()
}

func DealerShouldHit(hand *Hand) bool {
	return hand.Score() < 17
}

func ComputePayout(hand *Hand, dealerHand *Hand, balance int) (result RoundResult, profit int, newBalance int) {
	if hand.Status == StatusBust {
		return ResultLose, -hand.Bet, balance
	}

	if hand.IsBlackjack() && !dealerHand.IsBlackjack() {
		profit = int(float64(hand.Bet) * 1.5)
		return ResultBlackjack, profit, balance + hand.Bet + profit
	}

	if dealerHand.IsBust() {
		if hand.Status == StatusBust {
			return ResultLose, -hand.Bet, balance
		}
		return ResultWin, hand.Bet, balance + hand.Bet*2
	}

	if hand.Score() > dealerHand.Score() {
		return ResultWin, hand.Bet, balance + hand.Bet*2
	}

	if hand.Score() == dealerHand.Score() {
		return ResultPush, 0, balance + hand.Bet
	}

	return ResultLose, -hand.Bet, balance
}