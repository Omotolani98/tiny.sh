package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/ssh"
)

type screen int

const (
	screenSplash screen = iota
	screenMenu
	screenMonitor
)

const (
    tabSSHKeys = iota
    tabUsers
    tabAPT
    tabMonitor

    numTabs = 4
)

var (
    colorBackground = "#1a1a1a"
    colorPrimary    = "#eeeeee"
    colorAccent     = "#7D56F4"
    colorSuccess    = "#2cb67d"
    colorWarning    = "#fbbf24"
    colorError      = "#ef4444"
    colorBorder     = "#333333"
		coralPink				= "#FB9F89"
)

var (
    activeTabStyle = lipgloss.NewStyle().
				Border(lipgloss.NormalBorder()).
				PaddingLeft(5).
				PaddingRight(5).
				BorderForeground(lipgloss.Color(coralPink)).
				BorderStyle(lipgloss.RoundedBorder()).
        Bold(true)

    inactiveTabStyle = lipgloss.NewStyle().
        Foreground(lipgloss.Color(colorPrimary)).
				PaddingLeft(5).
				PaddingRight(5).
				Border(lipgloss.NormalBorder()).
        Bold(false)
)

var tabNames = []string{
    "history",
    "monitor",
    "files",
    "settings",
}

type model struct {
	term string
	cursor         int
	selected       map[int]string
	currentTab     int
	currentScreen  screen
	width          int
	height         int
}

func NewModel() model {
    return model{
				term: "main",
        cursor: 0,
        selected: make(map[int]string),
        currentTab: 0,
        currentScreen: screenSplash,
        width: 0,
        height: 0,
    }
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) View() string {
		switch m.currentScreen {
    case screenSplash:
        return m.renderSplashScreen()
    case screenMenu:
        return m.renderTabs()
    default:
        return ""
    }
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tea.KeyMsg:
		switch msg.String() {
			case "ctrl+c", "q":
				return m, tea.Quit
			case "enter":
				switch m.currentScreen {
					case screenSplash:
						m.currentScreen = screenMenu
						return m, nil
				}
			case "tab", "shift+tab":
   			switch m.currentScreen {
    		case screenMenu:
        // Tab Navigation between Tabs
        if msg.String() == "shift+tab" {
            m.currentTab = (m.currentTab - 1 + numTabs) % numTabs
        } else {
            m.currentTab = (m.currentTab + 1) % numTabs
        }
    	}
		}
	}

	return m, nil
}

func (m model) renderSplashScreen() string {
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

func (m model) renderTabs() string {
    var tabs []string

    for i, name := range tabNames {
        if i == m.currentTab {
					  tabs = append(tabs, activeTabStyle.Render(fmt.Sprintf(" %s ", name)))
        } else {
            tabs = append(tabs, inactiveTabStyle.Render(fmt.Sprintf(" %s ", name)))
        }
    }

    return lipgloss.Place(
			m.width,
			m.height,
			lipgloss.Center,
			lipgloss.Center,
			lipgloss.JoinHorizontal(lipgloss.Center, tabs...),
		)
}

func TeaHandler(s ssh.Session) (tea.Model, []tea.ProgramOption) {
	// This should never fail, as we are using the activeterm middleware.
	pty, _, _ := s.Pty()

	// When running a Bubble Tea app over SSH, you shouldn't use the default
	// lipgloss.NewStyle function.
	// That function will use the color profile from the os.Stdin, which is the
	// server, not the client.
	// We provide a MakeRenderer function in the bubbletea middleware package,
	// so you can easily get the correct renderer for the current session, and
	// use it to create the styles.
	// The recommended way to use these styles is to then pass them down to
	// your Bubble Tea model.
	//renderer := t
	//txtStyle := renderer.NewStyle().Foreground(lipgloss.Color("10"))
	//quitStyle := renderer.NewStyle().Foreground(lipgloss.Color("8"))


	m := NewModel()
	m.term = pty.Term
	return m, []tea.ProgramOption{tea.WithAltScreen()}
}
