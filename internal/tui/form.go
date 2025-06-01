package tui

import (
	"github.com/charmbracelet/bubbles/textinput"
)

type ConnectionFormModel struct {
    host     string
    port     string
    username string
    method   string
    step     int
}

type ConnectionForm struct {
    inputs []textinput.Model
    focus  int
		passwordInputIndex   int // Index of the "Password" input
		authMethodInputIndex int // Index of the "Auth Method" input
		authPathInputIndex   int // Index of the "Auth Path" input
}
