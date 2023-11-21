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

type (
	DetailListAction func(d DetailList) tea.Cmd
	SpinnerCmd       struct{ Active bool }
	MessageCmd       struct {
		Msg    string
		Active bool
	}
)

var (
	detailListItemStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFB86C"))
	detailListSelectedStyle = detailListItemStyle.Copy().Foreground(lipgloss.Color("#FF79C6")).Bold(true)
	pointerStyle            = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF79C6"))
	detailBoxStyle          = lipgloss.NewStyle().Border(lipgloss.HiddenBorder(), false, false, false, true).Padding(1)
	detailKeys              = list.DefaultKeyMap()
	mdConverter             = md.NewConverter("", true, nil)
	DefaultSelectAction     = func(d DetailList) tea.Cmd { return nil }
)

type DetailList struct {
	list         list.Model
	selectAction DetailListAction
	details      viewport.Model
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

func extraKeys() []key.Binding {
	return []key.Binding{
		key.NewBinding(
			key.WithKeys("H"),
			key.WithDisabled(),
			key.WithHelp("H", "Previous tab")),
		key.NewBinding(
			key.WithKeys("L"),
			key.WithDisabled(),
			key.WithHelp("L", "Next tab")),
		key.NewBinding(
			key.WithKeys("J"),
			key.WithDisabled(),
			key.WithHelp("J", "Scroll description down")),
		key.NewBinding(
			key.WithKeys("K"),
			key.WithDisabled(),
			key.WithHelp("K", "Scroll description up")),
		key.NewBinding(
			key.WithKeys("enter"),
			key.WithDisabled(),
			key.WithHelp("enter", "Select")),
	}
}

func NewDetailList(items []list.Item, width int, height int, action DetailListAction) tea.Model {
	list := list.New(items, DetailListItemDelegate{}, width/2, height)
	list.SetShowTitle(false)
	list.SetFilteringEnabled(true)
	list.SetShowHelp(true)
	list.DisableQuitKeybindings()
	list.AdditionalShortHelpKeys = extraKeys
	details := viewport.New(width/2, height)
	return DetailList{
		list:         list,
		details:      details,
		selectAction: action,
	}
}

func (d DetailList) Init() tea.Cmd {
	return func() tea.Msg { return d.list.StartSpinner() }
}

func (d DetailList) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd = nil
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
		case "enter":
			cmd = d.selectAction(d)
		}
	case SpinnerCmd:
		if msg.Active {
			d.list.StartSpinner()
		} else {
			d.list.StopSpinner()
		}
	case MessageCmd:
		d.list.SetShowStatusBar(msg.Active)
		if msg.Active {
			d.list.NewStatusMessage(msg.Msg)
		}
	}
	return d.UpdatePanels(msg, cmd)
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

func (d DetailList) UpdatePanels(msg tea.Msg, extra tea.Cmd) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	d.list, cmd = d.list.Update(msg)
	cmds = append(cmds, cmd)
	d.UpdateDetails()
	d.details, cmd = d.details.Update(msg)
	cmds = append(cmds, cmd)
	if extra != nil {
		cmds = append(cmds, extra)
	}
	return d, tea.Batch(cmds...)
}

func (d DetailList) SelectedItem() list.Item {
	return d.list.SelectedItem()
}
