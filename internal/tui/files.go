package tui

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"golang.org/x/crypto/ssh"
)
type FileListItem struct {
	Name string
	Info string 
}

func (f FileListItem) Title() string       { return f.Name }
func (f FileListItem) Description() string { return f.Info }
func (f FileListItem) FilterValue() string { return f.Name }


func runFileListCommand(client *ssh.Client) tea.Msg {
	session, err := client.NewSession()
	if err != nil {
		return fileListMsg{err: fmt.Errorf("session error: %w", err)}
	}
	defer session.Close()

	var out bytes.Buffer
	session.Stdout = &out

	err = session.Run("ls -lh --color=never --time-style=long-iso")
	if err != nil {
		return fileListMsg{err: fmt.Errorf("run error: %w", err)}
	}

	lines := strings.Split(out.String(), "\n")

	var items []list.Item
	for _, line := range lines[1:] {
		fields := strings.Fields(line)
		if len(fields) < 6 {
			continue
		}
		name := fields[len(fields)-1]
		info := fmt.Sprintf("%s %s %s", fields[0], fields[4], fields[5])
		items = append(items, FileListItem{Name: name, Info: info})
	}

	return fileListMsg{items: items}
}

type fileListMsg struct {
	items []list.Item
	err   error
}
