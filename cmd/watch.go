package cmd

import (
	"github.com/adwaith5002/download-helper/internal/watcher"
	"github.com/spf13/cobra"
)

var watchCmd = &cobra.Command{
	Use:   "watch",
	Short: "Watch the Downloads folder for new files",
	RunE: func(cmd *cobra.Command, args []string) error {
		path, err := cmd.Flags().GetString("path")
		if err != nil {
			return err
		}
		return watcher.Watch(path)
	},
}

func init() {
	watchCmd.Flags().String("path", "", "Path to watch (required)")
	watchCmd.MarkFlagRequired("path")
	rootCmd.AddCommand(watchCmd)
}
