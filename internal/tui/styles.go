package tui

import(
	"github.com/adwaith5002/download-helper/pkg/fileinfo"
	"github.com/charmbracelet/lipgloss"
)
var (
    titleStyle = lipgloss.NewStyle().
        Bold(true).
        Foreground(lipgloss.Color("#7D56F4"))

    selectedStyle = lipgloss.NewStyle().
        Background(lipgloss.Color("#7D56F4")).
        Foreground(lipgloss.Color("#FFFFFF"))

    dimStyle = lipgloss.NewStyle().
        Foreground(lipgloss.Color("#888888"))

    warningStyle = lipgloss.NewStyle().
        Foreground(lipgloss.Color("#FFA500"))
)
var categoryColors = map[fileinfo.Category]lipgloss.Style{
    fileinfo.Image:      lipgloss.NewStyle().Foreground(lipgloss.Color("#FF79C6")),
    fileinfo.Document:   lipgloss.NewStyle().Foreground(lipgloss.Color("#8BE9FD")),
    fileinfo.Video:      lipgloss.NewStyle().Foreground(lipgloss.Color("#FFB86C")),
    fileinfo.Audio:      lipgloss.NewStyle().Foreground(lipgloss.Color("#50FA7B")),
    fileinfo.Archive:    lipgloss.NewStyle().Foreground(lipgloss.Color("#F1FA8C")),
    fileinfo.Code:       lipgloss.NewStyle().Foreground(lipgloss.Color("#BD93F9")),
    fileinfo.Executable: lipgloss.NewStyle().Foreground(lipgloss.Color("#FF5555")),
}