package main

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mmcdole/gofeed"
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

type subData struct {
	feeds    []gofeed.Feed
	episodes []*gofeed.Item
}

type modelState int

const (
	intro modelState = iota
	subsPage
	episodespage
	searchpage
	playerPage
	configPage
)

type model struct {
	intro struct {
		spinner spinner.Model
		loader  chan subData
	}
	topMenu    menuBar
	modelState modelState
	sub_data   subData
}

func (m *model) initIntro() {
	s := spinner.New()
	s.Spinner = spinner.Globe
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Align(lipgloss.Center, lipgloss.Center)
	m.intro.spinner = s
	m.intro.loader = make(chan subData)
}

func initModel() tea.Model {
	var m model
	m.modelState = intro
	m.initIntro()
	return m
}

func (m model) Init() tea.Cmd {
	go loadSubData(m.intro.loader)
	return m.intro.spinner.Tick
}

func (m model) UpdateIntro(msg tea.Msg) (tea.Model, tea.Cmd) {
	select {
	case <-m.intro.loader:
		return m, tea.Quit
	default:
	}
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	}
	var cmd tea.Cmd
	m.intro.spinner, cmd = m.intro.spinner.Update(msg)
	return m, cmd
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m.modelState {
	case intro:
		return m.UpdateIntro(msg)
	}
	return m, nil
}

func (m model) ViewIntro() string {
	str := fmt.Sprintf("\n\n   %s Loading subscriptions... \n\n", m.intro.spinner.View())
	return str
}

func (m model) View() string {
	switch m.modelState {
	case intro:
		return m.ViewIntro()
	}
	return ""
}
