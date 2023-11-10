package pages

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/mikemehl/dripper/db"
)

type SubData struct {
	feeds    db.FeedList
	episodes db.EpisodeList
}

type SubsPage struct {
	subs SubData
}

func InitSubsPage() SubsPage {
	return SubsPage{
		subs: SubData{
			feeds:    nil,
			episodes: nil,
		},
	}
}

func (m SubsPage) Init() tea.Cmd {
	return nil
}

func (m SubsPage) Update(msg tea.Msg) (SubsPage, tea.Cmd) {
	return m, nil
}

func (m SubsPage) View() string {
	return "SubsPage.View()"
}
