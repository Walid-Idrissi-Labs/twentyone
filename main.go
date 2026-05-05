package main

import (
	"flag"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/twentyone/twentyone/ui"
)

var (
	balance   int
	noSplash  bool
	minBet    int
	maxBet    int
	startBal  int
)

func init() {
	flag.IntVar(&balance, "balance", 1000, "Starting balance in dollars")
	flag.BoolVar(&noSplash, "no-splash", false, "Skip the welcome screen")
	flag.IntVar(&minBet, "min-bet", 1, "Minimum bet")
	flag.IntVar(&maxBet, "max-bet", 0, "Maximum bet, 0 = no limit")
}

func main() {
	flag.Parse()

	if err := validateFlags(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	defer func() {
		if r := recover(); r != nil {
			logFile, _ := os.OpenFile("twentyone-crash.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
			fmt.Fprintln(logFile, r)
			logFile.Close()
			fmt.Fprintln(os.Stderr, "A fatal error occurred. See twentyone-crash.log for details.")
			os.Exit(1)
		}
	}()

	model := ui.New(balance, minBet, maxBet, noSplash)

	p := tea.NewProgram(model,
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	if _, err := p.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "Error running program:", err)
		os.Exit(1)
	}
}

func validateFlags() error {
	if balance < 1 {
		return fmt.Errorf("--balance must be at least 1")
	}
	if minBet < 1 {
		return fmt.Errorf("--min-bet must be at least 1")
	}
	if minBet > balance {
		return fmt.Errorf("--min-bet cannot be greater than --balance")
	}
	if maxBet > 0 && maxBet < minBet {
		return fmt.Errorf("--max-bet cannot be less than --min-bet")
	}
	return nil
}