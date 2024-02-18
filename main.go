package main

import (
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
)

func main() {
	cfg, err := LoadConfig("config.json")
	if err != nil {
		panic(err) // Simplified error handling for brevity
	}

	args := os.Args[1:]
	initialModel := InitialModel(cfg, args) // Pass cmdline args to the model

	p := tea.NewProgram(initialModel)
	if _, err := p.Run(); err != nil {
		log.Error("Oh no!", err)
		os.Exit(1)
	}
}
