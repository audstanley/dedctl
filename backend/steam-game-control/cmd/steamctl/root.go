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

var metadataDir string

// NewRootCmd creates and returns the root command
func NewRootCmd() *cobra.Command {
	return rootCmd
}

func init() {
	rootCmd.PersistentFlags().StringVar(&metadataDir, "metadata", "", "metadata directory (default is same directory as config.yaml)")

	rootCmd.AddCommand(NewCacheImagesCmd())
	rootCmd.AddCommand(NewHashCmd())
}

// resolveMetadataDir returns the metadata directory path.
func resolveMetadataDir() string {
	if metadataDir != "" {
		return metadataDir
	}
	// Default to ./configs
	return "./configs"
}
