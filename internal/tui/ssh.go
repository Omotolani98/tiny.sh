package tui

import (
	"bytes"
	"fmt"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"golang.org/x/crypto/ssh"
)


type sshConnectedMsg struct {
	client  *ssh.Client
	session *ssh.Session
	server  ServerHistoryItem
	err     error
	requestID string
}

type monitorDataMsg struct {
	cpu    string
	memory string
	disk   string
	err    error
}

func trySSHConnect(server ServerHistoryItem, requestID string) tea.Msg {
	config := &ssh.ClientConfig{
		User:            server.Username,
		Auth:            []ssh.AuthMethod{},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // ⚠️ TODO: replace for prod
		Timeout:         5 * time.Second,
	}

	if server.Auth.Method == "password" {
		config.Auth = append(config.Auth, ssh.Password(server.Auth.Value))
	} else if server.Auth.Method == "key" {
		key, err := os.ReadFile(server.Auth.Path)
		if err != nil {
			return sshConnectedMsg{err: fmt.Errorf("read key error: %w", err)}
		}
		signer, err := ssh.ParsePrivateKey(key)
		if err != nil {
			return sshConnectedMsg{err: fmt.Errorf("parse key error: %w", err)}
		}
		config.Auth = append(config.Auth, ssh.PublicKeys(signer))
	} else {
		return sshConnectedMsg{err: fmt.Errorf("unsupported auth method: %s", server.Auth.Method)}
	}

	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", server.Host, server.Port), config)
	if err != nil {
		return sshConnectedMsg{err: fmt.Errorf("ssh dial error: %w", err)}
	}

	session, err := client.NewSession()
	if err != nil {
		client.Close()
		return sshConnectedMsg{err: fmt.Errorf("ssh session error: %w", err)}
	}

	return sshConnectedMsg{
		client:  client,
		session: session,
		server:  server,
		err:     nil,
		requestID: requestID,
	}
}

func runMonitorCommand(client *ssh.Client) tea.Msg {
	var cpuOut, memOut, diskOut bytes.Buffer

	// CPU
	session, err := client.NewSession()
	if err != nil {
		return monitorDataMsg{err: err}
	}
	session.Stdout = &cpuOut
	err = session.Run(`top -bn1 | grep "load average"; top -bn1 | grep "Cpu(s)"`)
	session.Close()
	if err != nil {
		return monitorDataMsg{err: err}
	}

	// Memory
	memSession, err := client.NewSession()
	if err != nil {
		return monitorDataMsg{err: err}
	}
	memSession.Stdout = &memOut
	err = memSession.Run(`free -m | awk 'NR==2{printf "Used: %sMB / Total: %sMB", $3, $2}'`)
	memSession.Close()
	if err != nil {
		return monitorDataMsg{err: err}
	}

	// Disk
	diskSession, err := client.NewSession()
	if err != nil {
		return monitorDataMsg{err: err}
	}
	diskSession.Stdout = &diskOut
	err = diskSession.Run(`df -h / | awk 'NR==2{printf "Used: %s / Total: %s (%s)", $3, $2, $5}'`)
	diskSession.Close()
	if err != nil {
		return monitorDataMsg{err: err}
	}

	return monitorDataMsg{
		cpu:    cpuOut.String(),
		memory: memOut.String(),
		disk:   diskOut.String(),
	}
}
