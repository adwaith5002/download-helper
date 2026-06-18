package tui

import (
	"fmt"

	"github.com/adwaith5002/download-helper/pkg/fileinfo"
	tea "github.com/charmbracelet/bubbletea"
)
type Model struct {
    files       []fileinfo.FileInfo
    cursor      int
    windowStart int  // index of first visible file
}

func NewModel(files []fileinfo.FileInfo) Model {
	return Model{files: files, cursor: 0}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
				if m.cursor < m.windowStart {
					m.windowStart--
				}
			}
		case "down", "j":
			if m.cursor < len(m.files)-1 {
				m.cursor++
				if m.cursor >= m.windowStart+20 {
					m.windowStart++
				}
			}
		}
	}
	return m, nil
}
func (m Model) View() string {
    s := fmt.Sprintf("Download Helper — %d files | ↑/↓ to navigate, q to quit\n\n", len(m.files))

    windowSize := 20
    windowEnd := m.windowStart + windowSize
    if windowEnd > len(m.files) {
        windowEnd = len(m.files)
    }

    for i := m.windowStart; i < windowEnd; i++ {
        f := m.files[i]
        cursor := "  "
        if i == m.cursor {
            cursor = "> "
        }
        s += fmt.Sprintf("%s%-40s  %-12s  %d bytes\n", cursor, f.Name, f.Category, f.Size)
    }

    return s
}
