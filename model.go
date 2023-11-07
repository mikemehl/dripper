package main

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var globalKeybindings = map[string]key.Binding{
	"l":         key.NewBinding(key.WithKeys("l"), key.WithHelp("l", "move right")),
	"h":         key.NewBinding(key.WithKeys("h"), key.WithHelp("h", "move left")),
	"j":         key.NewBinding(key.WithKeys("j"), key.WithHelp("j", "move down")),
	"k":         key.NewBinding(key.WithKeys("k"), key.WithHelp("k", "move up")),
	"Tab":       key.NewBinding(key.WithKeys("Tab"), key.WithHelp("Tab", "next page")),
	"Shift-Tab": key.NewBinding(key.WithKeys("Shift-Tab"), key.WithHelp("Shift-Tab", "previous page")),
	"Enter":     key.NewBinding(key.WithKeys("Enter"), key.WithHelp("Enter", "select")),
	"o":         key.NewBinding(key.WithKeys("o"), key.WithHelp("o", "open")),
	"p":         key.NewBinding(key.WithKeys("p"), key.WithHelp("p", "play/pause")),
	"/":         key.NewBinding(key.WithKeys("/"), key.WithHelp("/", "search")),
	"q":         key.NewBinding(key.WithKeys("q"), key.WithHelp("q", "quit")),
}

type modelPage interface{}

type introPage struct {
	spinner spinner.Model
}

type homePage struct {
	menu *menuBar
}

type subscriptionsPage struct {
	menu *menuBar
}

type episodesPage struct {
	menu *menuBar
}

type playerPage struct {
	menu *menuBar
}

type configPage struct {
	menu *menuBar
}

type model struct {
	intro         introPage
	home          homePage
	subscriptions subscriptionsPage
	episodes      episodesPage
	player        playerPage
	config        configPage
	menu          menuBar
	currentPage   *modelPage
}

func initModel() model {
	var m model
	m.intro = introPage{}
	m.home = homePage{}
	m.subscriptions = subscriptionsPage{}
	m.episodes = episodesPage{}
	m.player = playerPage{}
	m.config = configPage{}
	m.menu = initMenubarModel()
	m.currentPage = nil
	return m
}

func initIntro() tea.Model {
	s := spinner.New()
	s.Spinner = spinner.Globe
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Align(lipgloss.Center, lipgloss.Center)
	return introPage{spinner: s}
}

func (m introPage) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m introPage) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.spinner, cmd = m.spinner.Update(msg)
	return m, cmd
}

func (m introPage) View() string {
	str := fmt.Sprintf("\n\n   %s Loading forever... \n\n", m.spinner.View())
	return str
}
