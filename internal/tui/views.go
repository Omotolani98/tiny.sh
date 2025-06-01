package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

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

func renderMenuScreen(m model) string {
	tabsView := renderTabs(m)
	currentTabView := renderCurrentTab(m)
	footerView := renderFooter(m)

	mainContent := lipgloss.JoinVertical(
		lipgloss.Left,
		tabsView,
		"", // Empty line for spacing
		currentTabView,
		"",
		renderConnectionStatus(m),
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

func renderConnectionForm(m model) string {
	var b strings.Builder

	b.WriteString(formTitleStyle.Render("🔌 Connect to a new server") + "\n\n")

	for i, input := range m.form.inputs {
		b.WriteString(input.View() + "\n")
		if i == m.form.focus {
			b.WriteString(focusIndicatorStyle.Render("↳") + "\n")
		}
	}

	helpText := "[Tab] to switch • [Enter] to submit • [Esc] to cancel"
	b.WriteString("\n" + lipgloss.NewStyle().Foreground(lipgloss.Color("#888888")).Render(helpText))

	formContent := contentAreaBorder.
		Width(m.getContainerWidth()).
		Render(b.String())

	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		formContent,
	)
}

func renderConnectionStatus(m model) string {
	selectedItem, ok := m.list.SelectedItem().(ServerListItem)
	if !ok {
		return ""
	}

	server := selectedItem.Server

	if m.connecting && m.connectingToHost == server.Host {
		return lipgloss.NewStyle().
			Foreground(lipgloss.Color(colorAccent)).
			Render(fmt.Sprintf("⏳ Connecting to %s...", server.Host))
	}

	if m.connectionErr != nil && m.currentServer != nil && m.currentServer.Host == server.Host {
		return lipgloss.NewStyle().
			Foreground(lipgloss.Color(colorError)).
			Render("🔌 Connection failed: " + m.connectionErr.Error())
	}

	if m.isConnected && m.currentServer != nil && m.currentServer.Host == server.Host {
		return lipgloss.NewStyle().
			Foreground(lipgloss.Color(colorSuccess)).
			Render(fmt.Sprintf("🟢 Connected to %s@%s", server.Username, server.Host))
	}

	return ""
}
