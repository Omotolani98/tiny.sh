package main


import (
	"os"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/Omotolani98/tiny.sh/internal/tui"
)

func main() {
	m := tui.NewModel()

	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running program: %v\n", err)
		os.Exit(1)
	}
}
/*m := tui.NewModel()

	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running program: %v\n", err)
    os.Exit(1)
	}*/

