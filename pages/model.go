package pages

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	log "github.com/charmbracelet/log"
	"github.com/mikemehl/dripper/db"
)

type ModelState int

const (
	StateIntro ModelState = iota
	StateSubsList
)

type Model struct {
	state ModelState
	intro IntroPage
	subs  SubsPage
	feeds *db.FeedList
}

func InitMainModel() tea.Model {
	return Model{
		state: StateIntro,
		intro: InitIntroPage(),
	}
}

func (m Model) Init() tea.Cmd {
	return m.intro.Init()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd = nil
	switch m.state {
	case StateIntro:
		switch msg.(type) {
		case db.FeedList:
			log.Info("Feed list received by main model")
			msg := msg.(db.FeedList)
			m.feeds = &msg
			m.state = StateSubsList
		}
		m.intro, cmd = m.intro.Update(msg)
	case StateSubsList:
		log.Debug("StateSubsList")
		m.subs, cmd = m.subs.Update(msg)
	}
	log.Info("Main Model Update", "state", m.state)
	return m, cmd
}

func (m Model) View() string {
	switch m.state {
	case StateIntro:
		return m.intro.View()
	case StateSubsList:
		return m.subs.View()
	default:
	}
	return lipgloss.NewStyle().Blink(true).Render("DANGER! DANGER! NO VIEW IMPLEMENTED FOR STATE: " + string(m.state))
}
