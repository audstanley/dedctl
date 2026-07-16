package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "dedctl",
	Short: "Dedctl - Dedicated Game Controller",
	Long:  `A CLI tool to manage dedicated Steam game servers`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("dedctl")
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
