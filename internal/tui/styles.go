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
	colorPending	= "#FDCA40"
	colorSubtle = "#888888"
	colorGrey = "#555555"
)

var (
	activeTabStyle = lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color(coralPink)).
		BorderStyle(lipgloss.RoundedBorder()).
		Bold(true)

	inactiveTabStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(colorPrimary)).
		Border(lipgloss.NormalBorder()).
		Bold(false)

	itemStyle = lipgloss.NewStyle().PaddingLeft(2).PaddingRight(2)

	selectedItemStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(coralPink)).
		PaddingLeft(1).
		PaddingRight(1).
		MarginBottom(0)

	contentAreaBorder = lipgloss.NewStyle().
		BorderForeground(lipgloss.Color(colorBorder)).
		Padding(1, 2)

	footerStyle = lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderTop(true).
		BorderForeground(lipgloss.Color(colorBorder)).
		Padding(1, 2).
		Foreground(lipgloss.Color(colorSubtle))

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
		BorderBottomForeground(lipgloss.Color(colorGrey))

	focusedPromptStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(coralPink))

	blurredPromptStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(colorSubtle))

	cursorStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(coralPink))

	focusIndicatorStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(coralPink))


	blockStyle = lipgloss.NewStyle().
		Padding(1, 2).
		MarginBottom(1).
		BorderLeft(true).
		BorderLeftForeground(lipgloss.Color(coralPink))

	titleText = lipgloss.NewStyle().Foreground(lipgloss.Color(colorPrimary)).Bold(true)
	detailText = lipgloss.NewStyle().Foreground(lipgloss.Color("#888888")).PaddingRight(1)
	labelStyle = lipgloss.NewStyle().Width(15).Align(lipgloss.Left).Foreground(lipgloss.Color(colorPrimary))
	
  tagLine = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("#E1F0C4")).Foreground(lipgloss.Color("#E1F0C4")).PaddingRight(1).PaddingLeft(1)
  tagLinePwd = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("#B8B8FF")).Foreground(lipgloss.Color("#B8B8FF")).PaddingRight(1).PaddingLeft(1)
)
