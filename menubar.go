package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type menuBarStyles struct {
	selected   lipgloss.Style
	unselected lipgloss.Style
	box        lipgloss.Style
	separator  string
}

type menuBar struct {
	menuItems []string
	selected  int
	styles    menuBarStyles
}

func initMenubarModel() menuBar {
	border := lipgloss.NewStyle().Align(lipgloss.Center).Border(lipgloss.RoundedBorder(), true, false, false, false).Padding(1).Margin(1)
	itemBase := lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Bold(true)
	return menuBar{
		menuItems: []string{"Subscriptions", "Episodes", "Search", "Player", "Config"},
		selected:  0,
		styles: menuBarStyles{
			selected:   itemBase,
			unselected: lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Inherit(itemBase),
			box:        border,
			separator:  lipgloss.NewStyle().Foreground(lipgloss.Color("177")).Padding(0, 1).Render("|"),
		},
	}
}

func (m menuBar) Init() tea.Cmd {
	return nil
}

func (m *menuBar) Update(msg tea.Msg) tea.Cmd {
	return nil
}

func (m menuBar) View() string {
	s := ""
	for i, item := range m.menuItems {
		if i == m.selected {
			s += m.styles.selected.Render(item)
		} else {
			s += m.styles.unselected.Render(item)
		}
		if i < len(m.menuItems)-1 {
			s += m.styles.separator
		}
	}
	return m.styles.box.Render(s)
}

func (m *menuBar) selectNext() {
	m.selected = (m.selected + 1) % len(m.menuItems)
}

func (m *menuBar) selectPrevious() {
	if m.selected == 0 {
		m.selected = len(m.menuItems) - 1
	} else {
		m.selected = (m.selected - 1) % len(m.menuItems)
	}
}
