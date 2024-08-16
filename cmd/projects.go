package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"sort"
	"strings"

	"github.com/spf13/cobra"
)

var addProject string
var deleteProject string

// projectsCmd represents the projects command
var projectsCmd = &cobra.Command{
	Use:   "projects",
	Short: "Lists all defined projects and their status",
	Long: `Lists all unique projects that have been tracked. The projects are retrieved 
from the history file and the state directory. Additionally, it indicates 
if a project has an ongoing session based on the presence of its state file.
You can also add or delete projects from the list using flags.`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Get the current user and construct necessary paths
		usr, err := user.Current()
		if err != nil {
			fmt.Printf("Couldn't get current user: %v\n", err)
			return
		}
		historyFilePath := filepath.Join(usr.HomeDir, ".local", "share", "ttg", "history.json")
		projectListPath := filepath.Join(usr.HomeDir, ".local", "share", "ttg", "projects.json")
		stateDir := filepath.Join(usr.HomeDir, ".local", "state", "ttg")

		// Load the existing project list
		var projectList []string
		projectData, err := os.ReadFile(projectListPath)
		if err == nil {
			err = json.Unmarshal(projectData, &projectList)
			if err != nil {
				fmt.Printf("Couldn't deserialize project list: %v\n", err)
				return
			}
		} else if !os.IsNotExist(err) {
			fmt.Printf("Error reading project list file: %v\n", err)
			return
		}

		// Handle adding a new project
		if addProject != "" {
			if !contains(projectList, addProject) {
				projectList = append(projectList, addProject)
				saveProjectList(projectListPath, projectList)
				fmt.Printf("Added project: %s\n", addProject)
			} else {
				fmt.Printf("Project %s already exists.\n", addProject)
			}
			return
		}

		// Handle deleting a project
		if deleteProject != "" {
			if contains(projectList, deleteProject) {
				projectList = remove(projectList, deleteProject)
				saveProjectList(projectListPath, projectList)
				fmt.Printf("Deleted project: %s\n", deleteProject)
			} else {
				fmt.Printf("Project %s not found in the list.\n", deleteProject)
			}
			return
		}

		// Use a map to store unique projects and their status
		projectsStatus := make(map[string]string)

		// Add projects from the history file
		var history []TimeEntry
		historyData, err := os.ReadFile(historyFilePath)
		if err == nil {
			err = json.Unmarshal(historyData, &history)
			if err != nil {
				fmt.Printf("Couldn't deserialize history: %v\n", err)
				return
			}

			for _, entry := range history {
				projectsStatus[entry.Project] = "no session"
			}
		}

		// Check if there are ongoing projects in the state directory
		files, err := os.ReadDir(stateDir)
		if err != nil {
			fmt.Printf("Error reading state directory: %v\n", err)
			return
		}

		for _, file := range files {
			if file.IsDir() {
				continue
			}

			project := strings.TrimSuffix(file.Name(), filepath.Ext(file.Name()))
			projectsStatus[project] = "session"
		}

		// Add projects from the manually maintained list
		for _, project := range projectList {
			if _, exists := projectsStatus[project]; !exists {
				projectsStatus[project] = "no entry"
			}
		}

		// Extract and sort the unique project names
		var projects []string
		for project := range projectsStatus {
			projects = append(projects, project)
		}
		sort.Strings(projects)

		// Print the projects with their status
		if len(projects) > 0 {
			fmt.Println("Tracked projects:")
			for _, project := range projects {
				statusSymbol := "-"
				switch projectsStatus[project] {
				case "session":
					statusSymbol = "[*]"
				case "no session":
					statusSymbol = "[ ]"
				case "no entry":
					statusSymbol = "[x]"
				}
				fmt.Printf(" %s %s\n", statusSymbol, project)
			}
		} else {
			fmt.Println("No projects found.")
		}
	},
}

// Save the updated project list to the file
func saveProjectList(path string, projectList []string) {
	data, err := json.Marshal(projectList)
	if err != nil {
		fmt.Printf("Couldn't serialize project list: %v\n", err)
		return
	}
	err = os.WriteFile(path, data, 0644)
	if err != nil {
		fmt.Printf("Couldn't write project list: %v\n", err)
	}
}

// Helper function to remove a string from a slice
func remove(slice []string, item string) []string {
	var newSlice []string
	for _, each := range slice {
		if each != item {
			newSlice = append(newSlice, each)
		}
	}
	return newSlice
}

func init() {
	// Add the projects command to the root command
	rootCmd.AddCommand(projectsCmd)

	// Add flags for adding and deleting projects
	projectsCmd.Flags().StringVarP(&addProject, "add", "a", "", "Add a new project to the project list")
	projectsCmd.Flags().StringVarP(&deleteProject, "delete", "d", "", "Delete a project from the project list")
}
