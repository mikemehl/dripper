package models

import (
	"fmt"
	"io"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"

	list "github.com/charmbracelet/bubbles/list"
)

type DetailListItem interface {
	Name() string
	Details() string
	FilterValue() string
}

var (
	detailListItemStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFB86C")).Align(lipgloss.Left)
	detailListSelectedStyle = detailListItemStyle.Copy().Foreground(lipgloss.Color("#FF79C6")).Bold(true)
	detailBoxStyle          = lipgloss.NewStyle().Border(lipgloss.HiddenBorder(), false, false, false, true).Padding(1)
	detailKeys              = list.DefaultKeyMap()
)

type DetailList struct {
	list        list.Model
	detailStyle lipgloss.Style
}

type DetailListItemDelegate struct{}

// Render renders the item's view.
func (d DetailListItemDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	style := detailListItemStyle
	if index == m.Index() {
		style = detailListSelectedStyle
	}
	fmt.Fprintf(w, "%s", style.Render(item.(DetailListItem).Name()))
}

// Height is the height of the list item.
func (d DetailListItemDelegate) Height() int { return 1 }

// Spacing is the size of the horizontal gap between list items in cells.
func (d DetailListItemDelegate) Spacing() int { return 0 }

// Update is the update loop for items. All messages in the list's update
// loop will pass through here except when the user is setting a filter.
// Use this method to perform item-level updates appropriate to this
// delegate.
func (d DetailListItemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }

func NewDetailList(items []list.Item, width int, height int) tea.Model {
	list := list.New(items, DetailListItemDelegate{}, width/2, height)
	list.SetShowTitle(false)
	list.SetFilteringEnabled(true)
	list.SetShowHelp(true)
	detailStyle := detailBoxStyle.Width(width / 2)
	detailStyle = detailStyle.Height(height)
	return DetailList{
		list:        list,
		detailStyle: detailStyle,
	}
}

func (d DetailList) Init() tea.Cmd {
	return func() tea.Msg { return d.list.StartSpinner() }
}

func (d DetailList) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		d.SetDimensions(msg.Width, msg.Height)
	case []list.Item:
		log.Debug("DetailList.Update() []list.Item")
		d.list.SetItems(msg)
		d.list.StopSpinner()
	}
	d.list, cmd = d.list.Update(msg)
	return d, cmd
}

func (d DetailList) View() string {
	details := ""
	if len(d.list.Items()) > 0 {
		details = d.list.SelectedItem().(DetailListItem).Details()
	}
	return lipgloss.JoinHorizontal(lipgloss.Left,
		d.list.View(), d.detailStyle.Render(details))
}

func (d *DetailList) SetDimensions(width int, height int) {
	width = width / 2
	height = height / 10 * 7
	d.list.SetWidth(width)
	d.list.SetHeight(height)
	d.list.SetSize(width, height)
	d.detailStyle = d.detailStyle.Width(width).MaxWidth(width).Height(height).MaxHeight(height)
}
