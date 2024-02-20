package main

import (
	"company-file-viewer/app"
	"company-file-viewer/config"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
)

func main() {
	cfg, err := config.LoadConfig("config/config.json")
	if err != nil {
		panic(err) // Simplified error handling for brevity
	}

	args := os.Args[1:]
	initialModel := app.InitialModel(cfg, args) // Pass cmdline args to the model

	p := tea.NewProgram(initialModel)
	if _, err := p.Run(); err != nil {
		log.Error("Oh no!", err)
		os.Exit(1)
	}
}
