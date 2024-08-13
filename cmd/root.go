package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// base command
var rootCmd = &cobra.Command{
	Use:   "ttg",
	Short: "A simple wy to track your time.",
	Long: `Time to Go is a CLI application that enables the user to track their time spent on defined projects. 
	The goal is to have a quick and reliable way to keep track and being able to receive a report.`,

	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) {
	// 	fmt.Printf("Your args: %v", args)
	// },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.time-to-go.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
