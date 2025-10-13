package images

import "github.com/yuasalily/harbor-tui/internal/app/ports"

type State struct {
	List   []ports.ImageSummary
	Err    string
	Filter ports.ImagesListOptions
}

func New() *State { return &State{Filter: ports.ImagesListOptions{All: true}} }
