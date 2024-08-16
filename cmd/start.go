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

		entry := TimeEntry{
			Project:   project,
			StartTime: time.Now(),
		}

		// Existing logic to start tracking...
		usr, err := user.Current()
		if err != nil {
			fmt.Printf("Couldn't get current user: %v\n", err)
			return
		}

		dir := filepath.Join(usr.HomeDir, ".local", "state", "ttg")
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			fmt.Printf("Couldn't create directory: %v\n", err)
			return
		}

		filePath := filepath.Join(dir, fmt.Sprintf("%s.json", project))

		data, err := json.Marshal(entry)
		if err != nil {
			fmt.Printf("Couldn't serialize entry: %v", err)
			return
		}

		err = os.WriteFile(filePath, data, 0644)
		if err != nil {
			fmt.Printf("Couldn't write file: %v", err)
			return
		}

		fmt.Printf("Started tracking time for project: %s\n", project)
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}
