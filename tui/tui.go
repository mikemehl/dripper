package tui

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	log "github.com/charmbracelet/log"
	"github.com/mikemehl/dripper/db"
	models "github.com/mikemehl/dripper/models"
)

var appBoxStyle = lipgloss.NewStyle().Border(lipgloss.DoubleBorder()).Padding(1)

type App struct {
	style    lipgloss.Style
	data     *db.SubData
	podcasts tea.Model
	episodes tea.Model
	active   *tea.Model
	menu     tea.Model
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
	podcasts := models.NewDetailList([]list.Item{}, 40, 40)
	episodes := models.NewDetailList([]list.Item{}, 40, 40)
	return App{
		menu: models.NewMenu([]models.MenuItem{
			{Name: "Podcasts", Action: nil},
			{Name: "Episodes", Action: nil},
		}),
		data:     nil,
		podcasts: podcasts,
		episodes: episodes,
		style:    appBoxStyle,
	}
}

func (app App) Init() tea.Cmd {
	log.Debug("I was called")
	return db.LoadFeeds
}

func (app App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		msg = app.SetDimensions(msg)
	case tea.KeyMsg:
		if cmd = app.processKey(msg); cmd != nil {
			return app, cmd
		}
	case *db.SubData:
		return app.UpdateFeeds(msg)
	default:
	}
	return app.UpdateSubModels(msg)
}

func (app App) View() string {
	return app.style.Render(lipgloss.JoinVertical(lipgloss.Top, app.menu.View(), app.podcasts.View()))
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

func ScaleDimensions(msg tea.WindowSizeMsg, bottom int, wtop int, htop int) tea.WindowSizeMsg {
	msg.Width = msg.Width / 10 * wtop
	msg.Height = msg.Height / 10 * htop
	return msg
}

func (app *App) SetDimensions(msg tea.WindowSizeMsg) tea.WindowSizeMsg {
	msg = ScaleDimensions(msg, 10, 9, 9)
	app.style = app.style.Width(msg.Width).Height(msg.Height)
	msg = ScaleDimensions(msg, 10, 9, 7)
	return msg
}

func (app App) UpdateSubModels(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	if app.menu, cmd = app.menu.Update(msg); cmd != nil {
		cmds = append(cmds, cmd)
	}
	if app.podcasts, cmd = app.podcasts.Update(msg); cmd != nil {
		cmds = append(cmds, cmd)
	}
	if app.episodes, cmd = app.episodes.Update(msg); cmd != nil {
		cmds = append(cmds, cmd)
	}
	return app, tea.Batch(cmds...)
}

func (app App) UpdateFeeds(data *db.SubData) (tea.Model, tea.Cmd) {
	app.data = data
	feeds := make([]list.Item, len(app.data.Feeds))
	episodes := make([]list.Item, len(app.data.Episodes))
	for i, feed := range app.data.Feeds {
		feeds[i] = list.Item(feed)
	}
	for i, episode := range app.data.Episodes {
		episodes[i] = list.Item(episode)
	}
	app.podcasts, _ = app.podcasts.Update(feeds)
	app.episodes, _ = app.episodes.Update(episodes)
	return app.UpdateSubModels(data)
}
