package overview

import "github.com/yuasalily/harbor-tui/internal/app/ports"

type State struct {
	Conn struct {
		OK  bool
		Err string
	}
	Daemon ports.DaemonInfo
}

func New() *State { return &State{} }
