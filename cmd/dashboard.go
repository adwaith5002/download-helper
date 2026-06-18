package cmd

import (
	"github.com/adwaith5002/download-helper/internal/scanner"
	"github.com/adwaith5002/download-helper/internal/tui"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

var dashboardCmd = &cobra.Command{
	Use:   "dashboard",
	Short: "Interactive dashboard for browsing scanned files",
	RunE: func(cmd *cobra.Command, args []string) error {
		path, err := cmd.Flags().GetString("path")
		if err != nil {
			return err
		}

		files, err := scanner.Scan(path)
		if err != nil {
			return err
		}

		model := tui.NewModel(files)
		p := tea.NewProgram(
			model,
			tea.WithAltScreen(),
		)
		_, err = p.Run()
		return err
	},
}

func init() {
	dashboardCmd.Flags().String("path", "", "Path to scan (required)")
	dashboardCmd.MarkFlagRequired("path")
	rootCmd.AddCommand(dashboardCmd)
}
