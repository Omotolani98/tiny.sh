package tui

import (
	"fmt"
	"time"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/log"
	"github.com/charmbracelet/bubbles/textinput"
	sshClient "golang.org/x/crypto/ssh"
)

type screen int

const (
	screenSplash screen = iota
	screenMenu
	screenMonitor
	screenForm
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
	form					ConnectionForm
	serverHistory []ServerHistoryItem
	activeSSHClient *sshClient.Client
	activeSession	*sshClient.Session
	connectionErr error
	isConnected		bool
	currentServer *ServerHistoryItem
	connecting 		bool
	connectingToHost string
	connectRequestID string
	monitorOutput string
	previousTab int
	cpuMetrics    string
	memoryMetrics string
	diskMetrics   string
	fileList      list.Model
	fileListItems []FileListItem
}

func NewModel() model {
	formInputs := make([]textinput.Model, 6)
	placeholders := []string{
		"Host (e.g., example.com)", 
		"Port (e.g., 22)", 
		"Username",
		"Password (optional, for key-based auth)",
		"Authentication Method (password, key, etc.)",
		"Auth Path (optional, for key-based auth, e.g. ~/.ssh/id_rsa)",
	}
	passwordIdx := 3
	authMethodIdx := 4
	authPathIdx := 5
	
	for i := range formInputs {
		ti := textinput.New()
		ti.Placeholder = placeholders[i]
		ti.CharLimit = 100
		ti.Width = 50

		ti.PromptStyle = blurredPromptStyle
		ti.TextStyle = blurredInputStyle
		ti.CursorStyle = cursorStyle
		ti.PlaceholderStyle = blurredPromptStyle

		if i == passwordIdx {
			ti.EchoMode = textinput.EchoPassword
			ti.EchoCharacter = '•'
		} else if i == authMethodIdx {
		} else if i == authPathIdx {
			ti.EchoMode = textinput.EchoNormal
		}
		formInputs[i] = ti
		formInputs[i] = ti
	}

	formInputs[0].Focus()
	formInputs[0].PromptStyle = focusedPromptStyle
	formInputs[0].TextStyle = focusedInputStyle
	formInputs[0].PlaceholderStyle = focusedPromptStyle

	form := ConnectionForm{
		inputs: formInputs,
		focus:  0,
		passwordInputIndex:   passwordIdx,
		authMethodInputIndex: authMethodIdx,
		authPathInputIndex:   authPathIdx,
	}

	l := list.New([]list.Item{}, itemDelegate{}, 0, 0) 	
	l.Title = "Server History"
	l.Styles.Title = lipgloss.NewStyle().Foreground(lipgloss.Color(coralPink)).Bold(true)

	history, err := LoadHistory()
	if err != nil {
    log.Info("failed to load history", err)
    history = []ServerHistoryItem{}
	}

	l.SetItems(ToListItems(history))
	
	fl := list.New([]list.Item{}, fileItemDelegate{}, 0, 0)
	fl.Title = "Files"
	return model{
		term:          "main",
		cursor:        0,
		selected:      make(map[int]string),
		currentTab:    0,
		currentScreen: screenSplash,
		width:         0,
		height:        0,
		list:          l,
		form:          form,
		serverHistory: history,
		fileList: fl,
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
		case "n":
			if m.currentScreen != screenForm {
				m.currentScreen = screenForm
				for i := range m.form.inputs {
					m.form.inputs[i].SetValue("")
					m.form.inputs[i].Blur()
					if i == m.form.passwordInputIndex {
						m.form.inputs[i].EchoMode = textinput.EchoPassword
					}
					if i == m.form.authPathInputIndex {
						m.form.inputs[i].EchoMode = textinput.EchoNormal
					}
					m.form.inputs[i].PromptStyle = blurredPromptStyle
					m.form.inputs[i].TextStyle = blurredInputStyle
					m.form.inputs[i].PlaceholderStyle = blurredPromptStyle
				}
				m.form.focus = 0
				m.form.inputs[0].Focus()
				m.form.inputs[0].PromptStyle = focusedPromptStyle
				m.form.inputs[0].TextStyle = focusedInputStyle
				m.form.inputs[0].PlaceholderStyle = focusedPromptStyle
				return m, textinput.Blink
			}
		case "esc":
			if m.currentScreen == screenForm {
				m.currentScreen = screenMenu
				m.updateLayout()
				return m, nil
			}
		case "ctrl+c", "q":
			return m, tea.Quit
		case "enter":
			switch m.currentScreen {
			case screenSplash:
				m.currentScreen = screenMenu
				m.updateLayout()
				return m, nil
			case screenForm:
				portVal := 22
				if p, err := fmt.Sscanf(m.form.inputs[1].Value(), "%d", &portVal); err != nil || p != 1 {
					log.Printf("Invalid port input: %s, defaulting to 22", m.form.inputs[1].Value())
					portVal = 22
				}

				authMethod := m.form.inputs[m.form.authMethodInputIndex].Value()
				authPath := ""
				authValue := ""

				if authMethod == "key" {
					authPath = m.form.inputs[m.form.authPathInputIndex].Value()
				} else if authMethod == "password" {
					authValue = m.form.inputs[m.form.passwordInputIndex].Value()
				}

				newServer := ServerHistoryItem{
					Title:         m.form.inputs[0].Value(), // Use host as title for now
					Description:   fmt.Sprintf("Last connected: %s", time.Now().Format("2006-01-02 15:04")),
					Host:          m.form.inputs[0].Value(),
					Port:          portVal,
					Username:      m.form.inputs[2].Value(),
					Auth: AuthDetails{Method: authMethod, Path: authPath, Value: authValue},
					LastConnected: time.Now().Format("2006-01-02 15:04"),
				}

				m.serverHistory = append(m.serverHistory, newServer)

				if err := SaveHistory(m.serverHistory); err != nil {
					log.Printf("Error saving history: %v", err)
					//TODO: display an error message in the TUI here
				}

				m.list.SetItems(ToListItems(m.serverHistory))
				m.list.Select(len(m.list.Items()) - 1)

				m.currentScreen = screenMenu
				m.updateLayout()
				return m, nil
			case screenMenu:	
				if m.currentTab == tabHistory {
					selected, ok := m.list.SelectedItem().(ServerListItem)
					if ok {
						requestID := fmt.Sprintf("%s-%d", selected.Server.Host, time.Now().UnixNano())
						m.connecting = true
						m.connectingToHost = selected.Server.Host
						m.connectRequestID = requestID
						return m, func() tea.Msg {
							return trySSHConnect(selected.Server, requestID)
						}
					}
				}
			}
		case "tab", "shift+tab":
			if m.currentScreen == screenForm {
				m.form.inputs[m.form.focus].Blur()
				m.form.inputs[m.form.focus].PromptStyle = blurredPromptStyle
				m.form.inputs[m.form.focus].TextStyle = blurredInputStyle
				m.form.inputs[m.form.focus].PlaceholderStyle = blurredPromptStyle

				nextFocus := m.form.focus
				if msg.String() == "tab" {
					nextFocus = (m.form.focus + 1) % len(m.form.inputs)
				} else {
					nextFocus = (m.form.focus - 1 + len(m.form.inputs)) % len(m.form.inputs)
				}

				currentAuthMethod := m.form.inputs[m.form.authMethodInputIndex].Value()
				for {
					if nextFocus == m.form.passwordInputIndex {
						if currentAuthMethod != "password" {
							if msg.String() == "tab" {
								nextFocus = (nextFocus + 1) % len(m.form.inputs)
							} else {
								nextFocus = (nextFocus - 1 + len(m.form.inputs)) % len(m.form.inputs)
							}
							continue
						}
					} else if nextFocus == m.form.authPathInputIndex {
						if currentAuthMethod != "key" {
							if msg.String() == "tab" {
								nextFocus = (nextFocus + 1) % len(m.form.inputs)
							} else {
								nextFocus = (nextFocus - 1 + len(m.form.inputs)) % len(m.form.inputs)
							}
							continue
						}
					}
					break
				}
				m.form.focus = nextFocus

				m.form.inputs[m.form.focus].Focus()
				m.form.inputs[m.form.focus].PromptStyle = focusedPromptStyle
				m.form.inputs[m.form.focus].TextStyle = focusedInputStyle
				m.form.inputs[m.form.focus].PlaceholderStyle = focusedPromptStyle
				return m, textinput.Blink
			} else if m.currentScreen == screenMenu {
				if msg.String() == "shift+tab" {
					m.currentTab = (m.currentTab - 1 + numTabs) % numTabs
				} else {
					m.currentTab = (m.currentTab + 1) % numTabs
				}
				if m.currentTab == tabMonitor && m.previousTab != tabMonitor && m.isConnected && m.activeSSHClient != nil {
					return m, func() tea.Msg {
						return runMonitorCommand(m.activeSSHClient)
					}
				}
				if m.currentTab == tabFiles && m.previousTab != tabFiles && m.isConnected && m.activeSSHClient != nil {
					return m, func() tea.Msg {
						return runFileListCommand(m.activeSSHClient)
					}
				}
			}
		case "d":
			if m.isConnected {
				if m.activeSession != nil {
					_ = m.activeSession.Close()
					m.activeSession = nil
				}
				
				if m.activeSSHClient != nil {
					_ = m.activeSSHClient.Close()
					m.activeSSHClient = nil
				}

				m.isConnected = false
				m.currentServer = nil
				m.connectionErr = nil
				m.monitorOutput = ""         				
				m.connecting = false
				m.connectingToHost = ""
				m.connectRequestID = ""

				m.currentTab = tabHistory

				return m, nil
			}
		case "r":
			if m.currentScreen == screenMenu && m.currentTab == tabHistory {
				index := m.list.Index()
				if index >= 0 && index < len(m.serverHistory) {
				err := RemoveHistoryItem(index, &m.serverHistory)
				if err != nil {
					m.connectionErr = fmt.Errorf("❌ Failed to remove server: %w", err)
					return m, nil
				}

				m.list.SetItems(ToListItems(m.serverHistory))

				if len(m.serverHistory) == 0 {
					m.list.Select(0)
				} else if index >= len(m.serverHistory) {
					m.list.Select(len(m.serverHistory) - 1)
				} else {
					m.list.Select(index)
				}

				if m.currentServer != nil {
					for _, s := range m.serverHistory {
						if s.Host == m.currentServer.Host {
							goto stillConnected
						}
					}

					m.isConnected = false
					m.activeSSHClient = nil
					m.activeSession = nil
					m.currentServer = nil

				}

				stillConnected:
				/*if index >= len(m.serverHistory) {
					m.list.Select(len(m.serverHistory) - 1)
				} else {
					m.list.Select(index)
				}

				if m.currentServer != nil && m.currentServer.Host == m.serverHistory[index].Host {
					m.isConnected = false
					m.activeSSHClient = nil
					m.activeSession = nil
					m.currentServer = nil
				}*/
			}
			return m, nil
		}
	}
	case sshConnectedMsg:
		if msg.requestID != m.connectRequestID {
			return m, nil
		}
		if m.currentServer != nil && m.currentServer.Host == msg.server.Host {
			m.connecting = false
			m.connectingToHost = ""
			return m, nil
		}
		if m.activeSession != nil {
			_ = m.activeSession.Close()
			m.activeSession = nil
		}
		if m.activeSSHClient != nil {
			_ = m.activeSSHClient.Close()
			m.activeSSHClient = nil
		}

		if msg.err != nil {
			m.connectionErr = msg.err
			m.isConnected = false
			m.connecting = false
			m.connectingToHost = ""
			return m, nil
		}

		m.activeSSHClient = msg.client
		m.activeSession = msg.session
		m.currentServer = &msg.server
		m.isConnected = true
		m.connecting = false
		m.connectingToHost = ""
		return m, nil
	case monitorDataMsg:
		if msg.err != nil {
			m.cpuMetrics = "❌ Failed to fetch: " + msg.err.Error()
			m.memoryMetrics = ""
			m.diskMetrics = ""
		} else {
			m.cpuMetrics = msg.cpu
			m.memoryMetrics = msg.memory
			m.diskMetrics = msg.disk
		}
		return m, nil
	case fileListMsg:
		if msg.err != nil {
			m.fileListItems = nil
			m.fileList.SetItems([]list.Item{})
			m.fileList.Title = "❌ Error loading files"
		} else {
			m.fileListItems = make([]FileListItem, len(msg.items))
			for i, item := range msg.items {
				m.fileListItems[i] = item.(FileListItem)
			}
			m.fileList.SetItems(msg.items)
			m.fileList.Select(0)
		}
		return m, nil
	}

	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch m.currentScreen {
	case screenMenu:
		if m.currentTab == tabHistory {
			m.list, cmd = m.list.Update(msg)
			cmds = append(cmds, cmd)
		} else if m.currentTab == tabFiles {
			m.fileList, cmd = m.fileList.Update(msg)
			cmds = append(cmds, cmd)
		}	
	case screenForm:
		m.form.inputs[m.form.focus], cmd = m.form.inputs[m.form.focus].Update(msg)
		cmds = append(cmds, cmd)

		if m.form.focus == m.form.authMethodInputIndex {
			currentAuthMethod := m.form.inputs[m.form.authMethodInputIndex].Value()
			passwordInput := &m.form.inputs[m.form.passwordInputIndex]
			authPathInput := &m.form.inputs[m.form.authPathInputIndex]

			if currentAuthMethod == "key" {
				authPathInput.Placeholder = "Auth Path (e.g., ~/.ssh/id_rsa)"
				authPathInput.EchoMode = textinput.EchoNormal
				passwordInput.SetValue("")
				passwordInput.Placeholder = ""
				passwordInput.EchoMode = textinput.EchoNormal
			} else if currentAuthMethod == "password" {
				passwordInput.Placeholder = "Password"
				passwordInput.EchoMode = textinput.EchoPassword
				authPathInput.SetValue("")
				authPathInput.Placeholder = ""
				authPathInput.EchoMode = textinput.EchoNormal
			} else {
				passwordInput.SetValue("")
				passwordInput.Placeholder = ""
				passwordInput.EchoMode = textinput.EchoNormal

				authPathInput.SetValue("")
				authPathInput.Placeholder = ""
				authPathInput.EchoMode = textinput.EchoNormal
			}
		}
	}

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	switch m.currentScreen {
	case screenSplash:
		return renderSplashScreen(m)
	case screenMenu:
		return renderMenuScreen(m)
	case screenForm:
		return renderConnectionForm(m)
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
	effectiveWidth := min(targetMaxWidth, m.width-8)
	if effectiveWidth < 60 {
		effectiveWidth = 60
	}
	return effectiveWidth
}

func (m model) getContainerHeight() int {
	const targetMaxHeight = 50
	effectiveHeight := min(targetMaxHeight, m.height-8)
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
