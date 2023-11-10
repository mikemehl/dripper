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
	menu  MenuBar
	feeds *db.SubData
}

func InitMainModel() tea.Model {
	m := Model{
		state: StateIntro,
		intro: InitIntroPage(),
	}
	m.subs = InitSubsPage(m.feeds)
	m.menu = InitMenuBar([]string{"Subscriptions", "Episodes"})
	return m
}

func (m Model) Init() tea.Cmd {
	return m.intro.Init()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd = nil
	switch m.state {
	case StateIntro:
		switch msg.(type) {
		case db.SubData:
			log.Info("Feed list received by main model")
			msg := msg.(db.SubData)
			m.feeds = &msg
			m.state = StateSubsList
			cmd = func() tea.Msg { return m.feeds }
			return m, cmd
		}
		m.intro, cmd = m.intro.Update(msg)
	case StateSubsList:
		log.Debug("StateSubsList")
		m.menu, cmd = m.menu.Update(msg)
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
		return m.menu.View() + m.subs.View()
	default:
	}
	return lipgloss.NewStyle().Blink(true).Render("DANGER! DANGER! NO VIEW IMPLEMENTED FOR STATE: " + string(m.state))
}
