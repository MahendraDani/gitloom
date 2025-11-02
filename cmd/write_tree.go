package cmd

import (
	"fmt"
	"os"

	"github.com/MahendraDani/gitloom.git/internal/repo"
	"github.com/MahendraDani/gitloom.git/internal/tree"
	"github.com/spf13/cobra"
)

var writeTreeCmd = &cobra.Command{
	Use:   "write-tree",
	Short: "Create a tree object from the current directory",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Get current working directory
		dir, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get current directory: %v", err)
		}

		// Find and load the repository
		r, err := repo.FindRepository(dir)
		if err != nil {
			return fmt.Errorf("not a gitloom repository: %v", err)
		}

		// Write the tree object
		hash, err := tree.WriteTree(dir, r)
		if err != nil {
			return fmt.Errorf("failed to write tree: %v", err)
		}

		// Print the resulting tree hash
		fmt.Println(hash)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(writeTreeCmd)
}
