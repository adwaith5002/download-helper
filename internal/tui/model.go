package tui

import (
	"fmt"
	"path/filepath"
	"sort"

	"github.com/adwaith5002/download-helper/internal/analyzer"
	"github.com/adwaith5002/download-helper/internal/organizer"
	"github.com/adwaith5002/download-helper/internal/scanner"
	"github.com/adwaith5002/download-helper/pkg/fileinfo"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	files       []fileinfo.FileInfo
	cursor      int
	windowStart int
	screen      Screen
	plans       []organizer.Plan
	dupes       [][]fileinfo.FileInfo
	root        string
	statusMsg   string
}
type Screen int

const (
	ScreenBrowse Screen = iota
	ScreenPlan
)

func NewModel(files []fileinfo.FileInfo, dupes [][]fileinfo.FileInfo, root string) Model {
	return Model{
		files:  files,
		dupes:  dupes,
		root:   root,
		screen: ScreenBrowse,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch m.screen {
		case ScreenBrowse:
			return m.updateBrowse(msg)
		case ScreenPlan:
			return m.updatePlan(msg)
		}
	}
	return m, nil
}
func (m Model) updateBrowse(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
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
	case "o":
		freshFiles, err := scanner.Scan(m.root)
		if err != nil {
			m.statusMsg = fmt.Sprintf("Scan error: %v", err)
			return m, nil
		}
		freshDupes, err := analyzer.FindDuplicates(freshFiles)
		if err != nil {
			m.statusMsg = fmt.Sprintf("Duplicate scan error: %v", err)
			return m, nil
		}
		m.files = freshFiles
		m.dupes = freshDupes
		m.plans = organizer.BuildPlan(m.files, m.dupes, m.root)
		m.screen = ScreenPlan
		return m, nil
	}
	return m, nil
}
func (m Model) updatePlan(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "y":
		err := organizer.Execute(m.plans)
		if err != nil {
			m.statusMsg = fmt.Sprintf("Error: %v", err)
		} else {
			m.statusMsg = fmt.Sprintf("Moved %d files successfully", len(m.plans))
		}
		m.screen = ScreenBrowse
	case "n", "esc":
		m.screen = ScreenBrowse
	case "up", "down", "j", "k":

		// no-op for now, just absorb the keypress
	}

	return m, nil
}

func (m Model) View() string {
	switch m.screen {
	case ScreenPlan:
		return m.viewPlan()
	default:
		return m.viewBrowse()
	}
}

func (m Model) viewBrowse() string {
	s := fmt.Sprintf("Download Helper — %d files | ↑/↓ navigate, o organize, q quit\n\n", len(m.files))

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

	if m.statusMsg != "" {
		s += "\n" + m.statusMsg + "\n"
	}
	return s
}

func (m Model) viewPlan() string {
	groups := make(map[string][]organizer.Plan)
	for _, p := range m.plans {
		folder := filepath.Base(filepath.Dir(p.To))
		groups[folder] = append(groups[folder], p)
	}

	var folderNames []string
	for folder := range groups {
		folderNames = append(folderNames, folder)
	}
	sort.Strings(folderNames)

	s := fmt.Sprintf("Proposed moves (%d total):\n\n", len(m.plans))
	for _, folder := range folderNames {
		plans := groups[folder]
		s += fmt.Sprintf("%s (%d files)\n", folder, len(plans))
		for _, p := range plans {
			s += fmt.Sprintf("  %s\n", filepath.Base(p.To))
		}
		s += "\n"
	}
	s += "y: confirm   n/esc: cancel"
	return s
}
