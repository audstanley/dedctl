package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "steamctl",
	Short: "Steam Game Server Control API",
	Long:  `A CLI tool to control Steam game servers via systemctl`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Steam Game Server Control API")
	},
}

// Execute executes the root command.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		fmt.Println(err)
	}
}

func init() {
	// Here you can add commands
}
