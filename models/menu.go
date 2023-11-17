package models

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type MenuItem struct {
	Action tea.Cmd
	Name   string
}

type Menu struct {
	items     []MenuItem
	active    int
	itemWidth int
}

var (
	menuItemStyle       = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Foreground(lipgloss.Color("#FFB86C")).Align(lipgloss.Center)
	menuItemActiveStyle = menuItemStyle.Copy().Foreground(lipgloss.Color("#FF79C6")).Underline(true)
	menuStyle           = lipgloss.NewStyle()
)

func NewMenu(items []MenuItem) Menu {
	width := 0
	for _, item := range items {
		if len(item.Name) > width {
			width = len(item.Name)
		}
	}
	width += 2
	return Menu{
		items:     items,
		active:    0,
		itemWidth: width,
	}
}

func (m Menu) Init() tea.Cmd {
	return nil
}

func (m Menu) Update(msg tea.Msg) (Menu, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "left", "h":
			if m.active > 0 {
				m.active--
			}
		case "right", "l":
			if m.active < len(m.items)-1 {
				m.active++
			}
		}
	}
	return m, nil
}

func (m Menu) View() string {
	view := ""
	for i, item := range m.items {
		if i == m.active {
			view = lipgloss.JoinHorizontal(lipgloss.Top, view, menuItemActiveStyle.Width(m.itemWidth).Render(item.Name))
		} else {
			view = lipgloss.JoinHorizontal(lipgloss.Top, view, menuItemStyle.Width(m.itemWidth).Render(item.Name))
		}
	}
	return menuStyle.Render(view)
}
