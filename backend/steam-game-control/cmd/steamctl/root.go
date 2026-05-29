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

// NewRootCmd creates and returns the root command
func NewRootCmd() *cobra.Command {
	return rootCmd
}

func init() {
	// Here you can add commands
}
