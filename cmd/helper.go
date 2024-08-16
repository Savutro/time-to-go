package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"time"
)

// Load the project list from the file
func loadProjectList() ([]string, error) {
	usr, err := user.Current()
	if err != nil {
		return nil, fmt.Errorf("couldn't get current user: %v", err)
	}
	projectListPath := filepath.Join(usr.HomeDir, ".local", "share", "ttg", "projects.json")

	var projectList []string
	projectData, err := os.ReadFile(projectListPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("project list file not found")
		}
		return nil, fmt.Errorf("error reading project list file: %v", err)
	}

	err = json.Unmarshal(projectData, &projectList)
	if err != nil {
		return nil, fmt.Errorf("couldn't deserialize project list: %v", err)
	}

	return projectList, nil
}

// Validate if the given project name is in the project list
func validateProjectName(project string, projectList []string) error {
	for _, p := range projectList {
		if p == project || p == "other" {
			return nil
		}
	}
	return fmt.Errorf("project '%s' is not in the project list", project)
}

func formatDuration(d time.Duration) string {
	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60
	seconds := int(d.Seconds()) % 60

	if hours > 0 {
		return fmt.Sprintf("%dh%02dm%02ds", hours, minutes, seconds)
	} else if minutes > 0 {
		return fmt.Sprintf("%dm%02ds", minutes, seconds)
	}
	return fmt.Sprintf("%ds", seconds)
}
