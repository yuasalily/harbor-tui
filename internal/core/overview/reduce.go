package overview

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/yuasalily/harbor-tui/internal/app/cmds"
)

func (s *State) Reduce(msg tea.Msg) bool {
	switch x := msg.(type) {
	case cmds.DaemonInfoMsg:
		if x.Err != nil {
			s.Conn.OK = false
			s.Conn.Err = x.Err.Error()
			return true
		}
		s.Conn.OK = true
		s.Conn.Err = ""
		s.Daemon = x.Info
		return true
	}
	return false
}
