package cmd

import (
	"fmt"
	"os"

	"github.com/adwaith5002/download-helper/internal/analyzer"
	"github.com/adwaith5002/download-helper/internal/organizer"
	"github.com/adwaith5002/download-helper/internal/scanner"
	"github.com/adwaith5002/download-helper/pkg/fileinfo"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "download-helper",
	Short: "A utility to manage your Downloads folder",
	RunE: func(cmd *cobra.Command, args []string) error {
		path, err := cmd.Flags().GetString("path")
		if err != nil {
			return err
		}
		files, err := scanner.Scan(path)
		if err != nil {
			return err
		}
		counts := make(map[fileinfo.Category]int)
		for _, f := range files {
			counts[f.Category]++
		}

		fmt.Printf("Found %d files\n\n", len(files))
		for category, count := range counts {
			fmt.Printf("  %-12s %d files\n", category, count)
		}
		for _, f := range files {
			if f.Category == fileinfo.Unknown {
				fmt.Printf("  UNKNOWN: %s\n", f.Name)
			}
		}
		duplicates, err := analyzer.FindDuplicates(files)
		if err != nil {
			return err
		}

		fmt.Printf("\nDuplicates found: %d groups\n", len(duplicates))
		for i, group := range duplicates {
			fmt.Printf("\n  Group %d:\n", i+1)
			for _, f := range group {
				fmt.Printf("    %s (%d bytes)\n", f.Name, f.Size)
			}
		}
		recommendations := analyzer.Recommend(files, duplicates)
		fmt.Printf("\nRecommendations:\n")
		for _, r := range recommendations {
			fmt.Printf("  [%s] %s\n", r.Priority, r.Message)
		}
		plans := organizer.BuildPlan(files, duplicates, path)
		fmt.Printf("\nOrganization plan (%d moves):\n", len(plans))
		for _, p := range plans {
			if p.IsDupe {
				fmt.Printf("  [DUPE] %s → %s\n", p.From, p.To)
			} else {
				fmt.Printf("  [MOVE] %s → %s\n", p.From, p.To)
			}
		}

		confirm, err := cmd.Flags().GetBool("confirm")
		if err != nil {
			return err
		}
		if confirm {
			err = organizer.Execute(plans)
			if err != nil {
				return err
			}
			fmt.Println("\nDone. Files organized.")
		} else {
			fmt.Println("\nDry run. Pass --confirm to execute.")
		}
		return nil
	},
	
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().String("path", "", "Path to scan (required)")
	rootCmd.MarkFlagRequired("path")
	rootCmd.Flags().Bool("confirm", false, "Execute the organization plan")

}
