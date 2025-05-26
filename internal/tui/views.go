package tui

import "github.com/charmbracelet/lipgloss"

// renderSplashScreen renders the initial splash screen.
func renderSplashScreen(m model) string {
	splash := `
████████╗██╗███╗   ██╗██╗   ██╗███████╗██╗  ██╗
╚══██╔══╝██║████╗  ██║╚██╗ ██╔╝██╔════╝██║  ██║
   ██║   ██║██╔██╗ ██║ ╚████╔╝ ███████╗███████║
   ██║   ██║██║╚██╗██║  ╚██╔╝  ╚════██║██╔══██║
   ██║   ██║██║ ╚████║   ██║██╗███████║██║  ██║
   ╚═╝   ╚═╝╚═╝  ╚═══╝   ╚═╝╚═╝╚══════╝╚═╝  ╚═╝
   Press Enter to continue
		`

	return lipgloss.NewStyle().
		Width(m.width).
		Height(m.height).
		Align(lipgloss.Center, lipgloss.Center).
		Foreground(lipgloss.Color(colorPrimary)).
		Render(splash)
}

// renderMenuScreen renders the main menu interface.
func renderMenuScreen(m model) string {
	// Render individual components
	tabsView := renderTabs(m)
	currentTabView := renderCurrentTab(m)
	footerView := renderFooter(m)

	// Join the components vertically.
	mainContent := lipgloss.JoinVertical(
		lipgloss.Left,
		tabsView,
		"", // Empty line for spacing
		currentTabView,
		"", // Empty line for spacing
		footerView,
	)

	containerHeight := m.getContainerHeight()

	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		lipgloss.NewStyle().
			Width(m.getContainerWidth()).
			Height(containerHeight). 
			Render(mainContent),
	)
}

func renderTabs(m model) string {
	var tabs []string
	containerWidth := m.getContainerWidth()
	tabWidth := (containerWidth - (len(tabNames) * 2)) / len(tabNames) 

	for i, name := range tabNames {
		tabContent := lipgloss.NewStyle().
			Width(tabWidth).
			Align(lipgloss.Center).
			PaddingLeft(5). 
			PaddingRight(5). 
			Render(name)

		if i == m.currentTab {
			tabs = append(tabs, activeTabStyle.Width(tabWidth).Render(tabContent))
		} else {
			tabs = append(tabs, inactiveTabStyle.Width(tabWidth).Render(tabContent))
		}
	}

	tabBar := lipgloss.NewStyle().
		Width(containerWidth).
		Render(lipgloss.JoinHorizontal(lipgloss.Left, tabs...))
	return tabBar
}

func renderCurrentTab(m model) string {
	containerWidth := m.getContainerWidth()
	var content string
	switch m.currentTab {
	case tabHistory:
		content = renderHistoryTabContent(m)
	case tabMonitor:
		content = "Monitor Tab Content"
	case tabFiles:
		content = "Files Tab Content"
	case tabSettings:
		content = "Settings Tab Content"
	default:
		content = "Unknown Tab"
	}

	contentArea := contentAreaBorder.
		Width(containerWidth).
		Render(content)

	return contentArea
}

func renderHistoryTabContent(m model) string {
	return m.list.View()
}

func renderFooter(m model) string {
	containerWidth := m.getContainerWidth()
	helpText := "Tab/Shift+Tab: Switch tabs • Enter: Select • Ctrl+C/q: Quit"

	footer := footerStyle.
		Width(containerWidth).
		Render(helpText)

	return footer
}
