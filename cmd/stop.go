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
	Use:   "stop [project]",
	Short: "Stop tracking time",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		dir := fmt.Sprintf("~/.local/state/ttg/%s.json", args[0])
		data, err := os.ReadFile(dir)
		if err != nil {
			fmt.Printf("No ongoing tracking found: %v", err)
			return
		}

		var entry TimeEntry
		err = json.Unmarshal(data, &entry)
		if err != nil {
			fmt.Printf("Couldn't serialize entry: %v", err)
			return
		}
		entry.EndTime = time.Now()

		// Save the entry to a history file
		var history []TimeEntry
		historyData, err := os.ReadFile("~/.local/share/ttg/history.json")
		if err != nil {
			fmt.Printf("Couldnt read file: %v", err)
			return
		}
		if len(historyData) > 0 {
			err = json.Unmarshal(historyData, &history)
			if err != nil {
				fmt.Printf("Couldn't serialize entry: %v", err)
				return
			}
		}
		history = append(history, entry)

		newHistoryData, err := json.Marshal(history)
		if err != nil {
			fmt.Printf("Couldn't serialize entry: %v", err)
			return
		}
		err = os.WriteFile("~/.local/share/ttg/history.json", newHistoryData, 0644)
		if err != nil {
			fmt.Printf("Couldnt write file: %v", err)
			return
		}

		// Delete the current tracking file
		err = os.Remove(dir)
		if err != nil {
			fmt.Printf("Couldnt remove directory: %v", err)
			return
		}

		fmt.Printf("Stopped tracking time for project: %s\n", entry.Project)
	},
}

func init() {
	rootCmd.AddCommand(stopCmd)
}
