package main

import (
	"fmt"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

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
	menu          *menuBar
	currentPage   *tea.Model
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
