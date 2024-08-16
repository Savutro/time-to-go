package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Lists all ongoing sessions with their start times and elapsed time",
	Long: `Lists all ongoing sessions. For each ongoing session, it provides the start time
and the elapsed time since the session started.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Get the current user and construct necessary path
		usr, err := user.Current()
		if err != nil {
			fmt.Printf("Couldn't get current user: %v\n", err)
			return
		}
		stateDir := filepath.Join(usr.HomeDir, ".local", "state", "ttg")

		// Read all files in the state directory
		files, err := os.ReadDir(stateDir)
		if err != nil {
			fmt.Printf("Error reading state directory: %v\n", err)
			return
		}

		// Process each file to get ongoing sessions
		var ongoingSessions []struct {
			Project   string
			StartTime time.Time
			Elapsed   time.Duration
		}

		for _, file := range files {
			if file.IsDir() {
				continue
			}

			project := strings.TrimSuffix(file.Name(), filepath.Ext(file.Name()))
			filePath := filepath.Join(stateDir, file.Name())

			// Read and unmarshal the JSON data
			data, err := os.ReadFile(filePath)
			if err != nil {
				fmt.Printf("Error reading file %s: %v\n", filePath, err)
				continue
			}

			var entry TimeEntry
			err = json.Unmarshal(data, &entry)
			if err != nil {
				fmt.Printf("Error unmarshaling data for project %s: %v\n", project, err)
				continue
			}

			// Calculate elapsed time
			elapsed := time.Since(entry.StartTime)

			ongoingSessions = append(ongoingSessions, struct {
				Project   string
				StartTime time.Time
				Elapsed   time.Duration
			}{
				Project:   project,
				StartTime: entry.StartTime,
				Elapsed:   elapsed,
			})
		}

		// Print the ongoing sessions with indentation
		if len(ongoingSessions) > 0 {
			fmt.Println("Ongoing sessions:")
			for _, session := range ongoingSessions {
				fmt.Printf("  Project: %s\n", session.Project)
				fmt.Printf("    Started: %s\n", session.StartTime.Format("2006-01-02 15:04:05"))
				fmt.Printf("    Elapsed: %s\n", formatDuration(session.Elapsed))
				fmt.Println()
			}
		} else {
			fmt.Println("No ongoing sessions.")
		}
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
