package tui

import (
	"fmt"
	"io"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/bubbles/list"
)

type fileItemDelegate struct{}

func (d fileItemDelegate) Height() int          { return 1 }
func (d fileItemDelegate) Spacing() int         { return 0 }
func (d fileItemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd {
    return nil
}

func (d fileItemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
    item, ok := listItem.(FileListItem)
    if !ok {
        return
    }

		content := lipgloss.JoinVertical(lipgloss.Left,
			titleText.Render(item.Name), 
			detailText.Render(item.Info),
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
