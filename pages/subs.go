package pages

import (
	"fmt"
	"io"

	key "github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	log "github.com/charmbracelet/log"
	"github.com/mikemehl/dripper/db"
)

type SubsPage struct {
	subs *db.SubData
	list list.Model
}

type itemDelegate struct{}

func (d itemDelegate) Height() int                             { return 5 }
func (d itemDelegate) Spacing() int                            { return 5 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, idx int, item list.Item) {
	fmt.Fprintf(w, "%s\n", item.(db.Feed).Title)
}

func initList() list.Model {
	l := list.New([]list.Item{}, itemDelegate{}, 20, 14)
	l.SetShowTitle(false)
	l.SetShowStatusBar(true)
	l.SetFilteringEnabled(false)
	return l
}

func InitSubsPage(subs *db.SubData) SubsPage {
	var m SubsPage
	m.subs = subs
	m.list = initList()
	return m
}

func (m SubsPage) Init() tea.Cmd {
	return nil
}

func (m SubsPage) Update(msg tea.Msg) (SubsPage, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, KeyBindings.Quit):
			return m, tea.Quit
		default:
		}
	case *db.SubData:
		log.Debug("SubData recevied by SubsPage")
		m.subs = msg
		newItems := newListItems(msg)
		m.list.SetItems(newItems)
	default:
	}
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m SubsPage) View() string {
	return m.list.View()
}

func newListItems(data *db.SubData) []list.Item {
	log.Info("NewItems() called")
	var items []list.Item
	for _, feed := range (data).Feeds {
		log.Debug("Adding item", "title", feed.Title)
		items = append(items, list.Item(feed))
	}
	return items
}
