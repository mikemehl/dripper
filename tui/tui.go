package tui

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	log "github.com/charmbracelet/log"
	"github.com/mikemehl/dripper/db"
	models "github.com/mikemehl/dripper/models"
)

var appBoxStyle = lipgloss.NewStyle().Border(lipgloss.RoundedBorder())

type App struct {
	data     *db.SubData
	loading  spinner.Model
	podcasts models.DetailList
	episodes models.DetailList
	menu     models.Menu
	width    int
	height   int
}

func Run() error {
	log.Debug("Run() called")
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
		podcasts: models.NewDetailList([]list.Item{}, 40, 40),
		episodes: models.NewDetailList([]list.Item{}, 40, 40),
		width:    40,
		height:   40,
		loading:  spinner.New(),
	}
}

func (app App) Init() tea.Cmd {
	log.Debug("I was called")
	return tea.Batch(app.loading.Tick, db.LoadFeeds)
}

func (app App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	log.Debug("Update got ", "msg", msg)
	log.Debug("HELLLLLLOOOO")
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if cmd = app.processKey(msg); cmd != nil {
			return app, cmd
		}
	case *db.SubData:
		log.Debug("Updating app data")
		app.data = msg
		feeds := make([]list.Item, len(app.data.Feeds))
		episodes := make([]list.Item, len(app.data.Episodes))
		for i, feed := range app.data.Feeds {
			feeds[i] = list.Item(feed)
		}
		for i, episode := range app.data.Episodes {
			episodes[i] = list.Item(episode)
		}
		app.podcasts, _ = app.podcasts.Update(app.data.Feeds)
		app.episodes, _ = app.episodes.Update(app.data.Episodes)
	default:
		if app.data == nil {
			app.loading, cmd = app.loading.Update(msg)
			return app, cmd
		}
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
	// return appBoxStyle.MaxWidth(app.width).MaxHeight(app.height).Render(lipgloss.JoinVertical(lipgloss.Top, app.menu.View(), app.podcasts.View()))
	if app.data == nil {
		return app.loading.View()
	}
	return app.podcasts.View()
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
