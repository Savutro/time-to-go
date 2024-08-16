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

var (
	fromDate        string
	toDate          string
	ignoreProjects  []string
	includeProjects []string
	outputJSON      bool
	outputCSV       bool
)

var reportCmd = &cobra.Command{
	Use:   "report",
	Short: "Generate a report for a specific time frame",
	Run: func(cmd *cobra.Command, args []string) {
		// Load and validate the project list
		projectList, err := loadProjectList()
		if err != nil {
			fmt.Println(err)
			return
		}

		// Validate included projects
		for _, project := range includeProjects {
			if err := validateProjectName(project, projectList); err != nil {
				fmt.Println(err)
				return
			}
		}

		// Validate ignored projects
		for _, project := range ignoreProjects {
			if err := validateProjectName(project, projectList); err != nil {
				fmt.Println(err)
				return
			}
		}

		// Get the current user and construct the history file path
		usr, err := user.Current()
		if err != nil {
			fmt.Printf("Couldn't get current user: %v\n", err)
			return
		}
		historyFilePath := filepath.Join(usr.HomeDir, ".local", "share", "ttg", "history.json")

		// Read the history file
		var history []TimeEntry
		historyData, err := os.ReadFile(historyFilePath)
		if err != nil {
			fmt.Println("No history found.")
			return
		}

		// Deserialize the history data
		err = json.Unmarshal(historyData, &history)
		if err != nil {
			fmt.Printf("Couldn't deserialize history: %v\n", err)
			return
		}

		// Determine the from and to dates
		var from, to time.Time

		if fromDate != "" {
			from, err = time.Parse("2006-01-02", fromDate)
			if err != nil {
				fmt.Printf("Invalid from date format: %v\n", err)
				return
			}
		}

		if toDate != "" {
			to, err = time.Parse("2006-01-02", toDate)
			if err != nil {
				fmt.Printf("Invalid to date format: %v\n", err)
				return
			}
		}

		// Default to last week if no dates are provided
		if fromDate == "" && toDate == "" {
			to = time.Now().Truncate(24 * time.Hour)
			from = to.AddDate(0, 0, -7)
		} else if fromDate == "" {
			from = to.AddDate(0, 0, -7)
		} else if toDate == "" {
			to = from.AddDate(0, 0, 7)
		}

		report := make(map[string]time.Duration)
		var detailedReport = make(map[string][]TimeEntry)

		// Generate the report
		for _, entry := range history {
			if entry.EndTime.After(from) && entry.EndTime.Before(to.AddDate(0, 0, 1)) {
				if (len(includeProjects) == 0 || contains(includeProjects, entry.Project)) &&
					!contains(ignoreProjects, entry.Project) {

					duration := entry.EndTime.Sub(entry.StartTime)
					report[entry.Project] += duration
					detailedReport[entry.Project] = append(detailedReport[entry.Project], entry)
				}
			}
		}

		// Output the report in the desired format
		if outputJSON {
			outputJSONReport(detailedReport)
		} else if outputCSV {
			outputCSVReport(detailedReport)
		} else {
			outputTextReport(from, to, report, detailedReport)
		}
	},
}

func outputTextReport(from, to time.Time, report map[string]time.Duration, detailedReport map[string][]TimeEntry) {
	fmt.Printf("\n%s -> %s\n\n",
		from.Format("Mon 02 January 2006"),
		to.Format("Mon 02 January 2006"))

	for project, duration := range report {
		fmt.Printf("%s - %s\n", project, formatDuration(duration))

		for _, entry := range detailedReport[project] {
			fmt.Printf("    %s to %s (%s)\n",
				entry.StartTime.Format("2006-01-02 15:04"),
				entry.EndTime.Format("2006-01-02 15:04"),
				formatDuration(entry.EndTime.Sub(entry.StartTime)))
		}
		fmt.Println()
	}
}

func outputJSONReport(detailedReport map[string][]TimeEntry) {
	reportData, err := json.MarshalIndent(detailedReport, "", "    ")
	if err != nil {
		fmt.Println("Error generating JSON:", err)
		return
	}
	fmt.Println(string(reportData))
}

func outputCSVReport(detailedReport map[string][]TimeEntry) {
	fmt.Println("Project,Start Time,End Time,Duration")
	for project, entries := range detailedReport {
		for _, entry := range entries {
			duration := formatDuration(entry.EndTime.Sub(entry.StartTime))
			fmt.Printf("%s,%s,%s,%s\n",
				project,
				entry.StartTime.Format("2006-01-02 15:04"),
				entry.EndTime.Format("2006-01-02 15:04"),
				duration)
		}
	}
}

func contains(slice []string, item string) bool {
	for _, each := range slice {
		if each == item {
			return true
		}
	}
	return false
}

func init() {
	// Add the flags for the report command with aliases
	reportCmd.Flags().StringVarP(&fromDate, "from", "f", "", "Start date for the report (format: YYYY-MM-DD)")
	reportCmd.Flags().StringVarP(&toDate, "to", "t", "", "End date for the report (format: YYYY-MM-DD)")
	reportCmd.Flags().StringSliceVarP(&ignoreProjects, "ignore-project", "i", []string{}, "Projects to ignore in the report")
	reportCmd.Flags().StringSliceVarP(&includeProjects, "project", "p", []string{}, "Projects to include in the report (ignore others)")
	reportCmd.Flags().BoolVarP(&outputJSON, "json", "j", false, "Output report in JSON format")
	reportCmd.Flags().BoolVarP(&outputCSV, "csv", "c", false, "Output report in CSV format")
	rootCmd.AddCommand(reportCmd)
}
