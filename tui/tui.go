package tui

import (
	"bufio"
	"os"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	log "github.com/charmbracelet/log"
	"github.com/mikemehl/dripper/db"
	models "github.com/mikemehl/dripper/models"
	"github.com/mikemehl/dripper/utils"
)

var (
	appBoxStyle = lipgloss.NewStyle().Border(lipgloss.DoubleBorder()).Padding(1)
	logoStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF0080")).Width(40).Align(lipgloss.Right)
	logoText    = `
 █▀▄ █▀█ █ █▀█ █▀█ █▀▀ █▀█
 █▄▀ █▀▄ █ █▀▀ █▀▀ ██▄ █▀▄
  `
	logo = logoStyle.Render(logoText)
)

type App struct {
	podcasts tea.Model
	episodes tea.Model
	add      tea.Model
	menu     tea.Model
	data     *db.SubData
	style    lipgloss.Style
}

func Run() error {
	log.Debug("Run() called")
	app := NewApp()
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
		log.Debug("SelectPodcast", "selected", selected)
		cmds = append(cmds, func() tea.Msg {
			episodes := make([]*db.Episode, len(selected.Items))
			for i, episode := range selected.Items {
				ep := db.Episode(*episode)
				episodes[i] = &ep
			}
			return episodes
		})
		cmds = append(cmds, func() tea.Msg { return models.MenuSetActive{Index: 1} })
	}
	return tea.Batch(cmds...)
}

func NewApp() tea.Model {
	podcasts := models.NewDetailList([]list.Item{}, 40, 40, SelectPodcast, models.ActionMap{})
	episodes := models.NewDetailList([]list.Item{}, 40, 40, models.DefaultSelectAction, models.ActionMap{"d": models.DownloadEpisodeAction})
	add := models.NewAdd()
	return App{
		menu: models.NewMenu([]models.MenuItem{
			{Name: "Podcasts", Action: nil},
			{Name: "Episodes", Action: nil},
		}),
		data:     nil,
		podcasts: podcasts,
		episodes: episodes,
		style:    appBoxStyle,
		add:      add,
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
		newApp, batchCmd := app.UpdateSubModels(msg)
		return newApp, tea.Batch(cmd, batchCmd)
	case tea.KeyMsg:
		cmd = app.processKey(msg)
	case *db.SubData:
		return app.UpdateFeeds(msg)
	case []*db.Episode:
		app.episodes, _ = app.UpdateEpisodes(msg)
		return app, cmd
	case models.DownloadEpisode:
		log.Debug("DownloadEpisode", "msg", msg)
		cmd = downloadEpisode(msg)
	case models.FocusRemove:
		if msg.Complete {
			context := msg.Context
			value := msg.Value
			cmd = context.(models.AddContext)(value)
		}
	default:
	}
	newApp, batchCmd := app.UpdateSubModels(msg)
	return newApp, tea.Batch(cmd, batchCmd)
}

func (app App) View() string {
	var currView string
	active := app.menu.(models.Menu).Active()
	switch {
	case app.AddFocused():
		log.Debug("Add View")
		currView = app.add.View()
	case active == 0:
		log.Debug("Podcasts View")
		currView = app.podcasts.View()
	case active == 1:
		log.Debug("Episodes View")
		currView = app.episodes.View()
	default:
		currView = "Error"
	}
	return app.style.Render(lipgloss.JoinVertical(lipgloss.Top, app.menu.View(), currView))
}

func (app *App) processKey(msg tea.KeyMsg) tea.Cmd {
	var cmd tea.Cmd = nil
	key := msg.String()
	switch {
	case app.AddFocused() && key != "ctrl+c":
		cmd = nil
	case key == "q" || key == "ctrl+c":
		cmd = tea.Quit
	case key == "u":
		cmd = CmdWithSpinnerMessage("Updating Feeds", db.UpdateFeeds)
	case key == "d":
	case key == "esc":
		// Reset episodes and return focus to podcasts
		cmd = tea.Batch(func() tea.Msg { return app.data }, func() tea.Msg { return models.MenuSetActive{Index: 0} })
	case key == "a":
		cmd = func() tea.Msg {
			return models.FocusAdd{
				Prompt: "Add Podcast URL", Context: models.AddContext(
					func(s string) tea.Cmd {
						log.Debug("Adding feed", "url", s)
						_, err := db.NewFeed(s)
						if err != nil {
							return nil
						}
						return db.LoadFeeds
					}),
			}
		}
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
	if app.add, cmd = app.add.Update(msg); cmd != nil {
		cmds = append(cmds, cmd)
	}
	switch add := app.add.(type) {
	case models.Add:
		if add.Focused() {
			switch msg.(type) {
			case tea.KeyMsg:
				msg = nil
			}
		}
	}
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

func (app App) AddFocused() bool {
	switch add := app.add.(type) {
	case models.Add:
		return add.Focused()
	}
	return false
}

func CmdWithSpinnerMessage(msg string, cmd tea.Cmd) tea.Cmd {
	return tea.Sequence(
		func() tea.Msg { return models.SpinnerCmd{Active: true} },
		func() tea.Msg { return models.MessageCmd{Msg: msg, Active: true} },
		cmd,
		func() tea.Msg { return models.SpinnerCmd{Active: false} },
	)
}

func downloadEpisode(episode models.DownloadEpisode) tea.Cmd {
	return tea.Batch(
		func() tea.Msg { return models.SpinnerCmd{Active: true} },
		func() tea.Msg { return models.MessageCmd{Msg: "Downloading Episode to download.mp3", Active: true} },
		func() tea.Msg {
			audio, err := utils.GetAudio(episode.Url)
			if err != nil {
				log.Error("Error downloading audio", "error", err)
				return nil
			}

			log.Debug("Creating file")
			file, err := os.Create(episode.Filename)
			if err != nil {
				log.Error("Create file failed", "error", err)
				return nil
			}
			defer file.Close()

			writer := bufio.NewWriter(file)
			if len, err := writer.ReadFrom(audio); len <= 0 || err != nil {
				log.Error("Read from body failed", "error", err)
				return nil
			}

			if err := writer.Flush(); err != nil {
				log.Error("Flush failed", "error", err)
				return nil
			}

			return nil
		},
		func() tea.Msg { return models.SpinnerCmd{Active: false} },
	)
}
