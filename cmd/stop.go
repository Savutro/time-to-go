package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
)

// stopCmd represents the stop command
var deleteSession bool

var stopCmd = &cobra.Command{
	Use:   "stop [project]",
	Short: "Stop tracking time",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Default project name to "other" if not provided
		project := "other"
		if len(args) > 0 {
			project = args[0]
		}

		// Load and validate the project list
		projectList, err := loadProjectList()
		if err != nil {
			fmt.Println(err)
			return
		}
		if err := validateProjectName(project, projectList); err != nil {
			fmt.Println(err)
			return
		}

		// Get the current user and construct the file paths
		usr, err := user.Current()
		if err != nil {
			fmt.Printf("Couldn't get current user: %v\n", err)
			return
		}

		dir := filepath.Join(usr.HomeDir, ".local", "state", "ttg")
		filePath := filepath.Join(dir, fmt.Sprintf("%s.json", project))

		// Read the current tracking file
		data, err := os.ReadFile(filePath)
		if err != nil {
			fmt.Printf("No ongoing tracking found: %v\n", err)
			return
		}

		// If delete flag is set, remove the tracking file and exit
		if deleteSession {
			err = os.Remove(filePath)
			if err != nil {
				fmt.Printf("Couldn't delete tracking session for project '%s': %v\n", project, err)
				return
			}
			fmt.Printf("Deleted tracking session for project: %s\n", project)
			return
		}

		// Existing logic to stop tracking...
		var entry TimeEntry
		err = json.Unmarshal(data, &entry)
		if err != nil {
			fmt.Printf("Couldn't deserialize entry: %v\n", err)
			return
		}
		entry.EndTime = time.Now()

		historyDir := filepath.Join(usr.HomeDir, ".local", "share", "ttg")
		historyFilePath := filepath.Join(historyDir, "history.json")

		err = os.MkdirAll(historyDir, 0755)
		if err != nil {
			fmt.Printf("Couldn't create history directory: %v\n", err)
			return
		}

		var history []TimeEntry
		historyData, err := os.ReadFile(historyFilePath)
		if err == nil && len(historyData) > 0 {
			err = json.Unmarshal(historyData, &history)
			if err != nil {
				fmt.Printf("Couldn't deserialize history: %v\n", err)
				return
			}
		}

		history = append(history, entry)

		newHistoryData, err := json.Marshal(history)
		if err != nil {
			fmt.Printf("Couldn't serialize history: %v\n", err)
			return
		}
		err = os.WriteFile(historyFilePath, newHistoryData, 0644)
		if err != nil {
			fmt.Printf("Couldn't write history file: %v\n", err)
			return
		}

		err = os.Remove(filePath)
		if err != nil {
			fmt.Printf("Couldn't remove tracking file: %v\n", err)
			return
		}

		fmt.Printf("Stopped tracking time for project: %s\n", entry.Project)
	},
}

func init() {
	// Add the delete flag to the stop command
	stopCmd.Flags().BoolVarP(&deleteSession, "delete", "d", false, "Delete the current tracking session without saving it to the history")
	rootCmd.AddCommand(stopCmd)
}
