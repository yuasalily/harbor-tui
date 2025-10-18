package core

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/yuasalily/harbor-tui/internal/app/ports"
	"github.com/yuasalily/harbor-tui/internal/core/containers"
	"github.com/yuasalily/harbor-tui/internal/core/images"
	"github.com/yuasalily/harbor-tui/internal/core/overview"
)

type Core struct {
	API        ports.DockerAPI
	Overview   *overview.State
	Images     *images.State
	Containers *containers.State
}

func New(api ports.DockerAPI) Core {
	return Core{
		API:        api,
		Overview:   overview.New(),
		Images:     images.New(),
		Containers: containers.New(),
	}
}

func (c *Core) Reduce(msg tea.Msg) {
	if c.Overview.Reduce(msg) {
		return
	}
	if c.Images.Reduce(msg) {
		return
	}
	if c.Containers.Reduce(msg) {
		return
	}
}
