package pages

import (
	"fmt"
	"io"

	key "github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	log "github.com/charmbracelet/log"
	"github.com/mikemehl/dripper/db"
)

type SubEpsPage struct {
	subTitle string
	list     list.Model
	// TODO: Add a viewport for the details panel, allow scrolling with captial letters.
	// TODO: Maybe extract this logic to a new component? You're reusing it in subs.go.
	showDetails bool
}

type epItemDelegate struct{}

func (d epItemDelegate) Height() int                             { return 1 }
func (d epItemDelegate) Spacing() int                            { return 0 }
func (d epItemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d epItemDelegate) Render(w io.Writer, m list.Model, idx int, item list.Item) {
	style := lipgloss.NewStyle().Inline(true).Faint(true)
	if idx == m.Index() {
		style = style.Faint(false).Foreground(lipgloss.Color("#F97137"))
	}
	fmt.Fprintf(w, "%s\n", style.MaxWidth(CurrWidth/2-2).Render(item.(db.Episode).Title))
}

func InitSubepsPage() SubEpsPage {
	var m SubEpsPage
	m.subTitle = ""
	m.list = initList(epItemDelegate{})
	m.list.SetShowTitle(true)
	return m
}

func (m SubEpsPage) Init() tea.Cmd {
	return nil
}

func (m SubEpsPage) Update(msg tea.Msg) (SubEpsPage, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, KeyBindings.Quit):
			return m, tea.Quit
		case key.Matches(msg, KeyBindings.Back):
			return m, func() tea.Msg { return StateChangeMsg{newState: StateSubsList} }
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
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m SubEpsPage) View() string {
	panel_width := CurrWidth / 2
	panel_height := CurrHeight - 15
	details := ""
	if selected := m.list.SelectedItem(); selected != nil {
		details = selected.(db.Episode).Description
	}
	details = detailsStyle.Width(panel_width).Height(panel_height).MaxWidth(panel_width).MaxHeight(panel_height).Render(details)
	m.list.SetWidth(panel_width)
	m.list.SetHeight(panel_height)
	return lipgloss.JoinHorizontal(lipgloss.Left, m.list.View(), details)
}

func newSubEpsListItems(data *db.Feed) []list.Item {
	var items []list.Item
	for _, ep := range (data).Items {
		log.Debug("Adding item", "title", ep.Title)
		items = append(items, list.Item(db.Episode(*ep)))
	}
	return items
}
