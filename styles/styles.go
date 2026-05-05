package styles

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	ColorBackground lipgloss.Color
	ColorSurface    lipgloss.Color
	ColorBorder     lipgloss.Color
	ColorText       lipgloss.Color
	ColorSubtle     lipgloss.Color
	ColorAccent     lipgloss.Color
	ColorSuccess    lipgloss.Color
	ColorDanger     lipgloss.Color
	ColorWarning    lipgloss.Color
	ColorRed        lipgloss.Color
	ColorNeutral    lipgloss.Color
)

var initialized = false

func EnsureInit() {
	if initialized {
		return
	}
	InitializeColors()
	initialized = true
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
	initialized = true
}

func GetStyleBackground() lipgloss.Style {
	EnsureInit()
	return lipgloss.NewStyle().
		Background(ColorBackground)
}

func GetHUDStyle() lipgloss.Style {
	EnsureInit()
	return lipgloss.NewStyle().
		Background(ColorSurface).
		Foreground(ColorText)
}

func GetCardStyle(red, highlighted, dimmed bool) lipgloss.Style {
	EnsureInit()
	var style lipgloss.Style
	if red {
		if highlighted {
			style = lipgloss.NewStyle().
				Background(ColorSurface).
				Foreground(ColorRed).
				Border(lipgloss.NormalBorder()).
				BorderForeground(ColorAccent)
		} else {
			style = lipgloss.NewStyle().
				Background(ColorSurface).
				Foreground(ColorRed).
				Border(lipgloss.NormalBorder()).
				BorderForeground(ColorBorder)
		}
	} else {
		if highlighted {
			style = lipgloss.NewStyle().
				Background(ColorSurface).
				Foreground(ColorText).
				Border(lipgloss.NormalBorder()).
				BorderForeground(ColorAccent)
		} else {
			style = lipgloss.NewStyle().
				Background(ColorSurface).
				Foreground(ColorText).
				Border(lipgloss.NormalBorder()).
				BorderForeground(ColorBorder)
		}
	}
	if dimmed {
		style = lipgloss.NewStyle().
			Background(ColorSurface).
			Foreground(ColorSubtle).
			Border(lipgloss.NormalBorder()).
			BorderForeground(ColorSubtle)
	}
	return style
}

func GetButtonStyle(disabled, focused bool) lipgloss.Style {
	EnsureInit()
	var style lipgloss.Style
	if disabled {
		style = lipgloss.NewStyle().
			Padding(0, 2).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorBorder).
			Foreground(ColorSubtle)
	} else if focused {
		style = lipgloss.NewStyle().
			Padding(0, 2).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorAccent).
			Foreground(ColorAccent).
			Bold(true)
	} else {
		style = lipgloss.NewStyle().
			Padding(0, 2).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorBorder).
			Foreground(ColorText)
	}
	return style
}

func GetModalStyle() lipgloss.Style {
	EnsureInit()
	return lipgloss.NewStyle().
		Background(ColorSurface).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ColorAccent).
		Padding(1, 2)
}

func GetSuccessStyle() lipgloss.Style {
	EnsureInit()
	return lipgloss.NewStyle().
		Foreground(ColorSuccess).
		Bold(true)
}

func GetDangerStyle() lipgloss.Style {
	EnsureInit()
	return lipgloss.NewStyle().
		Foreground(ColorDanger).
		Bold(true)
}

const (
	HUDHeight      = 1
	ActionBarHeight = 3
	CardWidth      = 5
	CardHeight     = 5
	CardGap        = 1
	MinTermWidth   = 80
	MinTermHeight  = 24
)