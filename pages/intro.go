package pages

import (
	"fmt"

	key "github.com/charmbracelet/bubbles/key"
	spinner "github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	log "github.com/charmbracelet/log"
	"github.com/mikemehl/dripper/db"
)

type IntroPage struct {
	spinner           spinner.Model
	feedLoadInitiated bool
	feeds             *db.FeedList
}

func InitIntroPage() IntroPage {
	m := IntroPage{
		spinner:           spinner.New(),
		feedLoadInitiated: false,
		feeds:             nil,
	}
	m.spinner.Spinner = spinner.Points
	return m
}

func (m IntroPage) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m IntroPage) Update(msg tea.Msg) (IntroPage, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, KeyBindings.Quit):
			return m, tea.Quit
		default:
		}
	case db.FeedList:
		log.Warn("Feed list received by intro page")
	default:
	}

	var cmd tea.Cmd
	if !m.feedLoadInitiated {
		cmd = db.LoadFeeds
		m.feedLoadInitiated = true
	} else {
		m.spinner, cmd = m.spinner.Update(msg)
	}
	log.Debug("IntroPage.Update() returns ", "cmd", cmd)
	return m, cmd
}

func (m IntroPage) View() string {
	log.Debug("IntroPage.View()")
	s := fmt.Sprintf("%s Loading feeds", m.spinner.View())
	return s
}
