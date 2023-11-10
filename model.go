package main

import (
	"fmt"

	"github.com/charmbracelet/log"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
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
	loader     chan subData
	addSubChan chan error
	intro      struct {
		spinner spinner.Model
	}
	topMenu    menuBar
	addSubForm textinput.Model
	subList    table.Model
	modelState modelState
	sub_data   subData
}

func (m *model) initIntro() {
	s := spinner.New()
	s.Spinner = spinner.Globe
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Align(lipgloss.Center, lipgloss.Center)
	m.intro.spinner = s
	m.loader = make(chan subData)
}

func initModel() tea.Model {
	var m model
	m.modelState = intro
	m.topMenu = initMenubarModel()
	m.addSubForm = textinput.New()
	m.subList = table.New(table.WithColumns([]table.Column{{Width: 120}}), table.WithStyles(table.DefaultStyles()), table.WithFocused(false))
	m.initIntro()
	return m
}

func (m model) Init() tea.Cmd {
	go loadSubData(m.loader)
	return m.intro.spinner.Tick
}

func (m model) UpdateIntro(msg tea.Msg) (tea.Model, tea.Cmd) {
	select {
	case subs := <-m.loader:
		log.Info("Got subs", "subs", subs)
		m.modelState = subsPage
		m.sub_data = subs
		return m.UpdateSubs(msg)
	default:
	}
	var cmd tea.Cmd
	m.intro.spinner, cmd = m.intro.spinner.Update(msg)
	return m, cmd
}

func (m model) UpdateSubs(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keyBindings.menuNext):
			m.topMenu.selectNext()
		case key.Matches(msg, keyBindings.menuPrev):
			m.topMenu.selectPrevious()
		case key.Matches(msg, keyBindings.addSub):
			m.modelState = addingSub
			m.addSubForm.Focus()
			return m, textinput.Blink
		case key.Matches(msg, keyBindings.down):
			m.subList.MoveDown(1)
		case key.Matches(msg, keyBindings.up):
			m.subList.MoveUp(1)
		}
	}
	var titles []table.Row
	for _, feed := range m.sub_data.feeds {
		titles = append(titles, table.Row{feed.Title})
	}
	m.subList.SetRows(titles)
	m.subList, cmd = m.subList.Update(msg)
	return m, cmd
}

func (m model) UpdateAddingSub(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keyBindings.quit):
			m.modelState = subsPage
			m.addSubForm.Blur()
			return m, nil
		case key.Matches(msg, keyBindings.confirm):
			m.modelState = subsPage
			m.addSubForm.Blur()
			go newFeed(m.addSubChan, m.addSubForm.Value())
			return m, nil
		}
	}
	var cmd tea.Cmd
	m.addSubForm, cmd = m.addSubForm.Update(msg)
	return m, cmd
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if key.Matches(msg, keyBindings.quit) && m.modelState != addingSub {
			return m, tea.Quit
		}
	}
	select {
	case status := <-m.addSubChan:
		if status == nil {
			log.Fatal(status)
		}
		return m, nil
	default:
	}
	switch m.modelState {
	case intro:
		return m.UpdateIntro(msg)
	case subsPage:
		return m.UpdateSubs(msg)
	case addingSub:
		return m.UpdateAddingSub(msg)
	default:

	}
	return m, nil
}

func (m model) ViewIntro() string {
	str := fmt.Sprintf("\n\n   %s Loading subscriptions... \n\n", m.intro.spinner.View())
	return str
}

func (m model) ViewAddingSub() string {
	return fmt.Sprintf("Enter feed URL:\n\n%s", m.addSubForm.View())
}

func (m model) ViewSubs() string {
	log.Info("ViewSubs")
	s := fmt.Sprintf("%s%s\n", m.topMenu.View(), m.subList.View())
	log.Info("ViewSubs", "s", s)
	return s
}

func (m model) View() string {
	switch m.modelState {
	case intro:
		return m.ViewIntro()
	case subsPage:
		return m.ViewSubs()
	case addingSub:
		return m.ViewAddingSub()
	default:
		return m.topMenu.View()
	}
	return ""
}
