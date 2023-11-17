package tui

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	log "github.com/charmbracelet/log"
	"github.com/mikemehl/dripper/db"
	models "github.com/mikemehl/dripper/models"
)

var appBoxStyle = lipgloss.NewStyle().Border(lipgloss.RoundedBorder())

type App struct {
	data     *db.SubData
	podcasts models.DetailList
	episodes models.DetailList
	menu     models.Menu
}

func Run() error {
	var app tea.Model = NewApp()
	program := tea.NewProgram(app)
	if _, err := program.Run(); err != nil {
		return err
	}
	return nil
}

func NewApp() tea.Model {
	return App{
		menu: models.NewMenu([]models.MenuItem{
			{Name: "Podcasts", Action: nil},
			{Name: "Episodes", Action: nil},
		}),
		data:     nil,
		podcasts: models.NewDetailList([]list.Item{}, 80, 80),
		episodes: models.NewDetailList([]list.Item{}, 80, 80),
	}
}

func (app App) Init() tea.Cmd {
	log.Debug("I was called")
	var cmd tea.Cmd = func() tea.Msg {
		log.Debug("Loading feeds")
		return tea.Msg(db.LoadFeeds())
	}

	log.Debug(cmd)

	return cmd
}

func (app App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if cmd = app.processKey(msg); cmd != nil {
			return app, cmd
		}
	case *db.SubData:
		log.Debug("Updating app data")
		app.data = msg
		app.podcasts, _ = app.podcasts.Update(app.data.Feeds)
		app.episodes, _ = app.episodes.Update(app.data.Episodes)
	}
	if app.menu, cmd = app.menu.Update(msg); cmd != nil {
		return app, cmd
	}
	if app.podcasts, cmd = app.podcasts.Update(msg); cmd != nil {
		return app, cmd
	}
	if app.episodes, cmd = app.episodes.Update(msg); cmd != nil {
		return app, cmd
	}
	return app, cmd
}

func (app App) View() string {
	return appBoxStyle.Render(app.menu.View())
}

func (app *App) processKey(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return tea.Quit
		}
	}
	return nil
}
