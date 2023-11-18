package models

import (
	textinput "github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/mikemehl/dripper/utils"
)

type AddContext func(input string) tea.Cmd

type Add struct {
	input       textinput.Model
	promptStyle lipgloss.Style
	boxStyle    lipgloss.Style
	prompt      string
	context     tea.Msg
	width       int
	height      int
}

type FocusAdd struct {
	Prompt  string
	Context tea.Msg
}

type FocusRemove struct {
	Value    string
	Complete bool
	Context  tea.Msg
}

var (
	defaultPromptStyle    = lipgloss.NewStyle().Align(lipgloss.Left)
	defaultPromptBoxStyle = lipgloss.NewStyle().Border(lipgloss.DoubleBorder()).Padding(1).Align(lipgloss.Center)
)

func NewAdd() tea.Model {
	input := textinput.New()
	input.Prompt = ""
	input.Blur()
	input.Cursor.Blink = true
	return Add{
		input:       textinput.New(),
		promptStyle: defaultPromptStyle,
		boxStyle:    defaultPromptBoxStyle,
		context:     nil,
		prompt:      "",
	}
}

func (a Add) Init() tea.Cmd {
	return nil
}

func (a Add) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd = nil
	switch msg := msg.(type) {
	case FocusAdd:
		log.Debug("FocusAdd message received", "prompt", msg.Prompt)
		if !a.input.Focused() {
			log.Debug("Focused")
			a.context = msg.Context
			a.prompt = msg.Prompt
			a.input.Focus()
		}
	case tea.KeyMsg:
		if a.input.Focused() {
			switch msg.String() {
			case "esc":
				return a.PromptClosed(false)
			case "enter":
				return a.PromptClosed(true)
			}
		}
	case tea.WindowSizeMsg:
		_ = a.SetDimensions(msg)
	}
	a.input, cmd = a.input.Update(msg)
	return a, cmd
}

func (a *Add) PromptClosed(finished bool) (tea.Model, tea.Cmd) {
	context := a.context
	a.context = nil
	a.input.Blur()
	a.input.Reset()
	return a, func() tea.Msg { return FocusRemove{Value: a.input.Value(), Complete: finished, Context: context} }
}

func (a Add) View() string {
	view := a.promptStyle.Align(lipgloss.Left).Render(a.prompt)
	view = lipgloss.JoinVertical(lipgloss.Top, view, a.input.View())
	view = a.boxStyle.Render(view)
	view = lipgloss.Place(a.width, a.height, lipgloss.Center, lipgloss.Center, view, lipgloss.WithWhitespaceChars(""), lipgloss.WithWhitespaceBackground(lipgloss.Color("#666666")))
	return view
}

func (a *Add) SetDimensions(msg tea.WindowSizeMsg) tea.WindowSizeMsg {
	msg = utils.ScaleDimensions(msg, 10, 8, 8)
	a.width = msg.Width
	a.height = msg.Height
	msg = utils.ScaleDimensions(msg, 10, 3, 3)
	a.boxStyle = a.boxStyle.Width(msg.Width).Height(msg.Height).Align(lipgloss.Center, lipgloss.Center)
	return msg
}

func (a Add) Focused() bool {
	return a.input.Focused()
}
