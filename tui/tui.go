package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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
	program := tea.NewProgram(NewApp())
	if _, err := program.Run(); err != nil {
		return err
	}
	return nil
}

func NewApp() App {
	return App{
		menu: models.NewMenu([]models.MenuItem{
			{Name: "Podcasts", Action: nil},
			{Name: "Episodes", Action: nil},
		}),
	}
}

func (app App) Init() tea.Cmd {
	return func() tea.Msg { return db.LoadFeeds() }
}

func (app App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if cmd = app.processKey(msg); cmd != nil {
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
