package app

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/yuasalily/harbor-tui/internal/app/cmds"
	"github.com/yuasalily/harbor-tui/internal/app/ports"
)

type Model struct {
	witdh, height int

	// Docker
	dockerOK       bool
	dockerErr      string
	serverVersion  string
	daemonPlatform string

	images []ports.ImageSummary

	api ports.DockerAPI
}

func New(api ports.DockerAPI) Model {
	return Model{api: api}
}

// Bubble Tea ライフサイクル
func (m Model) Init() tea.Cmd { return cmds.FetchDaemonInfoCmd(m.api, 3*time.Second) }

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "i":
			return m, cmds.FetchImagesCmd(m.api, ports.ImagesListOptions{All: true}, 5*time.Second)
		}
	case tea.WindowSizeMsg:
		m.witdh, m.height = msg.Width, msg.Height
	case cmds.DaemonInfoMsg:
		if msg.Err != nil {
			m.dockerErr = msg.Err.Error()
			m.dockerOK = false
			return m, nil
		}
		m.dockerOK = true
		m.serverVersion = msg.Info.Version
		m.daemonPlatform = msg.Info.OS
	case cmds.ImagesListedMsg:
		if msg.Err != nil {
			m.dockerErr = msg.Err.Error()
			return m, nil
		}
		m.dockerErr = ""
		m.images = msg.Items
	}
	return m, nil
}

func (m Model) View() string {
	lines := []string{
		"",
		"  Harbor-TUI: Bubble Tea + Docker SDK",
		"",
	}

	status := "NOT CONNECTED"
	info := ""
	if m.dockerOK {
		status = "CONNECTED"
		info = fmt.Sprintf("Docker %s (%s)", m.serverVersion, m.daemonPlatform)
	} else if m.dockerErr != "" {
		info = fmt.Sprintf("error: %s", m.dockerErr)
	}
	lines = append(lines,
		fmt.Sprintf("  Status: %s", status),
		fmt.Sprintf("  %s", info),
		"",
		"  [Keys] q: quit   i:fetch images",
		"",
	)
	lines = append(lines, m.renderImagesSection()...)

	return strings.Join(lines, "\n")
}

func (m Model) renderImagesSection() []string {
	out := []string{"  Images:"}
	count := len(m.images)
	out = append(out, fmt.Sprintf("  - total: %d", count))
	if count == 0 {
		out = append(out, "  (press 'i' to load images)")
		return out
	}

	const maxRows = 10
	rows := count
	if rows > maxRows {
		rows = maxRows
	}
	out = append(out, "")
	out = append(out, "  ID (short)       Tags                           Size     Created")
	out = append(out, "  ----------------  -----------------------------  -------  ------------")
	for i := 0; i < rows; i++ {
		it := m.images[i]
		id := shortID(it.ID)
		tags := "<none>"
		if len(it.RepoTags) > 0 {
			tags = it.RepoTags[0]
		}
		size := humanBytes(it.Size)
		created := it.CreatedAt.Local().Format("2006-01-02 15:04")
		out = append(out, fmt.Sprintf("  %-17s  %-30s  %7s  %s", id, tags, size, created))
	}
	if count > rows {
		out = append(out, fmt.Sprintf("  ... and %d more", count-rows))
	}
	return out
}

func shortID(id string) string {
	if len(id) > 12 {
		return id[:12]
	}
	return id
}

func humanBytes(n int64) string {
	const unit = 1024
	if n < unit {
		return fmt.Sprintf("%dB", n)
	}
	div, exp := int64(unit), 0
	for n >= unit && exp < 4 {
		n /= unit
		div *= unit
		exp++
	}
	suffix := []string{"KB", "MB", "GB", "TB"}[exp-1]
	return fmt.Sprintf("%d%s", n, suffix)
}
