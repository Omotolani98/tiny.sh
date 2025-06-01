package tui

import (
	"fmt"
	"io"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/bubbles/list"
)

type itemDelegate struct{}

func (d itemDelegate) Height() int { return 4 }

func (d itemDelegate) Spacing() int { return 0 }

func (d itemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }


func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	// Assert the listItem to our ServerHistoryItem type
	s, ok := listItem.(ServerListItem)
	if !ok {
		return 
	}

	// Styles for text and labels
	titleText := lipgloss.NewStyle().Foreground(lipgloss.Color(colorPrimary)).Bold(true)
	detailText := lipgloss.NewStyle().Foreground(lipgloss.Color("#888888")).PaddingRight(1) // Subtle grey for details
	labelStyle := lipgloss.NewStyle().Width(15).Align(lipgloss.Left).Foreground(lipgloss.Color(colorPrimary))

	
  tagLine := lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("#E1F0C4")).Foreground(lipgloss.Color("#E1F0C4")).PaddingRight(1).PaddingLeft(1)
  tagLinePwd := lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("#B8B8FF")).Foreground(lipgloss.Color("#B8B8FF")).PaddingRight(1).PaddingLeft(1)

	var tag string
	if s.Server.Auth.Method == "key" {
		tag = tagLine.Render("Key")
	} else if s.Server.Auth.Method == "password" {
		tag = tagLinePwd.Render("Password")
	}
	// Join the main title and the three formatted detail lines vertically
	content := lipgloss.JoinVertical(lipgloss.Left,
		titleText.Render(s.Server.Title), // Corrected: Use s.Title() as ServerHistoryItem now implements list.Item
		labelStyle.Render(s.Server.Username),
		detailText.Render(s.Server.LastConnected),
		lipgloss.JoinHorizontal(lipgloss.Left, tag),
	)

	var renderedItem string
	if index == m.Index() {
			renderedItem = selectedItemStyle.
			Width(m.Width() - selectedItemStyle.GetHorizontalFrameSize()).
			Render(content)
	} else {
			renderedItem = itemStyle.
			Width(m.Width() - itemStyle.GetHorizontalFrameSize()).
			Render(content)
	}

	fmt.Fprint(w, renderedItem)
}

