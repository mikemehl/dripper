package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
	pages "github.com/mikemehl/dripper/pages"
	tui "github.com/mikemehl/dripper/tui"
)

type errMsg error

func main() {
	err := tui.Run()
	if err != nil {
		fmt.Println(err)
	}
}

func old_main() {
	f, err := os.OpenFile("debug.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening log file")
	}
	log.SetLevel(log.DebugLevel)
	log.SetOutput(f)
	log.SetReportCaller(true)
	log.Info("Starting app========================")
	p := tea.NewProgram(pages.InitMainModel())
	if _, err := p.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
