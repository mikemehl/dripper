package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

type errMsg error

func main() {
	p := tea.NewProgram(initModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
