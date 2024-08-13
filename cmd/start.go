package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
)

type TimeEntry struct {
	Project   string    `json:"project"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
}

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start [project name]",
	Short: "Start tracking time for a project",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		project := args[0]
		entry := TimeEntry{
			Project:   project,
			StartTime: time.Now(),
		}

		data, _ := json.Marshal(entry)
		os.WriteFile("current.json", data, 0644)
		fmt.Printf("Started tracking time for project: %s\n", project)
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}
