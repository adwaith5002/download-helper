package tui

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/adwaith5002/download-helper/internal/analyzer"
	"github.com/adwaith5002/download-helper/internal/organizer"
	"github.com/adwaith5002/download-helper/internal/scanner"
	"github.com/adwaith5002/download-helper/pkg/fileinfo"
	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
)

var progressBar = progress.New(progress.WithDefaultGradient())

type Model struct {
	files       []fileinfo.FileInfo
	cursor      int
	windowStart int
	screen      Screen
	plans       []organizer.Plan
	dupes       [][]fileinfo.FileInfo
	root        string
	statusMsg   string
	moveIndex   int
	moving      bool
	failCount   int
}
type Screen int

const (
	ScreenBrowse Screen = iota
	ScreenPlan
)

type fileMovedMsg struct {
	index int
	err   error
}

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
	case fileMovedMsg:
		if msg.err != nil {
			m.failCount++
		}
		m.moveIndex++
		if m.moveIndex < len(m.plans) {
			return m, moveFileCmd(m.plans[m.moveIndex], m.moveIndex)
		}
		m.moving = false
		if m.failCount > 0 {
			m.statusMsg = fmt.Sprintf("Moved %d files, %d failed", len(m.plans)-m.failCount, m.failCount)
		} else {
			m.statusMsg = fmt.Sprintf("Moved %d files successfully", len(m.plans))
		}
		m.screen = ScreenBrowse
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
		if len(m.plans) == 0 {
			m.screen = ScreenBrowse
			return m, nil
		}
		m.moving = true
		m.moveIndex = 0
		return m, moveFileCmd(m.plans[0], 0)
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
	s := titleStyle.Render(fmt.Sprintf("Download Helper — %d files", len(m.files))) + "\n"
	s += dimStyle.Render("↑/↓ navigate, o organize, q quit") + "\n\n"
	windowSize := 20
	windowEnd := m.windowStart + windowSize
	if windowEnd > len(m.files) {
		windowEnd = len(m.files)
	}

	for i := m.windowStart; i < windowEnd; i++ {
		f := m.files[i]
		namePart := fmt.Sprintf("%-40s", f.Name)
		catPart := fmt.Sprintf("%-12s", f.Category)
		sizePart := fmt.Sprintf("%d bytes", f.Size)

		if i == m.cursor {
			line := namePart + "  " + catPart + "  " + sizePart
			s += selectedStyle.Render("> "+line) + "\n"
		} else {
			styledCat := categoryColors[f.Category].Render(catPart)
			s += "  " + namePart + "  " + styledCat + "  " + sizePart + "\n"
		}
	}

	if m.statusMsg != "" {
		s += "\n" + m.statusMsg + "\n"
	}
	return s
}

func (m Model) viewPlan() string {
	if m.moving {
		return m.viewMoving()
	}
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

	s := titleStyle.Render(fmt.Sprintf("Proposed moves (%d total)", len(m.plans))) + "\n\n"
	for _, folder := range folderNames {
		plans := groups[folder]
		header := fmt.Sprintf("%s (%d files)", folder, len(plans))
		if folder == "Duplicates" {
			s += warningStyle.Render(header) + "\n"
		} else {
			s += titleStyle.Render(header) + "\n"
		}
		for _, p := range plans {
			s += "  " + filepath.Base(p.To) + "\n"
		}
		s += "\n"
	}
	s += dimStyle.Render("y: confirm   n/esc: cancel")
	return s
}

func moveFileCmd(p organizer.Plan, index int) tea.Cmd {
	return func() tea.Msg {
		err := os.MkdirAll(filepath.Dir(p.To), 0755)
		if err == nil {
			err = os.Rename(p.From, p.To)
		}
		return fileMovedMsg{index: index, err: err}
	}
}

func (m Model) viewMoving() string {
	total := len(m.plans)
	done := m.moveIndex
	percent := float64(done) / float64(total)

	s := titleStyle.Render("Moving files...") + "\n\n"
	s += progressBar.ViewAs(percent) + "\n"
	s += fmt.Sprintf("%d/%d files\n", done, total)
	return s
}
