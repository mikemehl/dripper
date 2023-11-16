package pages

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type MenuBar struct {
	items    []string
	selected int
}

var (
	unselectedTabStyle = lipgloss.NewStyle().Border(lipgloss.RoundedBorder(), true).Padding(0, 1)
	selectedTabStyle   = lipgloss.NewStyle().Inherit(unselectedTabStyle).UnsetBorderBottom().Underline(true)
)

func InitMenuBar(items []string) MenuBar {
	return MenuBar{
		items:    items,
		selected: 0,
	}
}

func (m MenuBar) Init() tea.Cmd { return nil }
func (m MenuBar) Update(msg tea.Msg) (MenuBar, tea.Cmd) {
	return m, nil
}

func (m MenuBar) View() string {
	bar := ""
	for i, item := range m.items {
		if i == m.selected {
			bar = lipgloss.JoinHorizontal(lipgloss.Left, bar, selectedTabStyle.Render(item))
		} else {
			bar = lipgloss.JoinHorizontal(lipgloss.Left, bar, unselectedTabStyle.Render(item))
		}
	}
	return bar
}
