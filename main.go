package main

import (
	"os"
	"vision/app"
	"vision/config"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
)

func main() {
	homeDirectory, _ := os.UserHomeDir()
	configPath := homeDirectory + "/Code/vision/config/config.json"

	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		panic(err) // Simplified error handling for brevity
	}

	args := os.Args[1:]
	initialModel := app.InitialModel(cfg, args) // Pass cmdline args to the model

	p := tea.NewProgram(initialModel, tea.WithMouseCellMotion(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Error("Oh no!", err)
		os.Exit(1)
	}
}
