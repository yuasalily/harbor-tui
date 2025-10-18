package containers

import "github.com/yuasalily/harbor-tui/internal/app/ports"

type State struct {
	List   []ports.ContainerSummary
	Err    string
	Filter ports.ContainersListOptions
}

func New() *State { return &State{Filter: ports.ContainersListOptions{All: true}} }
