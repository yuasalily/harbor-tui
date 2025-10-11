package main

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/yuasalily/harbor-tui/internal/adapters/dockeradapter"
	"github.com/yuasalily/harbor-tui/internal/app/ports"
	"github.com/yuasalily/harbor-tui/internal/core"
	"github.com/yuasalily/harbor-tui/internal/docker"
	"github.com/yuasalily/harbor-tui/internal/ui"
)

func main() {
	cli, err := docker.New()
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()
	var api ports.DockerAPI = dockeradapter.New(cli)
	appCore := core.New(api)

	tui := ui.New(&appCore)
	p := tea.NewProgram(tui, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
