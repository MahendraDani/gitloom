package cmd

import (
	"fmt"
	"os"

	"github.com/MahendraDani/gitloom.git/internal/object"
	"github.com/MahendraDani/gitloom.git/internal/repo"
	"github.com/spf13/cobra"
)

var writeFlag bool

var hashObjectCmd = &cobra.Command{
	Use:   "hash-object <file>",
	Short: "Compute SHA-1 hash of a file and optionally store it as a gitloom object",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		file := args[0]

		// Find the gitloom repo starting from current directory
		r, err := repo.FindRepo(".")
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}

		hash, err := object.HashObject(file, r, writeFlag)
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}

		fmt.Println(hash) // prints the SHA-1 hash
	},
}

func init() {
	rootCmd.AddCommand(hashObjectCmd)
	hashObjectCmd.Flags().BoolVarP(&writeFlag, "write", "w", false, "Write object to gitloom repository")
}
