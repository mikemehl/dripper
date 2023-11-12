package pages

import (
	"fmt"
	"io"

	key "github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	log "github.com/charmbracelet/log"
	"github.com/mikemehl/dripper/db"
)

const (
	ListWidth   = 40
	ListHeight  = 20
	InputWidth  = (40 / 4)
	InputHeight = (20 / 4)
)

type SubsPage struct {
	subs  *db.SubData
	list  list.Model
	input textinput.Model
}

type subItemDelegate struct{}

func (d subItemDelegate) Height() int                             { return 1 }
func (d subItemDelegate) Spacing() int                            { return 1 }
func (d subItemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d subItemDelegate) Render(w io.Writer, m list.Model, idx int, item list.Item) {
	style := lipgloss.NewStyle().Inline(true).Faint(true)
	if idx == m.Cursor() {
		style = style.Faint(false).Foreground(lipgloss.Color("#F97137"))
	}
	fmt.Fprintf(w, "%s\n", style.Render(item.(db.Feed).Title))
}

func initList(delegate list.ItemDelegate) list.Model {
	l := list.New([]list.Item{}, delegate, ListWidth, ListHeight)
	l.SetShowTitle(false)
	l.SetShowStatusBar(true)
	l.SetFilteringEnabled(false)
	l.KeyMap.Quit.SetKeys(KeyBindings.Quit.Keys()...)
	l.KeyMap.CursorUp.SetKeys(KeyBindings.Up.Keys()...)
	l.KeyMap.CursorDown.SetKeys(KeyBindings.Down.Keys()...)
	l.InfiniteScrolling = true
	return l
}

func initInput() textinput.Model {
	input := textinput.New()
	input.Prompt = "Add subscription: "
	input.Blur()
	return input
}

func InitSubsPage(subs *db.SubData) SubsPage {
	var m SubsPage
	m.subs = subs
	m.list = initList(subItemDelegate{})
	m.input = initInput()
	return m
}

func (m SubsPage) Init() tea.Cmd {
	return nil
}

func (m SubsPage) UpdateAddingSub(msg tea.Msg) (SubsPage, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, KeyBindings.Confirm):
			var cmds []tea.Cmd
			m.input.Blur()
			url := m.input.Value()
			cmds = append(cmds, m.list.NewStatusMessage("Adding "+url))
			cmds = append(cmds, m.list.StartSpinner())
			m.list.SetShowStatusBar(true)
			cmd = func() tea.Msg {
				feed, err := db.NewFeed(m.input.Value())
				if err != nil {
					return err
				}
				m.subs.LoadFeed(feed)
				return m.subs
			}
			cmds = append(cmds, cmd)
			return m, tea.Batch(cmds...)
		default:
		}
	}
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func (m SubsPage) Update(msg tea.Msg) (SubsPage, tea.Cmd) {
	if m.input.Focused() {
		return m.UpdateAddingSub(msg)
	}
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, KeyBindings.Quit):
			return m, tea.Quit
		case key.Matches(msg, KeyBindings.AddSub):
			m.input.Focus()
		case key.Matches(msg, KeyBindings.Confirm):
			feed := m.list.SelectedItem().(db.Feed)
			return m, func() tea.Msg { return &feed }
		default:
		}
	case *db.SubData:
		log.Debug("SubData recevied by SubsPage")
		m.subs = msg
		newItems := newSubListItems(msg)
		m.list.SetItems(newItems)
	default:
	}
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m SubsPage) View() string {
	view := m.list.View()
	if m.input.Focused() {
		view = lipgloss.Place(InputWidth, InputHeight, lipgloss.Center, lipgloss.Left, m.input.View(), lipgloss.WithWhitespaceChars(view))
	}
	return view
}

func newSubListItems(data *db.SubData) []list.Item {
	log.Info("NewItems() called")
	var items []list.Item
	for _, feed := range (data).Feeds {
		log.Debug("Adding item", "title", feed.Title)
		items = append(items, list.Item(feed))
	}
	return items
}
