package models

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mikemehl/dripper/utils"
)

type MenuItem struct {
	Action tea.Cmd
	Name   string
}

type Menu struct {
	items     []MenuItem
	active    int
	itemWidth int
	style     lipgloss.Style
}

type MenuSetActive struct {
	Index int
}

var (
	menuItemStyle       = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Foreground(lipgloss.Color("#FFB86C")).Align(lipgloss.Center)
	menuItemActiveStyle = menuItemStyle.Copy().Foreground(lipgloss.Color("#FF79C6")).Underline(true)
	menuStyle           = lipgloss.NewStyle()
)

func MenuSetActiveCmd(index int) tea.Cmd {
	return func() tea.Msg { return MenuSetActive{Index: index} }
}

func NewMenu(items []MenuItem) tea.Model {
	width := 0
	for _, item := range items {
		if len(item.Name) > width {
			width = len(item.Name)
		}
	}
	width += 2
	return tea.Model(Menu{
		items:     items,
		active:    0,
		itemWidth: width,
		style:     menuStyle,
	})
}

func (m Menu) Init() tea.Cmd {
	return nil
}

func (m Menu) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "H":
			if m.active > 0 {
				m.active--
			}
		case "L":
			if m.active < len(m.items)-1 {
				m.active++
			}
		}
	case tea.WindowSizeMsg:
		_ = m.SetDimensions(msg)
	case MenuSetActive:
		m.setActive(msg.Index)
	}
	return m, nil
}

func (m Menu) View() string {
	view := ""
	for i, item := range m.items {
		if i == m.active {
			view = lipgloss.JoinHorizontal(lipgloss.Top, view, menuItemActiveStyle.Width(m.itemWidth).Align(lipgloss.Center, lipgloss.Center).Render(item.Name))
		} else {
			view = lipgloss.JoinHorizontal(lipgloss.Top, view, menuItemStyle.Width(m.itemWidth).Align(lipgloss.Center, lipgloss.Center).Render(item.Name))
		}
	}
	return menuStyle.Render(view)
}

func (m *Menu) setActive(index int) {
	if index < 0 || index >= len(m.items) {
		return
	}
	m.active = index
}

func (m Menu) Active() int {
	return m.active
}

func (m *Menu) SetDimensions(msg tea.WindowSizeMsg) tea.WindowSizeMsg {
	msg = utils.ScaleDimensions(msg, 10, 9, 2)
	m.style = m.style.Width(msg.Width).Height(msg.Height).Align(lipgloss.Left, lipgloss.Top)
	return msg
}
