package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	log "github.com/charmbracelet/log"
	tui "github.com/mikemehl/dripper/tui"
)

type errMsg error

func main() {
	if err := setupLogging(); err != nil {
		os.Exit(1)
	}
	IntroLog()
	log.Debug("Calling Run()")
	err := tui.Run()
	if err != nil {
		log.Error(err)
	} else {
		log.Info("Exiting Dripper")
	}
}

func setupLogging() error {
	f, err := tea.LogToFile("dripper.log", "debug")
	if err != nil {
		fmt.Println(err)
		return err
	}
	log.SetOutput(f)
	return nil
}

func IntroLog() {
	hr := lipgloss.NewStyle().Border(lipgloss.DoubleBorder(), false, false, true).Width(20)
	log.Info(hr.Render())
	log.Info("Starting Dripper")
}
