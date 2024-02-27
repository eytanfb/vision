package main

import (
	"os"
	"vision/app"
	"vision/config"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
)

func main() {
	cfg, err := config.LoadConfig("/Users/eytananjel/Code/vision/config/config.json")
	if err != nil {
		panic(err) // Simplified error handling for brevity
	}

	args := os.Args[1:]
	initialModel := app.InitialModel(cfg, args) // Pass cmdline args to the model

	p := tea.NewProgram(initialModel, tea.WithAltScreen(), tea.WithMouseCellMotion())
	if _, err := p.Run(); err != nil {
		log.Error("Oh no!", err)
		os.Exit(1)
	}
}
