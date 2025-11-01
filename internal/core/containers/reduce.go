package containers

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/yuasalily/harbor-tui/internal/app/cmds"
)

func (s *State) Reduce(msg tea.Msg) bool {
	switch x := msg.(type) {
	case cmds.ContainersListMsg:
		if x.Err != nil {
			s.Err = x.Err.Error()
			return true
		}
		s.Err = ""
		s.List = x.Items
		return true
	}
	return false
}
