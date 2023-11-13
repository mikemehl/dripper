package pages

import (
	"fmt"
	"io"

	key "github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	log "github.com/charmbracelet/log"
	"github.com/mikemehl/dripper/db"
)

type SubEpsPage struct {
	subTitle    string
	list        list.Model
	details     viewport.Model
	showDetails bool
}

type epItemDelegate struct{}

func (d epItemDelegate) Height() int                             { return 1 }
func (d epItemDelegate) Spacing() int                            { return 1 }
func (d epItemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d epItemDelegate) Render(w io.Writer, m list.Model, idx int, item list.Item) {
	style := lipgloss.NewStyle().Inline(true).Faint(true)
	if idx == m.Index() {
		style = style.Faint(false).Foreground(lipgloss.Color("#F97137"))
	}
	fmt.Fprintf(w, "%s\n", style.Render(item.(db.Episode).Title))
}

func InitSubepsPage() SubEpsPage {
	var m SubEpsPage
	m.subTitle = ""
	m.list = initList(epItemDelegate{})
	m.details.SetContent("")
	m.details.Width = m.list.Width()
	m.details.Height = m.list.Height()
	m.list.SetShowTitle(true)
	m.showDetails = false
	return m
}

func (m SubEpsPage) Init() tea.Cmd {
	return nil
}

func (m SubEpsPage) SetDetails(selected db.Episode) {
	content := selected.Title + "\n\n" + selected.Description
	content = lipgloss.NewStyle().Width(m.list.Width()).Height(m.list.Height()).Render(content)
	m.details.SetContent(content)
}

func (m SubEpsPage) Update(msg tea.Msg) (SubEpsPage, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, KeyBindings.Quit):
			return m, tea.Quit
		case key.Matches(msg, KeyBindings.Back):
			if m.showDetails {
				m.showDetails = false
				return m, nil
			}
			return m, func() tea.Msg { return StateChangeMsg{newState: StateSubsList} }
		case key.Matches(msg, KeyBindings.Confirm):
			if !m.showDetails {
				m.showDetails = true
				m.SetDetails(m.list.SelectedItem().(db.Episode))
				return m, nil
			}
		default:
		}
	case *db.Feed:
		log.Debug("SubData recevied by SubEpsPage")
		m.subTitle = msg.Title
		m.list.Title = msg.Title
		newItems := newSubEpsListItems(msg)
		m.list.SetItems(newItems)
	default:
	}
	if m.showDetails {
		m.details, cmd = m.details.Update(msg)
	} else {
		m.list, cmd = m.list.Update(msg)
	}
	return m, cmd
}

func (m SubEpsPage) View() string {
	return lipgloss.Place(CurrWidth/2, CurrHeight/2, lipgloss.Center, lipgloss.Top, m.details.View(), lipgloss.WithWhitespaceChars(m.list.View()))
}

func newSubEpsListItems(data *db.Feed) []list.Item {
	var items []list.Item
	for _, ep := range (data).Items {
		log.Debug("Adding item", "title", ep.Title)
		items = append(items, list.Item(db.Episode(*ep)))
	}
	return items
}
