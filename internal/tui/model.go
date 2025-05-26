package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/bubbles/list"
)

type screen int

const (
	screenSplash screen = iota
	screenMenu
	screenMonitor
)

const (
	tabHistory = iota
	tabMonitor
	tabFiles
	tabSettings

	numTabs = 4
)

var tabNames = []string{
	"history",
	"monitor",
	"files",
	"settings",
}

type model struct {
	term          string
	cursor        int
	selected      map[int]string
	currentTab    int
	currentScreen screen
	width         int
	height        int
	list          list.Model 
}

func NewModel() model {
	l := list.New([]list.Item{}, itemDelegate{}, 0, 0) 	
	l.Title = "Server History"
	l.Styles.Title = lipgloss.NewStyle().Foreground(lipgloss.Color(coralPink)).Bold(true)
	l.SetItems([]list.Item{
		item{title: "Server 1", description: "Last connected: 2023-10-01"},
		item{title: "Server 2", description: "Last connected: 2023-10-02"},
		item{title: "Server 3", description: "Last connected: 2023-10-03"},
		item{title: "Server 4", description: "Last connected: 2023-10-03"},
		item{title: "Server 5", description: "Last connected: 2023-10-03"},
		item{title: "Server 6", description: "Last connected: 2023-10-03"},
		item{title: "Server 7", description: "Last connected: 2023-10-03"},
		item{title: "Server 8", description: "Last connected: 2023-10-03"},
	})

	return model{
		term:          "main",
		cursor:        0,
		selected:      make(map[int]string),
		currentTab:    0,
		currentScreen: screenSplash,
		width:         0,
		height:        0,
		list:          l,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.updateLayout() 
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "enter":
			switch m.currentScreen {
			case screenSplash:
				m.currentScreen = screenMenu
				m.updateLayout()
				return m, nil
			}
		case "tab", "shift+tab":
			switch m.currentScreen {
			case screenMenu:
				if msg.String() == "shift+tab" {
					m.currentTab = (m.currentTab - 1 + numTabs) % numTabs
				} else {
					m.currentTab = (m.currentTab + 1) % numTabs
				}
			}
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

// View renders the TUI.
func (m model) View() string {
	switch m.currentScreen {
	case screenSplash:
		return renderSplashScreen(m)
	case screenMenu:
		return renderMenuScreen(m)
	default:
		return ""
	}
}

func (m *model) updateLayout() {
	containerWidth := m.getContainerWidth()
	containerHeight := m.getContainerHeight()

	availableHeight := containerHeight - 12
	availableWidth := containerWidth - 6 
	if availableHeight < 5 {
		availableHeight = 5
	}
	if availableWidth < 20 {
		availableWidth = 20
	}

	m.list.SetSize(availableWidth, availableHeight)
}


func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func (m model) getContainerWidth() int {
	const targetMaxWidth = 100
	effectiveWidth := min(targetMaxWidth, m.width-8) // Increased buffer for border visibility
	if effectiveWidth < 60 {
		effectiveWidth = 60
	}
	return effectiveWidth
}

func (m model) getContainerHeight() int {
	const targetMaxHeight = 30
	effectiveHeight := min(targetMaxHeight, m.height-4)
	if effectiveHeight < 15 {
		effectiveHeight = 15
	}
	return effectiveHeight
}

func TeaHandler(s ssh.Session) (tea.Model, []tea.ProgramOption) {
	pty, _, _ := s.Pty()
	m := NewModel()
	m.term = pty.Term
	return m, []tea.ProgramOption{tea.WithAltScreen()}
}
