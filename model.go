package main

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mmcdole/gofeed"
)

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
	addingSub
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
	m.topMenu = initMenubarModel()
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
		m.modelState = subsPage
		return m, nil
	default:
	}
	var cmd tea.Cmd
	m.intro.spinner, cmd = m.intro.spinner.Update(msg)
	return m, cmd
}

func (m model) UpdateSubs(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keyBindings.menuNext):
			m.topMenu.selectNext()
		case key.Matches(msg, keyBindings.menuPrev):
			m.topMenu.selectPrevious()
		case key.Matches(msg, keyBindings.addSub):
			m.modelState = addingSub
		}
	}
	return m, nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if key.Matches(msg, keyBindings.quit) {
			return m, tea.Quit
		}
	}
	switch m.modelState {
	case intro:
		return m.UpdateIntro(msg)
	case subsPage:
		return m.UpdateSubs(msg)
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
	case subsPage:
		return m.topMenu.View()
	}
	return ""
}
