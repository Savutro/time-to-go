package cmd

import (
	"encoding/json"
	"fmt"
	"os"
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
		var history []TimeEntry
		historyData, err := os.ReadFile("history.json")
		if err != nil {
			fmt.Println("No history found.")
			return
		}
		json.Unmarshal(historyData, &history)

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

		if fromDate == "" && toDate == "" {
			// Default: Last week, ending today
			to = time.Now().Truncate(24 * time.Hour)
			from = to.AddDate(0, 0, -7)
		} else if fromDate == "" {
			// If only toDate is specified, report from one week before toDate
			from = to.AddDate(0, 0, -7)
		} else if toDate == "" {
			// If only fromDate is specified, report from that day to one week later
			to = from.AddDate(0, 0, 7)
		}

		report := make(map[string]time.Duration)
		var detailedReport = make(map[string][]TimeEntry)

		for _, entry := range history {
			if entry.EndTime.After(from) && entry.EndTime.Before(to.AddDate(0, 0, 1)) {
				// Filter by project name if needed
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
