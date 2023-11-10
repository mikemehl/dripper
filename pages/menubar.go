package pages

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type MenuBar struct {
	items    []string
	selected int
}

type MenuBarStyles struct {
	selected   lipgloss.Style
	unselected lipgloss.Style
	box        lipgloss.Style
}

var topBarStyles = MenuBarStyles{
	selected:   lipgloss.NewStyle().Background(lipgloss.Color("#1A64C5")).Foreground(lipgloss.Color("#F97137")).Padding(0, 2, 0, 2).UnsetColorWhitespace(),
	unselected: lipgloss.NewStyle().Foreground(lipgloss.Color("#DAD9C7")).Padding(0, 2, 0, 2).UnsetColorWhitespace(),
	box:        lipgloss.NewStyle().Padding(1, 0, 1, 0).MaxWidth(250).Border(lipgloss.DoubleBorder(), false, false, true),
}

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
	var bar string
	for i, item := range m.items {
		if i == m.selected {
			bar += topBarStyles.selected.Render("👉 " + item)
		} else {
			bar += topBarStyles.unselected.Render(item)
		}
	}
	return topBarStyles.box.Render(bar) + "\n"
}
