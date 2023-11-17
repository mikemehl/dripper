package models

import (
	"fmt"
	"io"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

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
	detailBoxStyle          = lipgloss.NewStyle().Border(lipgloss.InnerHalfBlockBorder(), false, false, false, true)
	detailKeys              = list.DefaultKeyMap()
)

type DetailList struct {
	list   list.Model
	height int
	width  int
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

func NewDetailList(items []list.Item, width int, height int) DetailList {
	list := list.New(items, DetailListItemDelegate{}, width/2, height)
	list.StartSpinner()
	return DetailList{
		list:   list,
		height: height,
		width:  width,
	}
}

func (d DetailList) Init() tea.Cmd {
	return nil
}

func (d DetailList) Update(msg tea.Msg) (DetailList, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case []list.Item:
		d.list.SetItems(msg)
		d.list.StopSpinner()
	}
	d.list, cmd = d.list.Update(msg)
	return d, cmd
}

func (d DetailList) View() string {
	// details := ""
	// if len(d.list.Items()) > 0 {
	// 	details = d.list.SelectedItem().(DetailListItem).Details()
	// }
	// d.list.SetWidth(d.width / 2)
	// d.list.SetHeight(d.height)
	return d.list.View()
	// return lipgloss.JoinHorizontal(lipgloss.Left,
	// 	d.list.View(), detailBoxStyle.MaxWidth(d.width/2).MaxHeight(d.height).Render(details))
}
