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
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Default project name to "other" if not provided
		project := "other"
		if len(args) > 0 {
			project = args[0]
		}

		entry := TimeEntry{
			Project:   project,
			StartTime: time.Now(),
		}

		dir := fmt.Sprintf("~/.local/state/ttg/%s.json", project)

		data, err := json.Marshal(entry)
		if err != nil {
			fmt.Printf("Couldn't serialize entry: %v", err)
			return
		}

		err = os.WriteFile(dir, data, 0644)
		if err != nil {
			fmt.Printf("Couldnt write file: %v", err)
			return
		}

		fmt.Printf("Started tracking time for project: %s\n", project)
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}
