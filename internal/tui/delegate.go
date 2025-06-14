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
	s, ok := listItem.(ServerListItem)
	if !ok {
		return 
	}

	var tag string
	if s.Server.Auth.Method == "key" {
		tag = tagLine.Render("Key")
	} else if s.Server.Auth.Method == "password" {
		tag = tagLinePwd.Render("Password")
	}

	content := lipgloss.JoinVertical(lipgloss.Left,
		titleText.Render(s.Server.Title), 
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

