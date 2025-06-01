package tui

import "github.com/charmbracelet/lipgloss"

const (
	colorBackground = "#1a1a1a"
	colorPrimary    = "#eeeeee"
	colorAccent     = "#7D56F4"
	colorSuccess    = "#2cb67d"
	colorWarning    = "#fbbf24"
	colorError      = "#ef4444"
	colorBorder     = "#333333"
	coralPink       = "#FB9F89"
)

var (
	// Tabs
	activeTabStyle = lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color(coralPink)).
		BorderStyle(lipgloss.RoundedBorder()).
		Bold(true)

	inactiveTabStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(colorPrimary)).
		Border(lipgloss.NormalBorder()).
		Bold(false)

	// List Items (for history_list.go)
	itemStyle = lipgloss.NewStyle().PaddingLeft(2).PaddingRight(2)

	selectedItemStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(coralPink)).
		PaddingLeft(1).
		PaddingRight(1).
		MarginBottom(0)

	// Generic Content Area Border
	contentAreaBorder = lipgloss.NewStyle().
//		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(colorBorder)).
		Padding(1, 2) // top/bottom, left/right padding

	// Footer
	footerStyle = lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderTop(true).
		BorderForeground(lipgloss.Color(colorBorder)).
		Padding(1, 2).
		Foreground(lipgloss.Color("#888888")) // Subtle color for help text

	formTitleStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(coralPink)).
		Bold(true).
		PaddingBottom(1)

	focusedInputStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(colorPrimary)).
		Border(lipgloss.NormalBorder()).
		BorderBottom(true).
		BorderBottomForeground(lipgloss.Color(coralPink))

	blurredInputStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(colorPrimary)).
		Border(lipgloss.NormalBorder()).
		BorderBottom(true).
		BorderBottomForeground(lipgloss.Color("#555555"))

	focusedPromptStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(coralPink))

	blurredPromptStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#888888"))

	cursorStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(coralPink))

	// Style for the focus indicator (↳)
	focusIndicatorStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(coralPink))
)
