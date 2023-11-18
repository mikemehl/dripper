package models

import (
	"fmt"
	"io"

	md "github.com/JohannesKaufmann/html-to-markdown"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/mikemehl/dripper/utils"

	"github.com/charmbracelet/bubbles/key"
	list "github.com/charmbracelet/bubbles/list"
	viewport "github.com/charmbracelet/bubbles/viewport"
)

type DetailListItem interface {
	Name() string
	Details() string
	FilterValue() string
}

var (
	detailListItemStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFB86C"))
	detailListSelectedStyle = detailListItemStyle.Copy().Foreground(lipgloss.Color("#FF79C6")).Bold(true)
	pointerStyle            = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF79C6"))
	detailBoxStyle          = lipgloss.NewStyle().Border(lipgloss.HiddenBorder(), false, false, false, true).Padding(1)
	detailKeys              = list.DefaultKeyMap()
	mdConverter             = md.NewConverter("", true, nil)
)

type DetailList struct {
	list    list.Model
	details viewport.Model
}

type DetailListItemDelegate struct{}

// Render renders the item's view.
func (d DetailListItemDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	style := detailListItemStyle
	name := item.(DetailListItem).Name()
	if index == m.Index() {
		style = detailListSelectedStyle
		name = fmt.Sprintf("%s %s", "▶", name)
	}
	fmt.Fprintf(w, "%s", style.Width(m.Width()-1).Align(lipgloss.Left, lipgloss.Center).Render(name))
}

// Height is the height of the list item.
func (d DetailListItemDelegate) Height() int { return 1 }

// Spacing is the size of the horizontal gap between list items in cells.
func (d DetailListItemDelegate) Spacing() int { return 1 }

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
	list.DisableQuitKeybindings()
	list.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{
			key.NewBinding(key.WithHelp("H", "Previous tab")),
			key.NewBinding(key.WithHelp("L", "Next tab")),
			key.NewBinding(key.WithHelp("J", "Scroll description down")),
			key.NewBinding(key.WithHelp("K", "Scroll description up")),
		}
	}
	details := viewport.New(width/2, height)
	return DetailList{
		list:    list,
		details: details,
	}
}

func (d DetailList) Init() tea.Cmd {
	return func() tea.Msg { return d.list.StartSpinner() }
}

func (d DetailList) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		d.SetDimensions(msg)
	case []list.Item:
		log.Debug("DetailList.Update() []list.Item")
		d.list.SetItems(msg)
		d.list.StopSpinner()
	case tea.KeyMsg:
		switch msg.String() {
		case "J":
			d.details.HalfViewDown()
		case "K":
			d.details.HalfViewUp()
		}
	}
	return d.UpdatePanels(msg)
}

func (d DetailList) View() string {
	return lipgloss.JoinHorizontal(lipgloss.Center,
		d.list.View(), d.details.View())
}

func (d *DetailList) SetDimensions(msg tea.WindowSizeMsg) {
	msg = utils.ScaleDimensions(msg, 10, 4, 9)
	d.list.SetSize(msg.Width, msg.Height)
	d.details = viewport.New(msg.Width, msg.Height)
}

func (d *DetailList) UpdateDetails() {
	details := ""
	if len(d.list.Items()) > 0 {
		details = d.list.SelectedItem().(DetailListItem).Details()
		if mkdn, err := mdConverter.ConvertString(details); err == nil {
			if fancy, err := glamour.Render(mkdn, "dark"); err == nil {
				details = fancy
			}
		}
	}
	d.details.SetContent(details)
}

func (d DetailList) UpdatePanels(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	d.list, cmd = d.list.Update(msg)
	cmds = append(cmds, cmd)
	d.UpdateDetails()
	d.details, cmd = d.details.Update(msg)
	cmds = append(cmds, cmd)
	return d, tea.Batch(cmds...)
}
