package main

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	menuNextKey = key.NewBinding(key.WithKeys("l"))
	menuPrevKey = key.NewBinding(key.WithKeys("h"))
)

type menubarModel struct {
	menuItems []string
	selected  int
}

func initMenubarModel() menubarModel {
	return menubarModel{
		menuItems: []string{"Home", "Subscriptions", "Episodes", "Search", "Player", "Config"},
		selected:  0,
	}
}

func (m menubarModel) Init() tea.Cmd {
	return nil
}

func (m menubarModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		if key.Matches(msg, quitKeys) {
			return m, tea.Quit
		} else if key.Matches(msg, menuNextKey) {
			m.selectNext()
		} else if key.Matches(msg, menuPrevKey) {
			m.selectPrevious()
		}
		return m, nil
	case errMsg:
		return m, nil

	default:
		return m, nil
	}
}

func (m menubarModel) View() string {
	s := ""
	for i, item := range m.menuItems {
		if i == m.selected {
			s += lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Render(item)
		} else {
			s += lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render(item)
		}
		if i < len(m.menuItems)-1 {
			s += lipgloss.NewStyle().Foreground(lipgloss.Color("177")).Render("|")
		}
	}
	return s
}

func (m *menubarModel) selectNext() {
	if m.selected < len(m.menuItems)-1 {
		m.selected += 1
	}
}

func (m *menubarModel) selectPrevious() {
	if m.selected > 0 {
		m.selected -= 1
	}
}
