package tui

import (
	"fmt"
	"io"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
)

type item struct {
	title       string
	description string
}

func (i item) FilterValue() string {
	return i.title
}

func (i item) Title() string {
	return i.title
}

func (i item) Description() string {
	return i.description
}

type itemDelegate struct{}

func (d itemDelegate) Height() int { return 2 }

func (d itemDelegate) Spacing() int { return 0 }

func (d itemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }

func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return 
	}

	title := i.Title()
	description := i.Description()

	titleText := lipgloss.NewStyle().Foreground(lipgloss.Color(colorPrimary))
	descText := lipgloss.NewStyle().Foreground(lipgloss.Color("#888888"))

	content := lipgloss.JoinVertical(lipgloss.Left,
		titleText.Render(title),
		descText.Render(description),
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
