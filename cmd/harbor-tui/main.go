package main

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/yuasalily/harbor-tui/internal/adapters/dockeradapter"
	"github.com/yuasalily/harbor-tui/internal/app"
	"github.com/yuasalily/harbor-tui/internal/docker"
)

func main() {
	cli, err := docker.New()
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()
	api := dockeradapter.New(cli)
	m := app.New(api)
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
