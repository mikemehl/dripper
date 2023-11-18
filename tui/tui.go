package tui

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	log "github.com/charmbracelet/log"
	"github.com/mikemehl/dripper/db"
	models "github.com/mikemehl/dripper/models"
	"github.com/mikemehl/dripper/utils"
)

var appBoxStyle = lipgloss.NewStyle().Border(lipgloss.DoubleBorder()).Padding(1)

type App struct {
	style    lipgloss.Style
	data     *db.SubData
	podcasts tea.Model
	episodes tea.Model
	menu     tea.Model
}

func Run() error {
	log.Debug("Run() called")
	var app tea.Model = NewApp()
	program := tea.NewProgram(app)
	if _, err := program.Run(); err != nil {
		log.Error("Dripper quit with error", "error", err)
		return err
	}
	return nil
}

func SelectPodcast(d models.DetailList) tea.Cmd {
	var cmds []tea.Cmd
	selected := d.SelectedItem()
	switch selected := selected.(type) {
	case db.Feed:
		cmds = append(cmds, func() tea.Msg { return selected })
		cmds = append(cmds, func() tea.Msg { return models.MenuSetActive{Index: 1} })
	}
	return tea.Batch(cmds...)
}

func NewApp() tea.Model {
	podcasts := models.NewDetailList([]list.Item{}, 40, 40, SelectPodcast)
	episodes := models.NewDetailList([]list.Item{}, 40, 40, models.DefaultSelectAction)
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
	var cmd tea.Cmd = nil
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		msg = app.SetDimensions(msg)
	case tea.KeyMsg:
		if cmd = app.processKey(msg); cmd != nil {
			app, batchCmd := app.UpdateSubModels(msg)
			return app, tea.Batch(cmd, batchCmd)
		}
	case *db.SubData:
		return app.UpdateFeeds(msg)
	case []*db.Episode:
		app.episodes, _ = app.UpdateEpisodes(msg)
		return app, cmd
	default:
	}
	return app.UpdateSubModels(msg)
}

func (app App) View() string {
	var currView string
	switch menu := app.menu.(type) {
	case models.Menu:
		switch menu.Active() {
		case 0:
			log.Debug("Podcasts View")
			currView = app.podcasts.View()
		case 1:
			log.Debug("Episodes View")
			currView = app.episodes.View()
		default:
			currView = "Error"
		}
	}
	return app.style.Render(lipgloss.JoinVertical(lipgloss.Top, app.menu.View(), currView))
}

func (app *App) processKey(msg tea.KeyMsg) tea.Cmd {
	var cmd tea.Cmd = nil
	switch msg.String() {
	case "q", "ctrl+c":
		cmd = tea.Quit
	case "esc":
		cmd = func() tea.Msg { return app.data }
	}
	return cmd
}

func (app *App) SetDimensions(msg tea.WindowSizeMsg) tea.WindowSizeMsg {
	msg = utils.ScaleDimensions(msg, 10, 9, 9)
	app.style = app.style.Width(msg.Width).Height(msg.Height).Align(lipgloss.Left, lipgloss.Top)
	msg = utils.ScaleDimensions(msg, 10, 7, 5)
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
	app.podcasts, _ = app.UpdatePodcasts(data.Feeds)
	app.episodes, _ = app.UpdateEpisodes(data.Episodes)
	return app.UpdateSubModels(data)
}

func (app *App) UpdatePodcasts(data []db.Feed) (tea.Model, tea.Cmd) {
	feeds := make([]list.Item, len(data))
	for i, feed := range data {
		feeds[i] = list.Item(feed)
	}
	return app.podcasts.Update(feeds)
}

func (app *App) UpdateEpisodes(data []*db.Episode) (tea.Model, tea.Cmd) {
	episodes := make([]list.Item, len(data))
	for i, episode := range data {
		episodes[i] = list.Item(episode)
	}
	return app.episodes.Update(episodes)
}
