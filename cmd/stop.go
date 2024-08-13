package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
)

// stopCmd represents the stop command
var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop tracking time",
	Run: func(cmd *cobra.Command, args []string) {
		data, err := os.ReadFile("current.json")
		if err != nil {
			fmt.Println("No ongoing tracking found.")
			return
		}

		var entry TimeEntry
		json.Unmarshal(data, &entry)
		entry.EndTime = time.Now()

		// Save the entry to a history file
		var history []TimeEntry
		historyData, _ := os.ReadFile("history.json")
		if len(historyData) > 0 {
			json.Unmarshal(historyData, &history)
		}
		history = append(history, entry)

		newHistoryData, _ := json.Marshal(history)
		os.WriteFile("history.json", newHistoryData, 0644)

		// Delete the current tracking file
		os.Remove("current.json")

		fmt.Printf("Stopped tracking time for project: %s\n", entry.Project)
	},
}

func init() {
	rootCmd.AddCommand(stopCmd)
}
