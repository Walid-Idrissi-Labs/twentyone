# twentyone

Vegas Strip Blackjack — A terminal-based Blackjack game written in Go.

![Terminal Blackjack](https://img.shields.io/badge/TUI-Blackjack-2E86AB?style=flat-square)
![Go Version](https://img.shields.io/badge/Go-1.22+-00ADD8?style=flat-square)
![License](https://img.shields.io/badge/License-MIT-green?style=flat-square)

A full-screen terminal Blackjack game following Vegas Strip rules, with support for splitting, doubling down, insurance, and more.

## Features

- Full-screen TUI that adapts to any terminal size (minimum 80x24)
- Complete Vegas Strip Blackjack rules faithfully implemented
- 6-deck shoe with automatic reshuffling
- Split up to 4 hands
- Double down on any two cards
- Insurance when dealer shows an Ace
- Blackjack pays 3:2
- Keyboard shortcuts and mouse click support
- Modern minimal visual style using Unicode box-drawing characters

## Installation

### Prerequisites

- Go 1.22 or later
- A terminal emulator with Unicode support
- Terminal must support 256 colors

### Build from Source

```bash
git clone https://github.com/twentyone/twentyone.git
cd twentyone
go mod download
go build -o twentyone .
```

### Run

```bash
./twentyone --balance 500
```

## Usage

```
twentyone [flags]

Flags:
  --balance int    Starting balance in dollars (default 1000)
  --min-bet int    Minimum bet (default 1)
  --max-bet int    Maximum bet, 0 = no limit (default 0)
  --no-splash      Skip the welcome screen
```

## Controls

### Bet Screen
| Key | Action |
|-----|--------|
| 0-9 | Type bet amount |
| Backspace | Delete last digit |
| Arrow Up | +$5 |
| Arrow Down | -$5 |
| Arrow Right | +$25 |
| Arrow Left | -$25 |
| Enter | Confirm bet |
| Q / Esc | Quit |

### Game Table
| Key | Action |
|-----|--------|
| H | Hit |
| S | Stand |
| D | Double |
| P | Split |
| Y | Take Insurance |
| N | Decline Insurance |
| Q / Esc | Quit to summary |

## Game Rules

twentyone implements Vegas Strip Blackjack rules:

- 6-deck shoe (312 cards)
- Dealer stands on ALL 17s (including soft 17)
- Blackjack pays 3:2
- Insurance offered when dealer shows Ace
- Split up to 3 times (4 hands maximum)
- Split Aces receive only one card
- Double down on any two cards
- Double after split allowed

## License

MIT License - see LICENSE file for details.

## Acknowledgments

Built with [Charmbracelet](https://charm.sh/) libraries:
- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - TUI framework
- [Lip Gloss](https://github.com/charmbracelet/lipgloss) - Styling
- [Bubbles](https://github.com/charmbracelet/bubbles) - TUI components