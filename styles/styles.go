package styles

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	ColorBackground lipgloss.Color
	ColorSurface    lipgloss.Color
	ColorBorder    lipgloss.Color
	ColorText      lipgloss.Color
	ColorSubtle    lipgloss.Color
	ColorAccent    lipgloss.Color
	ColorSuccess   lipgloss.Color
	ColorDanger    lipgloss.Color
	ColorWarning   lipgloss.Color
	ColorRed       lipgloss.Color
	ColorNeutral   lipgloss.Color
)

func init() {
	InitializeColors()
}

func InitializeColors() {
	if lipgloss.HasDarkBackground() {
		ColorBackground = lipgloss.Color("#1A1A2E")
		ColorSurface = lipgloss.Color("#16213E")
		ColorBorder = lipgloss.Color("#3A3A5C")
		ColorText = lipgloss.Color("#E0E0E0")
		ColorSubtle = lipgloss.Color("#6C6C8A")
		ColorAccent = lipgloss.Color("#60A5FA")
		ColorSuccess = lipgloss.Color("#4ADE80")
		ColorDanger = lipgloss.Color("#F87171")
		ColorWarning = lipgloss.Color("#FBBF24")
		ColorRed = lipgloss.Color("#F87171")
		ColorNeutral = lipgloss.Color("#94A3B8")
	} else {
		ColorBackground = lipgloss.Color("#F5F5F5")
		ColorSurface = lipgloss.Color("#FFFFFF")
		ColorBorder = lipgloss.Color("#CCCCCC")
		ColorText = lipgloss.Color("#1A1A1A")
		ColorSubtle = lipgloss.Color("#888888")
		ColorAccent = lipgloss.Color("#2563EB")
		ColorSuccess = lipgloss.Color("#16A34A")
		ColorDanger = lipgloss.Color("#DC2626")
		ColorWarning = lipgloss.Color("#D97706")
		ColorRed = lipgloss.Color("#DC2626")
		ColorNeutral = lipgloss.Color("#64748B")
	}
}

var (
	StyleBackground = lipgloss.NewStyle().
			Background(ColorBackground)

	StyleSurface = lipgloss.NewStyle().
			Background(ColorSurface).
			Foreground(ColorText)

	StyleBorder = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorBorder)

	StyleTitle = lipgloss.NewStyle().
			Foreground(ColorText).
			Bold(true)

	StyleSubtle = lipgloss.NewStyle().
			Foreground(ColorSubtle)

	StyleAccent = lipgloss.NewStyle().
			Foreground(ColorAccent).
			Bold(true)

	StyleSuccess = lipgloss.NewStyle().
			Foreground(ColorSuccess).
			Bold(true)

	StyleDanger = lipgloss.NewStyle().
			Foreground(ColorDanger).
			Bold(true)

	StyleWarning = lipgloss.NewStyle().
			Foreground(ColorWarning)

	HUDStyle = lipgloss.NewStyle().
			Background(ColorSurface).
			Foreground(ColorText)

	CardStyle = lipgloss.NewStyle().
			Background(ColorSurface).
			Border(lipgloss.NormalBorder()).
			BorderForeground(ColorBorder)

	CardHighlightedStyle = lipgloss.NewStyle().
				Background(ColorSurface).
				Border(lipgloss.NormalBorder()).
				BorderForeground(ColorAccent)

	CardDimmedStyle = lipgloss.NewStyle().
			Background(ColorSurface).
			Foreground(ColorSubtle).
			Border(lipgloss.NormalBorder()).
			BorderForeground(ColorSubtle)

	ButtonStyle = lipgloss.NewStyle().
			Padding(0, 2).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorBorder).
			Foreground(ColorText)

	ButtonFocusedStyle = lipgloss.NewStyle().
				Padding(0, 2).
				Border(lipgloss.RoundedBorder()).
				BorderForeground(ColorAccent).
				Foreground(ColorAccent).
				Bold(true)

	ButtonDisabledStyle = lipgloss.NewStyle().
				Padding(0, 2).
				Border(lipgloss.RoundedBorder()).
				BorderForeground(ColorBorder).
				Foreground(ColorSubtle)

	ModalStyle = lipgloss.NewStyle().
			Background(ColorSurface).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorAccent).
			Padding(1, 2)

	RedCardStyle = lipgloss.NewStyle().
			Background(ColorSurface).
			Foreground(ColorRed).
			Border(lipgloss.NormalBorder()).
			BorderForeground(ColorBorder)

	BlackCardStyle = lipgloss.NewStyle().
			Background(ColorSurface).
			Foreground(ColorText).
			Border(lipgloss.NormalBorder()).
			BorderForeground(ColorBorder)

	RedCardHighlightedStyle = lipgloss.NewStyle().
			Background(ColorSurface).
			Foreground(ColorRed).
			Border(lipgloss.NormalBorder()).
			BorderForeground(ColorAccent)

	BlackCardHighlightedStyle = lipgloss.NewStyle().
			Background(ColorSurface).
			Foreground(ColorText).
			Border(lipgloss.NormalBorder()).
			BorderForeground(ColorAccent)
)

const (
	HUDHeight       = 1
	ActionBarHeight = 3
	CardWidth       = 5
	CardHeight      = 5
	CardGap         = 1
	MinTermWidth    = 80
	MinTermHeight   = 24
)