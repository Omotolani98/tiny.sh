package tui

import (
	"encoding/json"
	"fmt"

	"os"
	"path/filepath"

	"github.com/charmbracelet/bubbles/list"
)

type ServerHistoryItem struct {
	Title string `json:"title"`
	Description string `json:"description"`
	Host string `json:"host"`
	Port int `json:"port"`
	Username string `json:"username"`
	AuthMethod string `json:"auth_method"`
	Auth AuthDetails `json:"auth"`
	LastConnected string `json:"last_connected"`
}

type AuthDetails struct {
	Method string `json:"method"`
	Path   string `json:"path,omitempty"`    
	Value  string `json:"value,omitempty"`
}

type ServerListItem struct {
	Server ServerHistoryItem
}

func (i ServerListItem) Title() string {
	return i.Server.Title
}

func (i ServerListItem) Description() string {
	return fmt.Sprintf("Host: %s, Port: %d, User: %s, Auth: %s, Last Connected: %s",
		i.Server.Host, i.Server.Port, i.Server.Username, i.Server.Auth.Method, i.Server.LastConnected)
}

func (i ServerListItem) FilterValue() string {
	return i.Server.Title + " " + i.Server.Host + " " + i.Server.Username + " " + i.Server.Auth.Method
}

func LoadHistory() ([]ServerHistoryItem, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("could not get user home directory: %w", err)
	}
	historyDir := filepath.Join(home, ".tiny")
	path := filepath.Join(historyDir, "history.json")

	content, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return []ServerHistoryItem{}, nil
		}
		return nil, fmt.Errorf("could not read history file: %w", err)
	}

	var servers []ServerHistoryItem
	if err := json.Unmarshal(content, &servers); err != nil {	
		return nil, fmt.Errorf("could not unmarshal history JSON: %w", err)
	}
	
	return servers, nil
}

func SaveHistory(servers []ServerHistoryItem) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("could not get user home directory: %w", err)
	}
	historyDir := filepath.Join(home, ".tiny")
	path := filepath.Join(historyDir, "history.json")

	if err := os.MkdirAll(historyDir, 0755); err != nil {
		return fmt.Errorf("could not create history directory: %w", err)
	}

	data, err := json.MarshalIndent(servers, "", "  ") 	
	if err != nil {
		return fmt.Errorf("could not marshal history to JSON: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("could not write history file: %w", err)
	}

	return nil
}

func ToListItems(history []ServerHistoryItem) []list.Item {
    var items []list.Item
    for _, s := range history {
			items = append(items, ServerListItem{Server: s})
    }
    return items
}

func RemoveHistoryItem(index int, history *[]ServerHistoryItem) error {
	if index < 0 || index >= len(*history) {
		return fmt.Errorf("invalid index %d", index)
	}

	removed := (*history)[index]

	*history = append((*history)[:index], (*history)[index+1:]...)

	if err := SaveHistory(*history); err != nil {
		return fmt.Errorf("failed to save updated history: %w", err)
	}

	return AppendToBackup(removed)
}

func AppendToBackup(item ServerHistoryItem) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("could not get user home directory: %w", err)
	}
	historyDir := filepath.Join(home, ".tiny")
	backupPath := filepath.Join(historyDir, "history_backup.json")

	var backup []ServerHistoryItem

	if _, err := os.Stat(backupPath); err == nil {
		data, err := os.ReadFile(backupPath)
		if err == nil {
			_ = json.Unmarshal(data, &backup) 
		}
	}

	backup = append(backup, item)

	data, err := json.MarshalIndent(backup, "", "  ")
	if err != nil {
		return fmt.Errorf("could not marshal backup: %w", err)
	}

	if err := os.WriteFile(backupPath, data, 0644); err != nil {
		return fmt.Errorf("could not write backup file: %w", err)
	}

	return nil
}
