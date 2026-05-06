package styles

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	ColorBackground lipgloss.Color = lipgloss.Color("#1A1A2E")
	ColorSurface    lipgloss.Color = lipgloss.Color("#16213E")
	ColorBorder     lipgloss.Color = lipgloss.Color("#3A3A5C")
	ColorText       lipgloss.Color = lipgloss.Color("#E0E0E0")
	ColorSubtle     lipgloss.Color = lipgloss.Color("#6C6C8A")
	ColorAccent     lipgloss.Color = lipgloss.Color("#60A5FA")
	ColorSuccess    lipgloss.Color = lipgloss.Color("#4ADE80")
	ColorDanger     lipgloss.Color = lipgloss.Color("#F87171")
	ColorWarning    lipgloss.Color = lipgloss.Color("#FBBF24")
	ColorRed        lipgloss.Color = lipgloss.Color("#F87171")
	ColorNeutral    lipgloss.Color = lipgloss.Color("#94A3B8")
)

func EnsureInit() {}

var (
	StyleBackground = lipgloss.NewStyle().
				Background(ColorBackground)

	StyleHUD = lipgloss.NewStyle().
			Background(ColorSurface).
			Foreground(ColorText).
			Width(80)

	StyleCard = lipgloss.NewStyle().
			Background(ColorSurface).
			Foreground(ColorText).
			Border(lipgloss.NormalBorder()).
			BorderForeground(ColorBorder)

	StyleCardRed = lipgloss.NewStyle().
			Background(ColorSurface).
			Foreground(ColorRed).
			Border(lipgloss.NormalBorder()).
			BorderForeground(ColorBorder)

	StyleCardHighlighted = lipgloss.NewStyle().
				Background(ColorSurface).
				Foreground(ColorText).
				Border(lipgloss.NormalBorder()).
				BorderForeground(ColorAccent)

	StyleCardRedHighlighted = lipgloss.NewStyle().
				Background(ColorSurface).
				Foreground(ColorRed).
				Border(lipgloss.NormalBorder()).
				BorderForeground(ColorAccent)

	StyleCardDimmed = lipgloss.NewStyle().
			Background(ColorSurface).
			Foreground(ColorSubtle).
			Border(lipgloss.NormalBorder()).
			BorderForeground(ColorSubtle)

	StyleButton = lipgloss.NewStyle().
			Padding(0, 2).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorBorder).
			Foreground(ColorText)

	StyleButtonFocused = lipgloss.NewStyle().
				Padding(0, 2).
				Border(lipgloss.RoundedBorder()).
				BorderForeground(ColorAccent).
				Foreground(ColorAccent).
				Bold(true)

	StyleButtonDisabled = lipgloss.NewStyle().
				Padding(0, 2).
				Border(lipgloss.RoundedBorder()).
				BorderForeground(ColorSubtle).
				Foreground(ColorSubtle)

	StyleModal = lipgloss.NewStyle().
			Background(ColorSurface).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorAccent).
			Padding(1, 3)

	StyleTitle = lipgloss.NewStyle().
			Foreground(ColorText).
			Bold(true)

	StyleDanger = lipgloss.NewStyle().
			Foreground(ColorDanger).
			Bold(true)

	StyleSuccess = lipgloss.NewStyle().
			Foreground(ColorSuccess).
			Bold(true)

	StyleDim = lipgloss.NewStyle().
			Foreground(ColorSubtle)
)

const (
	MinTermWidth  = 80
	MinTermHeight = 24
)