package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
)

// reportCmd represents the report command
var reportCmd = &cobra.Command{
	Use:   "report",
	Short: "Generate a report of the last week",
	Run: func(cmd *cobra.Command, args []string) {
		var history []TimeEntry
		historyData, err := os.ReadFile("history.json")
		if err != nil {
			fmt.Println("No history found.")
			return
		}
		json.Unmarshal(historyData, &history)

		oneWeekAgo := time.Now().AddDate(0, 0, -7)
		report := make(map[string]time.Duration)

		for _, entry := range history {
			if entry.EndTime.After(oneWeekAgo) {
				duration := entry.EndTime.Sub(entry.StartTime)
				report[entry.Project] += duration
			}
		}

		for project, duration := range report {
			fmt.Printf("Project: %s, Time Spent: %v\n", project, duration)
		}
	},
}

func init() {
	rootCmd.AddCommand(reportCmd)
}
