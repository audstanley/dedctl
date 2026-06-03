package cmd

import (
	"crypto/sha512"
	"fmt"

	"github.com/spf13/cobra"
	"golang.org/x/crypto/bcrypt"
)

var hashCmd = &cobra.Command{
	Use:   "hash <password> <type>",
	Short: "Generate a password hash",
	Long:  `Generate a password hash for use in config.yaml. Supported types: sha512, bcrypt, plain`,
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		password := args[0]
		hashType := args[1]

		var result string

		switch hashType {
		case "sha512":
			result = fmt.Sprintf("%x", sha512.Sum512([]byte(password)))
		case "bcrypt":
			bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
			if err != nil {
				fmt.Fprintf(cmd.OutOrStderr(), "error generating bcrypt hash: %v\n", err)
				return
			}
			result = string(bytes)
		case "plain":
			result = password
		default:
			cmd.Printf("error: unsupported hash type '%s'. Supported types: sha512, bcrypt, plain\n", hashType)
			return
		}

		cmd.Println(result)
	},
}

// NewHashCmd returns the hash command
func NewHashCmd() *cobra.Command {
	return hashCmd
}

func init() {
	rootCmd.AddCommand(hashCmd)
}
